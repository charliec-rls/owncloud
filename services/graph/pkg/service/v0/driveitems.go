package svc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"

	"github.com/cs3org/reva/v2/pkg/publicshare"
	"github.com/cs3org/reva/v2/pkg/share"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"

	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	"github.com/owncloud/ocis/v2/services/graph/pkg/validate"
)

// GetRootDriveChildren implements the Service interface.
func (g Graph) GetRootDriveChildren(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling GetRootDriveChildren")
	ctx := r.Context()

	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, "could not select next gateway client, aborting")
		return
	}

	currentUser := revactx.ContextMustGetUser(r.Context())
	// do we need to list all or only the personal drive
	filters := []*storageprovider.ListStorageSpacesRequest_Filter{}
	filters = append(filters, listStorageSpacesUserFilter(currentUser.GetId().OpaqueId))
	filters = append(filters, listStorageSpacesTypeFilter("personal"))

	res, err := gatewayClient.ListStorageSpaces(ctx, &storageprovider.ListStorageSpacesRequest{
		Filters: filters,
	})
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg("error making ListStorageSpaces grpc call")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
			return
		}
		g.logger.Error().Err(err).Msg("error sending ListStorageSpaces grpc request")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	var space *storageprovider.StorageSpace
	for _, s := range res.StorageSpaces {
		if utils.UserIDEqual(currentUser.GetId(), s.GetOwner().GetId()) {
			space = s
		}
	}

	lRes, err := gatewayClient.ListContainer(ctx, &storageprovider.ListContainerRequest{
		Ref: &storageprovider.Reference{ResourceId: space.Root},
	})
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg("error making ListContainer grpc call")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case lRes.Status.Code != cs3rpc.Code_CODE_OK:
		if lRes.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, lRes.Status.Message)
			return
		}
		if lRes.Status.Code == cs3rpc.Code_CODE_PERMISSION_DENIED {
			// TODO check if we should return 404 to not disclose existing items
			errorcode.AccessDenied.Render(w, r, http.StatusForbidden, lRes.Status.Message)
			return
		}
		g.logger.Error().Err(err).Msg("error sending list container grpc request")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	files, err := formatDriveItems(g.logger, lRes.Infos)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response as json")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: files})
}

// GetDriveItem returns a driveItem
func (g Graph) GetDriveItem(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling GetDriveItem")
	ctx := r.Context()

	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}
	driveItemID, err := parseIDParam(r, "driveItemID")
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}
	if driveID.StorageId != driveItemID.StorageId || driveID.SpaceId != driveItemID.SpaceId {
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, "Item does not exist")
		return
	}
	/*
		sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
		// Parse the request with odata parser
		odataReq, err := godata.ParseRequest(ctx, sanitizedPath, r.URL.Query())
		if err != nil {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
			return
		}
	*/

	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	res, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: &driveItemID}})
	switch {
	case err != nil:
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code == cs3rpc.Code_CODE_OK:
		// ok
	case res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND:
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
		return
	case res.Status.Code == cs3rpc.Code_CODE_PERMISSION_DENIED:
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message) // do not leak existence? check what graph does
		return
	case res.Status.Code == cs3rpc.Code_CODE_UNAUTHENTICATED:
		errorcode.Unauthenticated.Render(w, r, http.StatusUnauthorized, res.Status.Message) // do not leak existence? check what graph does
		return
	default:
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}
	driveItem, err := cs3ResourceToDriveItem(g.logger, res.Info)
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &driveItem)
}

// GetDriveItemChildren lists the children of a driveItem
func (g Graph) GetDriveItemChildren(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling GetDriveItemChildren")
	ctx := r.Context()

	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}
	driveItemID, err := parseIDParam(r, "driveItemID")
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}
	if driveID.StorageId != driveItemID.StorageId || driveID.SpaceId != driveItemID.SpaceId {
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, "Item does not exist")
		return
	}
	/*
		sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
		// Parse the request with odata parser
		odataReq, err := godata.ParseRequest(ctx, sanitizedPath, r.URL.Query())
		if err != nil {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
			return
		}
	*/

	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := gatewayClient.ListContainer(ctx, &storageprovider.ListContainerRequest{
		Ref: &storageprovider.Reference{ResourceId: &driveItemID},
	})
	switch {
	case err != nil:
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code == cs3rpc.Code_CODE_OK:
		// ok
	case res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND:
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
		return
	case res.Status.Code == cs3rpc.Code_CODE_PERMISSION_DENIED:
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message) // do not leak existence? check what graph does
		return
	case res.Status.Code == cs3rpc.Code_CODE_UNAUTHENTICATED:
		errorcode.Unauthenticated.Render(w, r, http.StatusUnauthorized, res.Status.Message) // do not leak existence? check what graph does
		return
	default:
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	files, err := formatDriveItems(g.logger, res.Infos)
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: files})
}

