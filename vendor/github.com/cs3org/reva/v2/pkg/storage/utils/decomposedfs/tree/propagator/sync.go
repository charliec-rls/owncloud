// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package propagator

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/rogpeppe/go-internal/lockedfile"
)

// SyncPropagator implements synchronous treetime & treesize propagation
type SyncPropagator struct {
	treeSizeAccounting bool
	treeTimeAccounting bool
	lookup             node.PathLookup
}

// NewSyncPropagator returns a new AsyncPropagator instance
func NewSyncPropagator(treeSizeAccounting, treeTimeAccounting bool, lookup node.PathLookup) SyncPropagator {
	return SyncPropagator{
		treeSizeAccounting: treeSizeAccounting,
		treeTimeAccounting: treeTimeAccounting,
		lookup:             lookup,
	}
}

// Propagate triggers a propagation
func (p SyncPropagator) Propagate(ctx context.Context, n *node.Node, sizeDiff int64) error {
	ctx, span := tracer.Start(ctx, "Propagate")
	defer span.End()
	sublog := appctx.GetLogger(ctx).With().
		Str("method", "sync.Propagate").
		Str("spaceid", n.SpaceID).
		Str("nodeid", n.ID).
		Int64("sizeDiff", sizeDiff).
		Logger()

	if !p.treeTimeAccounting && (!p.treeSizeAccounting || sizeDiff == 0) {
		// no propagation enabled
		sublog.Debug().Msg("propagation disabled or nothing to propagate")
		return nil
	}

	// is propagation enabled for the parent node?
	root := n.SpaceRoot

	// use a sync time and don't rely on the mtime of the current node, as the stat might not change when a rename happened too quickly
	sTime := time.Now().UTC()

	// we loop until we reach the root
	var err error
	for err == nil && n.ID != root.ID {
		sublog.Debug().Msg("propagating")

		attrs := node.Attributes{}

		var f *lockedfile.File
		// lock parent before reading treesize or tree time

		_, subspan := tracer.Start(ctx, "lockedfile.OpenFile")
		parentFilename := p.lookup.MetadataBackend().LockfilePath(n.ParentPath())
		f, err = lockedfile.OpenFile(parentFilename, os.O_RDWR|os.O_CREATE, 0600)
		subspan.End()
		if err != nil {
			sublog.Error().Err(err).
				Str("parent filename", parentFilename).
				Msg("Propagation failed. Could not open metadata for parent with lock.")
			return err
		}
		// always log error if closing node fails
		defer func() {
			// ignore already closed error
			cerr := f.Close()
			if err == nil && cerr != nil && !errors.Is(cerr, os.ErrClosed) {
				err = cerr // only overwrite err with en error from close if the former was nil
			}
			_ = os.Remove(f.Name())
		}()

		if n, err = n.Parent(ctx); err != nil {
			sublog.Error().Err(err).
				Msg("Propagation failed. Could not read parent node.")
			return err
		}

		if !n.HasPropagation(ctx) {
			sublog.Debug().Str("attr", prefixes.PropagationAttr).Msg("propagation attribute not set or unreadable, not propagating")
			// if the attribute is not set treat it as false / none / no propagation
			return nil
		}

		sublog = sublog.With().Str("spaceid", n.SpaceID).Str("nodeid", n.ID).Logger()

		if p.treeTimeAccounting {
			// update the parent tree time if it is older than the nodes mtime
			updateSyncTime := false

			var tmTime time.Time
			tmTime, err = n.GetTMTime(ctx)
			switch {
			case err != nil:
				// missing attribute, or invalid format, overwrite
				sublog.Debug().Err(err).
					Msg("could not read tmtime attribute, overwriting")
				updateSyncTime = true
			case tmTime.Before(sTime):
				sublog.Debug().
					Time("tmtime", tmTime).
					Time("stime", sTime).
					Msg("parent tmtime is older than node mtime, updating")
				updateSyncTime = true
			default:
				sublog.Debug().
					Time("tmtime", tmTime).
					Time("stime", sTime).
					Dur("delta", sTime.Sub(tmTime)).
					Msg("parent tmtime is younger than node mtime, not updating")
			}

			if updateSyncTime {
				// update the tree time of the parent node
				attrs.SetString(prefixes.TreeMTimeAttr, sTime.UTC().Format(time.RFC3339Nano))
			}

			attrs.SetString(prefixes.TmpEtagAttr, "")
		}

		// size accounting
		if p.treeSizeAccounting && sizeDiff != 0 {
			var newSize uint64

			// read treesize
			treeSize, err := n.GetTreeSize(ctx)
			switch {
			case metadata.IsAttrUnset(err):
				// fallback to calculating the treesize
				sublog.Warn().Msg("treesize attribute unset, falling back to calculating the treesize")
				newSize, err = calculateTreeSize(ctx, p.lookup, n.InternalPath())
				if err != nil {
					return err
				}
			case err != nil:
				sublog.Error().Err(err).
					Msg("Faild to propagate treesize change. Error when reading the treesize attribute from parent")
				return err
			case sizeDiff > 0:
				newSize = treeSize + uint64(sizeDiff)
			case uint64(-sizeDiff) > treeSize:
				// The sizeDiff is larger than the current treesize. Which would result in
				// a negative new treesize. Something must have gone wrong with the accounting.
				// Reset the current treesize to 0.
				sublog.Error().Uint64("treeSize", treeSize).Int64("sizeDiff", sizeDiff).
					Msg("Error when updating treesize of parent node. Updated treesize < 0. Reestting to 0")
				newSize = 0
			default:
				newSize = treeSize - uint64(-sizeDiff)
			}

			// update the tree size of the node
			attrs.SetString(prefixes.TreesizeAttr, strconv.FormatUint(newSize, 10))
			sublog.Debug().Uint64("newSize", newSize).Msg("updated treesize of parent node")
		}

		if err = n.SetXattrsWithContext(ctx, attrs, false); err != nil {
			sublog.Error().Err(err).Msg("Failed to update extend attributes of parent node")
			return err
		}

		// Release node lock early, ignore already closed error
		_, subspan = tracer.Start(ctx, "f.Close")
		cerr := f.Close()
		subspan.End()
		if cerr != nil && !errors.Is(cerr, os.ErrClosed) {
			sublog.Error().Err(cerr).Msg("Failed to close parent node and release lock")
			return cerr
		}
	}
	if err != nil {
		sublog.Error().Err(err).Msg("error propagating")
		return err
	}
	return nil

}
