package svc

import (
	"encoding/json"
	"net/http"
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revaCtx "github.com/cs3org/reva/v2/pkg/ctx"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	revaTags "github.com/cs3org/reva/v2/pkg/tags"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	"go-micro.dev/v4/metadata"
)

func (g Graph) GetTags(w http.ResponseWriter, r *http.Request) {
	th := r.Header.Get(revactx.TokenHeader)
	ctx := revactx.ContextSetToken(r.Context(), th)
	ctx = metadata.Set(ctx, revactx.TokenHeader, th)
	sr, err := g.searchService.Search(ctx, &searchsvc.SearchRequest{
		Query:    "Tags:*",
		PageSize: -1,
	})
	if err != nil {
		g.logger.Error().Err(err).Msg("Could not search for tags")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tagList := revaTags.FromList("")
	for _, match := range sr.Matches {
		for _, tag := range match.Entity.Tags {
			tagList.AddList(tag)
		}
	}

	tagCollection := libregraph.NewCollectionOfTags()
	tagCollection.Value = tagList.AsSlice()

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tagCollection)
}

func (g Graph) AssignTags(w http.ResponseWriter, r *http.Request) {
	var (
		assignment libregraph.TagAssignment
		ctx        = r.Context()
	)

	if err := json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		g.logger.Debug().Err(err).Interface("body", r.Body).Msg("could not decode tag assignment request")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid body schema definition")
		return
	}

	rid, err := storagespace.ParseID(assignment.ResourceId)
	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse resourceId")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid resourceId")
		return
	}

	sres, err := g.gatewayClient.Stat(ctx, &provider.StatRequest{
		Ref: &provider.Reference{ResourceId: &rid},
	})
	if err != nil {
		g.logger.Error().Err(err).Msg("error stating file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sres.GetStatus().GetCode() != rpc.Code_CODE_OK {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "can't stat resource")
		return
	}

	pm := sres.GetInfo().GetPermissionSet()
	if pm == nil {
		g.logger.Error().Err(err).Msg("no permissionset on file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// it says we need "write access" to set tags. One of those should do
	if !pm.InitiateFileUpload && !pm.CreateContainer {
		g.logger.Info().Msg("no permission to create a tag")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var currentTags string
	if m := sres.GetInfo().GetArbitraryMetadata().GetMetadata(); m != nil {
		currentTags = m["tags"]
	}

	allTags := revaTags.FromList(currentTags)
	newTags := strings.Join(assignment.Tags, ",")
	if !allTags.AddList(newTags) {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "no new tags in createtagsrequest or maximum reached")
		return
	}

	resp, err := g.gatewayClient.SetArbitraryMetadata(ctx, &provider.SetArbitraryMetadataRequest{
		Ref: &provider.Reference{ResourceId: &rid},
		ArbitraryMetadata: &provider.ArbitraryMetadata{
			Metadata: map[string]string{
				"tags": allTags.AsList(),
			},
		},
	})
	if err != nil || resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		g.logger.Error().Err(err).Msg("error setting tags")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if g.eventsPublisher != nil {
		ev := events.TagsAdded{
			Tags: newTags,
			Ref: &provider.Reference{
				ResourceId: &rid,
				Path:       ".",
			},
			SpaceOwner: sres.Info.Owner,
			Executant:  revaCtx.ContextMustGetUser(r.Context()).Id,
		}
		if err := events.Publish(g.eventsPublisher, ev); err != nil {
			g.logger.Error().Err(err).Msg("Failed to publish TagsAdded event")
		}
	}
}

func (g Graph) UnassignTags(w http.ResponseWriter, r *http.Request) {
	var (
		unassignment libregraph.TagUnassignment
		ctx          = r.Context()
	)

	if err := json.NewDecoder(r.Body).Decode(&unassignment); err != nil {
		g.logger.Debug().Err(err).Interface("body", r.Body).Msg("could not decode tag assignment request")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid body schema definition")
		return
	}

	rid, err := storagespace.ParseID(unassignment.ResourceId)
	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse resourceId")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid resourceId")
		return
	}

	sres, err := g.gatewayClient.Stat(ctx, &provider.StatRequest{
		Ref: &provider.Reference{ResourceId: &rid},
	})
	if err != nil {
		g.logger.Error().Err(err).Msg("error stating file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sres.GetStatus().GetCode() != rpc.Code_CODE_OK {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "can't stat resource")
		return
	}

	pm := sres.GetInfo().GetPermissionSet()
	if pm == nil {
		g.logger.Error().Err(err).Msg("no permissionset on file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// it says we need "write access" to set tags. One of those should do
	if !pm.InitiateFileUpload && !pm.CreateContainer {
		g.logger.Info().Msg("no permission to create a tag")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var currentTags string
	if m := sres.GetInfo().GetArbitraryMetadata().GetMetadata(); m != nil {
		currentTags = m["tags"]
	}

	allTags := revaTags.FromList(currentTags)
	toDelete := strings.Join(unassignment.Tags, ",")
	if !allTags.RemoveList(toDelete) {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "no new tags in createtagsrequest or maximum reached")
		return
	}

	resp, err := g.gatewayClient.SetArbitraryMetadata(ctx, &provider.SetArbitraryMetadataRequest{
		Ref: &provider.Reference{ResourceId: &rid},
		ArbitraryMetadata: &provider.ArbitraryMetadata{
			Metadata: map[string]string{
				"tags": allTags.AsList(),
			},
		},
	})
	if err != nil || resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		g.logger.Error().Err(err).Msg("error setting tags")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if g.eventsPublisher != nil {
		ev := events.TagsRemoved{
			Tags: toDelete,
			Ref: &provider.Reference{
				ResourceId: &rid,
				Path:       ".",
			},
			SpaceOwner: sres.Info.Owner,
			Executant:  revaCtx.ContextMustGetUser(r.Context()).Id,
		}
		if err := events.Publish(g.eventsPublisher, ev); err != nil {
			g.logger.Error().Err(err).Msg("Failed to publish TagsAdded event")
		}
	}
}