// ListPermissions lists the permissions of a driveItem
func (g Graph) ListPermissions(w http.ResponseWriter, r *http.Request) {
	gatewayClient, ok := g.GetGatewayClient(w, r)
	if !ok {
		return
	}

	_, itemID, err := g.GetDriveAndItemIDParam(r)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	ctx := r.Context()

	statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: &itemID}})
	if errCode := errorcode.FromStat(statResponse, err); errCode != nil {
		g.logger.Warn().Err(errCode).Interface("stat.res", statResponse)
		errCode.Render(w, r)
		return
	}

	permissionSet := *statResponse.GetInfo().GetPermissionSet()
	allowedActions := unifiedrole.CS3ResourcePermissionsToLibregraphActions(permissionSet)

	collectionOfPermissions := libregraph.CollectionOfPermissionsWithAllowedValues{
		LibreGraphPermissionsActionsAllowedValues: allowedActions,
		LibreGraphPermissionsRolesAllowedValues: conversions.ToValueSlice(
			unifiedrole.GetApplicableRoleDefinitionsForActions(
				allowedActions,
				unifiedrole.UnifiedRoleConditionGrantee,
				g.config.FilesSharing.EnableResharing,
				false,
			),
		),
	}

	for i, definition := range collectionOfPermissions.LibreGraphPermissionsRolesAllowedValues {
		// the openapi spec defines that the rolePermissions should not be part of the response
		definition.RolePermissions = nil
		collectionOfPermissions.LibreGraphPermissionsRolesAllowedValues[i] = definition
	}

	driveItems := make(driveItemsByResourceID)
	driveItems, err = g.listUserShares(ctx, []*collaboration.Filter{
		share.ResourceIDFilter(conversions.ToPointer(itemID)),
	}, driveItems)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	driveItems, err = g.listPublicShares(ctx, []*link.ListPublicSharesRequest_Filter{
		publicshare.ResourceIDFilter(conversions.ToPointer(itemID)),
	}, driveItems)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	for _, driveItem := range driveItems {
		collectionOfPermissions.Value = append(collectionOfPermissions.Value, driveItem.Permissions...)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, collectionOfPermissions)
}

