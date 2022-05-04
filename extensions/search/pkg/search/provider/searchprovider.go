package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/walker"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/cs3org/reva/v2/pkg/utils/resourceid"
	"github.com/owncloud/ocis/extensions/search/pkg/search"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"google.golang.org/grpc/metadata"

	searchmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
)

type Provider struct {
	logger            log.Logger
	gwClient          gateway.GatewayAPIClient
	indexClient       search.IndexClient
	machineAuthAPIKey string
}

func New(gwClient gateway.GatewayAPIClient, indexClient search.IndexClient, machineAuthAPIKey string, eventsChan <-chan interface{}, logger log.Logger) *Provider {
	p := &Provider{
		gwClient:          gwClient,
		indexClient:       indexClient,
		machineAuthAPIKey: machineAuthAPIKey,
		logger:            logger,
	}

	go func() {
		for {
			ev := <-eventsChan
			go func() {
				time.Sleep(1 * time.Second) // Give some time to let everything settle down before trying to access it when indexing
				p.handleEvent(ev)
			}()
		}
	}()

	return p
}

func (p *Provider) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if req.Query == "" {
		return nil, errtypes.PreconditionFailed("empty query provided")
	}

	listSpacesRes, err := p.gwClient.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{SpaceType: "+grant"},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	matches := []*searchmsg.Match{}
	for _, space := range listSpacesRes.StorageSpaces {
		pathPrefix := ""
		if space.SpaceType == "grant" {
			gpRes, err := p.gwClient.GetPath(ctx, &provider.GetPathRequest{
				ResourceId: space.Root,
			})
			if err != nil {
				return nil, err
			}
			if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
				return nil, errtypes.NewErrtypeFromStatus(gpRes.Status)
			}
			pathPrefix = utils.MakeRelativePath(gpRes.Path)
		}

		rootStorageID, _ := resourceid.StorageIDUnwrap(space.Root.StorageId)
		res, err := p.indexClient.Search(ctx, &searchsvc.SearchIndexRequest{
			Query: req.Query,
			Ref: &searchmsg.Reference{
				ResourceId: &searchmsg.ResourceID{
					StorageId: space.Root.StorageId,
					OpaqueId:  rootStorageID,
				},
				Path: pathPrefix,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, match := range res.Matches {
			if pathPrefix != "" {
				match.Entity.Ref.Path = utils.MakeRelativePath(strings.TrimPrefix(match.Entity.Ref.Path, pathPrefix))
			}
			matches = append(matches, match)
		}
	}

	return &searchsvc.SearchResponse{
		Matches: matches,
	}, nil
}

func (p *Provider) IndexSpace(ctx context.Context, req *searchsvc.IndexSpaceRequest) (*searchsvc.IndexSpaceResponse, error) {
	// get user
	res, err := p.gwClient.GetUserByClaim(context.Background(), &user.GetUserByClaimRequest{
		Claim: "username",
		Value: req.UserId,
	})
	if err != nil || res.Status.Code != rpc.Code_CODE_OK {
		fmt.Println("error: Could not get user by userid")
		return nil, err
	}

	// Get auth context
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), res.User)
	authRes, err := p.gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + res.User.Id.OpaqueId,
		ClientSecret: p.machineAuthAPIKey,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, err
	}

	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not get authenticated context for user")
	}
	ownerCtx = metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token)

	// Walk the space and index all files
	walker := walker.NewWalker(p.gwClient)
	rootId := &provider.ResourceId{StorageId: req.SpaceId, OpaqueId: req.SpaceId}
	err = walker.Walk(ownerCtx, rootId, func(wd string, info *provider.ResourceInfo, err error) error {
		if err != nil {
			p.logger.Error().Err(err).Msg("error walking the tree")
		}
		ref := &provider.Reference{
			Path:       utils.MakeRelativePath(filepath.Join(wd, info.Path)),
			ResourceId: rootId,
		}
		err = p.indexClient.Add(ref, info)
		if err != nil {
			p.logger.Error().Err(err).Msg("error adding resource to the index")
		} else {
			p.logger.Debug().Interface("ref", ref).Msg("added resource to index")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	p.logDocCount()
	return &searchsvc.IndexSpaceResponse{}, nil
}

func (p *Provider) logDocCount() {
	c, err := p.indexClient.DocCount()
	if err != nil {
		p.logger.Error().Err(err).Msg("error getting document count from the index")
	}
	p.logger.Debug().Interface("count", c).Msg("new document count")
}