// Invite invites a user to a storage drive (space).
func (g Graph) Invite(w http.ResponseWriter, r *http.Request) {
	gatewayClient, ok := g.GetGatewayClient(w, r)
	if !ok {
		return
	}

	_, itemID, err := g.GetDriveAndItemIDParam(r)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	driveItemInvite := &libregraph.DriveItemInvite{}
	if err := StrictJSONUnmarshal(r.Body, driveItemInvite); err != nil {
		g.logger.Debug().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()

	if err := validate.StructCtx(ctx, driveItemInvite); err != nil {
		g.logger.Debug().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	unifiedRolePermissions := []*libregraph.UnifiedRolePermission{{AllowedResourceActions: driveItemInvite.LibreGraphPermissionsActions}}
	for _, roleId := range driveItemInvite.GetRoles() {
		role, err := unifiedrole.NewUnifiedRoleFromID(roleId, g.config.FilesSharing.EnableResharing)
		if err != nil {
			g.logger.Debug().Err(err).Interface("role", driveItemInvite.GetRoles()[0]).Msg("unable to convert requested role")
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		unifiedRolePermissions = append(unifiedRolePermissions, conversions.ToPointerSlice(role.GetRolePermissions())...)
	}

	statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: &itemID}})
	if errCode := errorcode.FromStat(statResponse, err); errCode != nil {
		g.logger.Warn().Err(errCode).Interface("stat.res", statResponse)
		errCode.Render(w, r)
		return
	}

	createShareErrors := sync.Map{}
	createShareSuccesses := sync.Map{}

	shareCreateGroup, ctx := errgroup.WithContext(ctx)

	for _, driveRecipient := range driveItemInvite.GetRecipients() {
		// not needed anymore with go 1.22 and higher
		driveRecipient := driveRecipient // https://golang.org/doc/faq#closures_and_goroutines,

		shareCreateGroup.Go(func() error {
			objectId := driveRecipient.GetObjectId()

			if objectId == "" {
				return nil
			}

			cs3ResourcePermissions := unifiedrole.PermissionsToCS3ResourcePermissions(unifiedRolePermissions)

			createShareRequest := &collaboration.CreateShareRequest{
				ResourceInfo: statResponse.GetInfo(),
				Grant: &collaboration.ShareGrant{
					Permissions: &collaboration.SharePermissions{
						Permissions: cs3ResourcePermissions,
					},
				},
			}

			permission := &libregraph.Permission{}
			if role := unifiedrole.CS3ResourcePermissionsToUnifiedRole(*cs3ResourcePermissions, unifiedrole.UnifiedRoleConditionGrantee, g.config.FilesSharing.EnableResharing); role != nil {
				permission.Roles = []string{role.GetId()}
			}

			if len(permission.GetRoles()) == 0 {
				permission.LibreGraphPermissionsActions = unifiedrole.CS3ResourcePermissionsToLibregraphActions(*cs3ResourcePermissions)
			}

			switch driveRecipient.GetLibreGraphRecipientType() {
			case "group":
				group, err := g.identityCache.GetGroup(ctx, objectId)
				if err != nil {
					g.logger.Debug().Err(err).Interface("groupId", objectId).Msg("failed group lookup")
					createShareErrors.Store(objectId, errorcode.GeneralException.CreateOdataError(r.Context(), http.StatusText(http.StatusInternalServerError)))
					return nil
				}
				createShareRequest.GetGrant().Grantee = &storageprovider.Grantee{
					Type: storageprovider.GranteeType_GRANTEE_TYPE_GROUP,
					Id: &storageprovider.Grantee_GroupId{GroupId: &grouppb.GroupId{
						OpaqueId: group.GetId(),
					}},
				}
				permission.GrantedToV2 = &libregraph.SharePointIdentitySet{
					Group: &libregraph.Identity{
						DisplayName: group.GetDisplayName(),
						Id:          conversions.ToPointer(group.GetId()),
					},
				}
			default:
				user, err := g.identityCache.GetUser(ctx, objectId)
				if err != nil {
					g.logger.Debug().Err(err).Interface("userId", objectId).Msg("failed user lookup")
					createShareErrors.Store(objectId, errorcode.GeneralException.CreateOdataError(r.Context(), http.StatusText(http.StatusInternalServerError)))
					return nil
				}
				createShareRequest.GetGrant().Grantee = &storageprovider.Grantee{
					Type: storageprovider.GranteeType_GRANTEE_TYPE_USER,
					Id: &storageprovider.Grantee_UserId{UserId: &userpb.UserId{
						OpaqueId: user.GetId(),
					}},
				}
				permission.GrantedToV2 = &libregraph.SharePointIdentitySet{
					User: &libregraph.Identity{
						DisplayName: user.GetDisplayName(),
						Id:          conversions.ToPointer(user.GetId()),
					},
				}
			}

			if driveItemInvite.ExpirationDateTime != nil {
				createShareRequest.GetGrant().Expiration = utils.TimeToTS(*driveItemInvite.ExpirationDateTime)
			}

			createShareResponse, err := gatewayClient.CreateShare(ctx, createShareRequest)
			switch {
			case err != nil:
				fallthrough
			case createShareResponse.GetStatus().GetCode() != cs3rpc.Code_CODE_OK:
				g.logger.Debug().Err(err).Msg("share creation failed")
				createShareErrors.Store(objectId, errorcode.GeneralException.CreateOdataError(r.Context(), http.StatusText(http.StatusInternalServerError)))
				return nil
			}

			if id := createShareResponse.GetShare().GetId().GetOpaqueId(); id != "" {
				permission.Id = conversions.ToPointer(id)
			}

			if expiration := createShareResponse.GetShare().GetExpiration(); expiration != nil {
				permission.SetExpirationDateTime(utils.TSToTime(expiration))
			}

			createShareSuccesses.Store(objectId, permission)

			return nil
		})
	}

	if err := shareCreateGroup.Wait(); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	value := make([]interface{}, 0, len(driveItemInvite.Recipients))

	hasSuccesses := false
	createShareSuccesses.Range(func(key, permission interface{}) bool {
		value = append(value, permission)
		hasSuccesses = true
		return true
	})

	hasErrors := false
	createShareErrors.Range(func(key, err interface{}) bool {
		value = append(value, err)
		hasErrors = true
		return true
	})

	switch {
	case hasErrors && hasSuccesses:
		render.Status(r, http.StatusMultiStatus)
	case hasSuccesses:
		render.Status(r, http.StatusCreated)
	default:
		render.Status(r, http.StatusInternalServerError)
	}

	render.JSON(w, r, &ListResponse{Value: value})
}

// DeletePermission removes a Permission from a Drive item
func (g Graph) DeletePermission(w http.ResponseWriter, r *http.Request) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Debug().Err(err).Msg("selecting gatewaySelector failed")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	_, itemID, err := g.extractDriveAndDriveItem(r)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))

	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse driveID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid driveID")
		return
	}

	ctx := r.Context()
	getShareResp, err := gatewayClient.GetShare(ctx,
		&collaboration.GetShareRequest{
			Ref: &collaboration.ShareReference{
				Spec: &collaboration.ShareReference_Id{
					Id: &collaboration.ShareId{
						OpaqueId: permissionID,
					},
				},
			},
		})
	switch {
	case err != nil:
		fallthrough
	case getShareResp.Status.GetCode() != cs3rpc.Code_CODE_OK:
		g.logger.Debug().Err(err).Interface("permissionID", permissionID).Interface("GetShare", getShareResp).Msg("GetShare failed")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	sharedResourceId := getShareResp.GetShare().GetResourceId()
	// The resourceID of the shared resource need to matched the item ID from the Request Path
	// otherwise this is an invalid Request.
	if sharedResourceId.GetStorageId() != itemID.GetStorageId() ||
		sharedResourceId.GetSpaceId() != itemID.GetSpaceId() ||
		sharedResourceId.GetOpaqueId() != itemID.GetOpaqueId() {
		g.logger.Debug().Msg("resourceID of shared does not match itemID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "permissionID and itemID do not match")
		return
	}

	removeShareResp, err := gatewayClient.RemoveShare(ctx,
		&collaboration.RemoveShareRequest{
			Ref: &collaboration.ShareReference{
				Spec: &collaboration.ShareReference_Id{
					Id: &collaboration.ShareId{
						OpaqueId: permissionID,
					},
				},
			},
		})
	switch {
	case err != nil:
		fallthrough
	case removeShareResp.Status.GetCode() != cs3rpc.Code_CODE_OK:
		g.logger.Debug().Err(err).Interface("permissionID", permissionID).Interface("GetShare", getShareResp).Msg("GetShare failed")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)

	return
}

func (g Graph) extractDriveAndDriveItem(r *http.Request) (driveID storageprovider.ResourceId, itemID storageprovider.ResourceId, err error) {
	driveID, err = parseIDParam(r, "driveID")
	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse driveID")
		return driveID, itemID, errorcode.New(errorcode.InvalidRequest, "invalid driveID")
	}

	itemID, err = parseIDParam(r, "itemID")
	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse itemID")
		return driveID, itemID, errorcode.New(errorcode.InvalidRequest, "invalid itemID")
	}

	if driveID.GetStorageId() != itemID.GetStorageId() || driveID.GetSpaceId() != itemID.GetSpaceId() {
		g.logger.Debug().Interface("driveID", driveID).Interface("itemID", itemID).Msg("driveID and itemID do not match")
		return driveID, itemID, errorcode.New(errorcode.InvalidRequest, "driveID and itemID do not match")
	}
	return driveID, itemID, nil
}

func (g Graph) getDriveItem(ctx context.Context, ref storageprovider.Reference) (*libregraph.DriveItem, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	res, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &ref})
	if err != nil {
		return nil, err
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		refStr, _ := storagespace.FormatReference(&ref)
		return nil, fmt.Errorf("could not stat %s: %s", refStr, res.Status.Message)
	}
	return cs3ResourceToDriveItem(g.logger, res.Info)
}

func (g Graph) getRemoteItem(ctx context.Context, root *storageprovider.ResourceId, baseURL *url.URL) (*libregraph.RemoteItem, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	ref := &storageprovider.Reference{
		ResourceId: root,
	}
	res, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: ref})
	if err != nil {
		return nil, err
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		// Only log this, there could be mountpoints which have no grant
		g.logger.Debug().Msg(res.Status.Message)
		return nil, errors.New("could not fetch grant resource for the mountpoint")
	}
	item, err := cs3ResourceToRemoteItem(res.Info)
	if err != nil {
		return nil, err
	}

	if baseURL != nil && res.GetInfo() != nil && res.GetInfo().GetSpace() != nil {
		// TODO read from StorageSpace ... needs Opaque for now
		// TODO how do we build the url?
		// for now: read from request
		item.Name = libregraph.PtrString(res.GetInfo().GetName())
		if res.GetInfo().GetSpace().GetRoot() != nil {
			webDavURL := *baseURL
			relativePath := res.GetInfo().GetPath()
			webDavURL.Path = path.Join(webDavURL.Path, storagespace.FormatResourceID(*res.GetInfo().GetSpace().GetRoot()), relativePath)
			item.WebDavUrl = libregraph.PtrString(webDavURL.String())
		}
	}
	return item, nil
}

func formatDriveItems(logger *log.Logger, mds []*storageprovider.ResourceInfo) ([]*libregraph.DriveItem, error) {
	responses := make([]*libregraph.DriveItem, 0, len(mds))
	for i := range mds {
		res, err := cs3ResourceToDriveItem(logger, mds[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	return responses, nil
}

func cs3TimestampToTime(t *types.Timestamp) time.Time {
	return time.Unix(int64(t.Seconds), int64(t.Nanos))
}

func cs3ResourceToDriveItem(logger *log.Logger, res *storageprovider.ResourceInfo) (*libregraph.DriveItem, error) {
	size := new(int64)
	*size = int64(res.Size) // TODO lurking overflow: make size of libregraph drive item use uint64

	driveItem := &libregraph.DriveItem{
		Id:   libregraph.PtrString(storagespace.FormatResourceID(*res.Id)),
		Size: size,
	}

	if name := path.Base(res.Path); name != "" {
		driveItem.Name = &name
	}
	if res.Etag != "" {
		driveItem.ETag = &res.Etag
	}
	if res.Mtime != nil {
		lastModified := cs3TimestampToTime(res.Mtime)
		driveItem.LastModifiedDateTime = &lastModified
	}
	if res.ParentId != nil {
		parentRef := libregraph.NewItemReference()
		parentRef.SetDriveType(res.Space.SpaceType)
		parentRef.SetDriveId(storagespace.FormatStorageID(res.ParentId.StorageId, res.ParentId.SpaceId))
		parentRef.SetId(storagespace.FormatResourceID(*res.ParentId))
		driveItem.ParentReference = parentRef
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_FILE && res.MimeType != "" {
		// We cannot use a libregraph.File here because the openapi codegenerator autodetects 'File' as a go type ...
		driveItem.File = &libregraph.OpenGraphFile{
			MimeType: &res.MimeType,
		}
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER {
		driveItem.Folder = &libregraph.Folder{}
	}

	driveItem.Audio = cs3ResourceToDriveItemAudioFacet(logger, res)

	return driveItem, nil
}

func cs3ResourceToDriveItemAudioFacet(logger *log.Logger, res *storageprovider.ResourceInfo) *libregraph.Audio {
	if !strings.HasPrefix(res.MimeType, "audio/") {
		return nil
	}

	k := res.ArbitraryMetadata.Metadata
	if k == nil {
		return nil
	}

	var audio = &libregraph.Audio{}
	if ok := unmarshalStringMap(logger, audio, k, "libre.graph.audio."); ok {
		return audio
	}

	return nil
}

func getFieldName(structField reflect.StructField) string {
	tag := structField.Tag.Get("json")
	if tag == "" {
		return structField.Name
	}

	return strings.Split(tag, ",")[0]
}

func unmarshalStringMap(logger *log.Logger, out any, flatMap map[string]string, prefix string) bool {
	nonEmpty := false
	obj := reflect.ValueOf(out).Elem()
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		structField := obj.Type().Field(i)
		mapKey := prefix + getFieldName(structField)

		if value, ok := flatMap[mapKey]; ok {
			if field.Kind() == reflect.Ptr {
				newValue := reflect.New(field.Type().Elem())
				var tmp any
				var err error
				switch t := newValue.Type().Elem().Kind(); t {
				case reflect.String:
					tmp = value
				case reflect.Int32:
					tmp, err = strconv.ParseInt(value, 10, 32)
				case reflect.Int64:
					tmp, err = strconv.ParseInt(value, 10, 64)
				case reflect.Bool:
					tmp, err = strconv.ParseBool(value)
				default:
					err = errors.New("unsupported type")
					logger.Error().Err(err).Str("type", t.String()).Str("mapKey", mapKey).Msg("target field type for value of mapKey is not supported")
				}
				if err != nil {
					logger.Error().Err(err).Str("mapKey", mapKey).Msg("unmarshalling failed")
					continue
				}
				newValue.Elem().Set(reflect.ValueOf(tmp).Convert(field.Type().Elem()))
				field.Set(newValue)
				nonEmpty = true
			}
		}
	}

	return nonEmpty
}

func cs3ResourceToRemoteItem(res *storageprovider.ResourceInfo) (*libregraph.RemoteItem, error) {
	size := new(int64)
	*size = int64(res.Size) // TODO lurking overflow: make size of libregraph drive item use uint64

	remoteItem := &libregraph.RemoteItem{
		Id:   libregraph.PtrString(storagespace.FormatResourceID(*res.Id)),
		Size: size,
	}

	if res.GetPath() != "" {
		remoteItem.Path = libregraph.PtrString(path.Clean(res.GetPath()))
	}
	if res.Etag != "" {
		remoteItem.ETag = &res.Etag
	}
	if res.Mtime != nil {
		lastModified := cs3TimestampToTime(res.Mtime)
		remoteItem.LastModifiedDateTime = &lastModified
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_FILE && res.MimeType != "" {
		// We cannot use a libregraph.File here because the openapi codegenerator autodetects 'File' as a go type ...
		remoteItem.File = &libregraph.OpenGraphFile{
			MimeType: &res.MimeType,
		}
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER {
		remoteItem.Folder = &libregraph.Folder{}
	}
	if res.GetSpace() != nil && res.GetSpace().GetRoot() != nil {
		remoteItem.RootId = libregraph.PtrString(storagespace.FormatResourceID(*res.GetSpace().GetRoot()))
		grantSpaceAlias := utils.ReadPlainFromOpaque(res.GetSpace().GetOpaque(), "spaceAlias")
		if grantSpaceAlias != "" {
			remoteItem.DriveAlias = libregraph.PtrString(grantSpaceAlias)
		}
	}
	return remoteItem, nil
}

func (g Graph) getPathForResource(ctx context.Context, id storageprovider.ResourceId) (string, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return "", err
	}

	res, err := gatewayClient.GetPath(ctx, &storageprovider.GetPathRequest{ResourceId: &id})
	if err != nil {
		return "", err
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		return "", fmt.Errorf("could not stat %v: %s", id, res.Status.Message)
	}
	return res.Path, err
}

// getSpecialDriveItems reads properties from the opaque and transforms them into driveItems
func (g Graph) getSpecialDriveItems(ctx context.Context, baseURL *url.URL, space *storageprovider.StorageSpace) []libregraph.DriveItem {
	if space.GetRoot().GetStorageId() == utils.ShareStorageProviderID {
		return nil // no point in stating the ShareStorageProvider
	}
	if space.Opaque == nil {
		return nil
	}

	imageNode := utils.ReadPlainFromOpaque(space.Opaque, SpaceImageSpecialFolderName)
	readmeNode := utils.ReadPlainFromOpaque(space.Opaque, ReadmeSpecialFolderName)

	cachekey := spaceRootStatKey(space.Root, imageNode, readmeNode)
	// if the root is older or equal to our cache we can reuse the cached extended spaces properties
	if entry := g.specialDriveItemsCache.Get(cachekey); entry != nil {
		if cached, ok := entry.Value().(specialDriveItemEntry); ok {
			if cached.rootMtime != nil && space.Mtime != nil {
				// beware, LaterTS does not handle equalness. it returns t1 if t1 > t2, else t2, so a >= check looks like this
				if utils.LaterTS(space.Mtime, cached.rootMtime) == cached.rootMtime {
					return cached.specialDriveItems
				}
			}
		}
	}

	var spaceItems []libregraph.DriveItem

	spaceItems = g.fetchSpecialDriveItem(ctx, spaceItems, SpaceImageSpecialFolderName, imageNode, space, baseURL)
	spaceItems = g.fetchSpecialDriveItem(ctx, spaceItems, ReadmeSpecialFolderName, readmeNode, space, baseURL)

	// cache properties
	spacePropertiesEntry := specialDriveItemEntry{
		specialDriveItems: spaceItems,
		rootMtime:         space.Mtime,
	}
	g.specialDriveItemsCache.Set(cachekey, spacePropertiesEntry, time.Duration(g.config.Spaces.ExtendedSpacePropertiesCacheTTL))

	return spaceItems
}

func (g Graph) fetchSpecialDriveItem(ctx context.Context, spaceItems []libregraph.DriveItem, itemName string, itemNode string, space *storageprovider.StorageSpace, baseURL *url.URL) []libregraph.DriveItem {
	var ref storageprovider.Reference
	if itemNode != "" {
		rid, _ := storagespace.ParseID(itemNode)

		rid.StorageId = space.GetRoot().GetStorageId()
		ref = storageprovider.Reference{
			ResourceId: &rid,
		}
		spaceItem := g.getSpecialDriveItem(ctx, ref, itemName, baseURL, space)
		if spaceItem != nil {
			spaceItems = append(spaceItems, *spaceItem)
		}
	}
	return spaceItems
}

// generates a space root stat cache key used to detect changes in a space
// takes into account the special nodes because changing metadata does not affect the etag / mtime
func spaceRootStatKey(id *storageprovider.ResourceId, imagenode, readmeNode string) string {
	if id == nil {
		return ""
	}
	shakeHash := sha3.NewShake256()
	_, _ = shakeHash.Write([]byte(id.GetStorageId()))
	_, _ = shakeHash.Write([]byte(id.GetSpaceId()))
	_, _ = shakeHash.Write([]byte(id.GetOpaqueId()))
	_, _ = shakeHash.Write([]byte(imagenode))
	_, _ = shakeHash.Write([]byte(readmeNode))
	h := make([]byte, 64)
	_, _ = shakeHash.Read(h)
	return fmt.Sprintf("%x", h)
}

type specialDriveItemEntry struct {
	specialDriveItems []libregraph.DriveItem
	rootMtime         *types.Timestamp
}

func (g Graph) getSpecialDriveItem(ctx context.Context, ref storageprovider.Reference, itemName string, baseURL *url.URL, space *storageprovider.StorageSpace) *libregraph.DriveItem {
	var spaceItem *libregraph.DriveItem
	if ref.GetResourceId().GetSpaceId() == "" && ref.GetResourceId().GetOpaqueId() == "" {
		return nil
	}

	// FIXME we should send a fieldmask 'path' and return it as the Path property to save an additional call to the storage.
	// To do that we need to align the useg of ResourceInfo.Name vs ResourceInfo.Path. By default, only the name should be set
	// and Path should always be relative to the space root OR the resource the current user can access ...
	spaceItem, err := g.getDriveItem(ctx, ref)
	if err != nil {
		g.logger.Debug().Err(err).Str("ID", ref.GetResourceId().GetOpaqueId()).Str("name", itemName).Msg("Could not get item info")
		return nil
	}
	itemPath := ref.Path
	if itemPath == "" {
		// lookup by id
		itemPath, err = g.getPathForResource(ctx, *ref.ResourceId)
		if err != nil {
			g.logger.Debug().Err(err).Str("ID", ref.GetResourceId().GetOpaqueId()).Str("name", itemName).Msg("Could not get item path")
			return nil
		}
	}
	spaceItem.SpecialFolder = &libregraph.SpecialFolder{Name: libregraph.PtrString(itemName)}
	webdavURL := *baseURL
	webdavURL.Path = path.Join(webdavURL.Path, space.Id.OpaqueId, itemPath)
	spaceItem.WebDavUrl = libregraph.PtrString(webdavURL.String())

	return spaceItem
}
