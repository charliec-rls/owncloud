package svc

import (
	"context"
	"net/http"
	"net/url"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	"github.com/owncloud/ocis/v2/services/graph/pkg/validate"
)

const invalidIdMsg = "invalid driveID or itemID"

type DriveItemPermissionsProvider interface {
	Invite(ctx context.Context, resourceId storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error)
	SpaceRootInvite(ctx context.Context, driveID storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error)
	ListPermissions(ctx context.Context, itemID storageprovider.ResourceId) (libregraph.CollectionOfPermissionsWithAllowedValues, error)
	ListSpaceRootPermissions(ctx context.Context, driveID storageprovider.ResourceId) (libregraph.CollectionOfPermissionsWithAllowedValues, error)
	DeletePermission(ctx context.Context, itemID storageprovider.ResourceId, permissionID string) error
	UpdatePermission(ctx context.Context, itemID storageprovider.ResourceId, permissionID string, newPermission libregraph.Permission) (libregraph.Permission, error)
}

// DriveItemPermissionsService contains the production business logic for everything that relates to permissions on drive items.
type DriveItemPermissionsService struct {
	BaseGraphService
}

type permissionType int

const (
	Unknown permissionType = iota
	Public
	User
	Space
)

// NewDriveItemPermissionsService creates a new DriveItemPermissionsService
func NewDriveItemPermissionsService(logger log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], identityCache identity.IdentityCache, config *config.Config) (DriveItemPermissionsService, error) {
	return DriveItemPermissionsService{
		BaseGraphService: BaseGraphService{
			logger:          &log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemService").Logger()},
			gatewaySelector: gatewaySelector,
			identityCache:   identityCache,
			config:          config,
		},
	}, nil
}

// Invite invites a user to a drive item.
func (s DriveItemPermissionsService) Invite(ctx context.Context, resourceId storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return libregraph.Permission{}, err
	}

	statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: &resourceId}})
	if errCode := errorcode.FromStat(statResponse, err); errCode != nil {
		s.logger.Warn().Err(errCode).Interface("stat.res", statResponse).Msg("stat failed")
		return libregraph.Permission{}, *errCode
	}

	resourceInfo := statResponse.GetInfo()
	condition := unifiedrole.UnifiedRoleConditionGrantee
	if IsSpaceRoot(resourceInfo.GetId()) {
		condition = unifiedrole.UnifiedRoleConditionOwner
	}

	unifiedRolePermissions := []*libregraph.UnifiedRolePermission{{AllowedResourceActions: invite.LibreGraphPermissionsActions}}
	for _, roleID := range invite.GetRoles() {
		role, err := unifiedrole.NewUnifiedRoleFromID(roleID, s.config.FilesSharing.EnableResharing)
		if err != nil {
			s.logger.Debug().Err(err).Interface("role", invite.GetRoles()[0]).Msg("unable to convert requested role")
			return libregraph.Permission{}, err
		}

		allowedResourceActions := unifiedrole.GetAllowedResourceActions(role, condition)
		if len(allowedResourceActions) == 0 {
			return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, "role not applicable to this resource")
		}

		unifiedRolePermissions = append(unifiedRolePermissions, conversions.ToPointerSlice(role.GetRolePermissions())...)
	}

	driveRecipient := invite.GetRecipients()[0]

	objectID := driveRecipient.GetObjectId()
	cs3ResourcePermissions := unifiedrole.PermissionsToCS3ResourcePermissions(unifiedRolePermissions)

	createShareRequest := &collaboration.CreateShareRequest{
		ResourceInfo: resourceInfo,
		Grant: &collaboration.ShareGrant{
			Permissions: &collaboration.SharePermissions{
				Permissions: cs3ResourcePermissions,
			},
		},
	}

	permission := &libregraph.Permission{}
	if role := unifiedrole.CS3ResourcePermissionsToUnifiedRole(*cs3ResourcePermissions, condition, s.config.FilesSharing.EnableResharing); role != nil {
		permission.Roles = []string{role.GetId()}
	}

	if len(permission.GetRoles()) == 0 {
		permission.LibreGraphPermissionsActions = unifiedrole.CS3ResourcePermissionsToLibregraphActions(*cs3ResourcePermissions)
	}

	switch driveRecipient.GetLibreGraphRecipientType() {
	case "group":
		group, err := s.identityCache.GetGroup(ctx, objectID)
		if err != nil {
			s.logger.Debug().Err(err).Interface("groupId", objectID).Msg("failed group lookup")
			return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, err.Error())
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
		user, err := s.identityCache.GetUser(ctx, objectID)
		if err != nil {
			s.logger.Debug().Err(err).Interface("userId", objectID).Msg("failed user lookup")
			return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, err.Error())
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

	if invite.ExpirationDateTime != nil {
		createShareRequest.GetGrant().Expiration = utils.TimeToTS(*invite.ExpirationDateTime)
	}

	createShareResponse, err := gatewayClient.CreateShare(ctx, createShareRequest)
	if errCode := errorcode.FromCS3Status(createShareResponse.GetStatus(), err); errCode != nil {
		s.logger.Debug().Err(err).Msg("share creation failed")
		return libregraph.Permission{}, *errCode
	}

	if id := createShareResponse.GetShare().GetId().GetOpaqueId(); id != "" {
		permission.Id = conversions.ToPointer(id)
	} else if IsSpaceRoot(resourceInfo.GetId()) {
		// permissions on a space root are not handled by a share manager so
		// they don't get a share-id
		permission.SetId(identitySetToSpacePermissionID(permission.GetGrantedToV2()))
	}

	if expiration := createShareResponse.GetShare().GetExpiration(); expiration != nil {
		permission.SetExpirationDateTime(utils.TSToTime(expiration))
	}

	return *permission, nil
}

// SpaceRootInvite handles invitation request on project spaces
func (s DriveItemPermissionsService) SpaceRootInvite(ctx context.Context, driveID storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return libregraph.Permission{}, err
	}

	space, err := utils.GetSpace(ctx, storagespace.FormatResourceID(driveID), gatewayClient)
	if err != nil {
		return libregraph.Permission{}, err
	}

	if space.SpaceType != "project" {
		return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, "unsupported space type")
	}

	rootResourceID := space.GetRoot()
	return s.Invite(ctx, *rootResourceID, invite)
}

// ListPermissions lists the permissions of a driveItem
func (s DriveItemPermissionsService) ListPermissions(ctx context.Context, itemID storageprovider.ResourceId) (libregraph.CollectionOfPermissionsWithAllowedValues, error) {
	collectionOfPermissions := libregraph.CollectionOfPermissionsWithAllowedValues{}
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return collectionOfPermissions, err
	}

	statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: &itemID}})
	if errCode := errorcode.FromStat(statResponse, err); errCode != nil {
		s.logger.Warn().Err(errCode).Interface("stat.res", statResponse).Msg("stat failed")
		return collectionOfPermissions, err
	}

	condition := unifiedrole.UnifiedRoleConditionGrantee
	if IsSpaceRoot(statResponse.GetInfo().GetId()) {
		condition = unifiedrole.UnifiedRoleConditionOwner
	}

	permissionSet := *statResponse.GetInfo().GetPermissionSet()
	allowedActions := unifiedrole.CS3ResourcePermissionsToLibregraphActions(permissionSet)

	collectionOfPermissions = libregraph.CollectionOfPermissionsWithAllowedValues{
		LibreGraphPermissionsActionsAllowedValues: allowedActions,
		LibreGraphPermissionsRolesAllowedValues: conversions.ToValueSlice(
			unifiedrole.GetApplicableRoleDefinitionsForActions(
				allowedActions,
				condition,
				s.config.FilesSharing.EnableResharing,
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
	if IsSpaceRoot(statResponse.GetInfo().GetId()) {
		permissions, err := s.getSpaceRootPermissions(ctx, statResponse.GetInfo().GetSpace().GetId())
		if err != nil {
			return collectionOfPermissions, err
		}
		collectionOfPermissions.Value = permissions
	} else {
		// "normal" driveItem, populate user  permissions via share providers
		driveItems, err = s.listUserShares(ctx, []*collaboration.Filter{
			share.ResourceIDFilter(conversions.ToPointer(itemID)),
		}, driveItems)
		if err != nil {
			return collectionOfPermissions, err
		}
	}
	// finally get public shares, which are possible for spaceroots and "normal" resources
	driveItems, err = s.listPublicShares(ctx, []*link.ListPublicSharesRequest_Filter{
		publicshare.ResourceIDFilter(conversions.ToPointer(itemID)),
	}, driveItems)
	if err != nil {
		return collectionOfPermissions, err
	}

	for _, driveItem := range driveItems {
		collectionOfPermissions.Value = append(collectionOfPermissions.Value, driveItem.Permissions...)
	}

	return collectionOfPermissions, nil
}

// ListSpaceRootPermissions handles ListPermissions request on project spaces
func (s DriveItemPermissionsService) ListSpaceRootPermissions(ctx context.Context, driveID storageprovider.ResourceId) (libregraph.CollectionOfPermissionsWithAllowedValues, error) {
	collectionOfPermissions := libregraph.CollectionOfPermissionsWithAllowedValues{}
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return collectionOfPermissions, err
	}

	space, err := utils.GetSpace(ctx, storagespace.FormatResourceID(driveID), gatewayClient)
	if err != nil {
		return collectionOfPermissions, err
	}

	if space.SpaceType != "project" {
		return collectionOfPermissions, errorcode.New(errorcode.InvalidRequest, "unsupported space type")
	}

	rootResourceID := space.GetRoot()
	return s.ListPermissions(ctx, *rootResourceID)
}

func (s DriveItemPermissionsService) DeletePermission(ctx context.Context, itemID storageprovider.ResourceId, permissionID string) error {
	var permissionType permissionType

	sharedResourceID, err := s.getLinkPermissionResourceID(ctx, permissionID)
	switch {
	// Check if the ID is referring to a public share
	case err == nil:
		permissionType = Public
	// If the item id is referring to a space root and this is not a public share
	// we have to deal with space permissions
	case IsSpaceRoot(&itemID):
		permissionType = Space
		sharedResourceID = &itemID
		err = nil
	// If this is neither a public share nor a space permission, check if this is a
	// user share
	default:
		sharedResourceID, err = s.getUserPermissionResourceID(ctx, permissionID)
		if err == nil {
			permissionType = User
		}
	}

	switch {
	case err != nil:
		return err
	case permissionType == Unknown:
		return errorcode.New(errorcode.ItemNotFound, "permission not found")
	case sharedResourceID == nil:
		return errorcode.New(errorcode.ItemNotFound, "failed to resolve resource id for shared resource")
	}

	// The resourceID of the shared resource need to match the item ID from the Request Path
	// otherwise this is an invalid Request.
	if !utils.ResourceIDEqual(sharedResourceID, &itemID) {
		s.logger.Debug().Msg("resourceID of shared does not match itemID")
		return errorcode.New(errorcode.InvalidRequest, "permissionID and itemID do not match")
	}

	switch permissionType {
	case User:
		return s.removeUserShare(ctx, permissionID)
	case Public:
		return s.removePublicShare(ctx, permissionID)
	case Space:
		return s.removeSpacePermission(ctx, permissionID, sharedResourceID)
	}

	// This should never be reached
	return errorcode.New(errorcode.GeneralException, "failed to delete permission")
}

func (s DriveItemPermissionsService) UpdatePermission(ctx context.Context, itemID storageprovider.ResourceId, permissionID string, newPermission libregraph.Permission) (libregraph.Permission, error) {
	oldPermission, sharedResourceID, err := s.getPermissionByID(ctx, permissionID, &itemID)
	if err != nil {
		return libregraph.Permission{}, err
	}

	// The resourceID of the shared resource need to match the item ID from the Request Path
	// otherwise this is an invalid Request.
	if !utils.ResourceIDEqual(sharedResourceID, &itemID) {
		s.logger.Debug().Msg("resourceID of shared does not match itemID")
		return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, "permissionID and itemID do not match")
	}

	// This is a public link
	if _, ok := oldPermission.GetLinkOk(); ok {
		updatedPermission, err := s.updatePublicLinkPermission(ctx, permissionID, &itemID, &newPermission)
		if err != nil {
			return libregraph.Permission{}, err
		}
		return *updatedPermission, nil
	}

	// This is a user share
	updatedPermission, err := s.updateUserShare(ctx, permissionID, sharedResourceID, &newPermission)
	if err != nil {
		return libregraph.Permission{}, err
	}
	return *updatedPermission, nil
}

// DriveItemPermissionsService is the api that registers the http endpoints which expose needed operation to the graph api.
// the business logic is delegated to the permissions service and further down to the cs3 client.
type DriveItemPermissionsApi struct {
	logger                      log.Logger
	driveItemPermissionsService DriveItemPermissionsProvider
}

// NewDriveItemPermissionsApi creates a new DriveItemPermissionsApi
func NewDriveItemPermissionsApi(driveItemPermissionService DriveItemPermissionsProvider, logger log.Logger) (DriveItemPermissionsApi, error) {
	return DriveItemPermissionsApi{
		logger:                      log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemApi").Logger()},
		driveItemPermissionsService: driveItemPermissionService,
	}, nil
}

func (api DriveItemPermissionsApi) Invite(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Debug().Err(err).Msg(invalidIdMsg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, invalidIdMsg)
		return
	}

	driveItemInvite := &libregraph.DriveItemInvite{}
	if err = StrictJSONUnmarshal(r.Body, driveItemInvite); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()

	if err = validate.StructCtx(ctx, driveItemInvite); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	permission, err := api.driveItemPermissionsService.Invite(ctx, itemID, *driveItemInvite)

	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: []interface{}{permission}})
}

func (api DriveItemPermissionsApi) SpaceRootInvite(w http.ResponseWriter, r *http.Request) {
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		msg := "could not parse driveID"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	driveItemInvite := &libregraph.DriveItemInvite{}
	if err = StrictJSONUnmarshal(r.Body, driveItemInvite); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	if err = validate.StructCtx(ctx, driveItemInvite); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	permission, err := api.driveItemPermissionsService.SpaceRootInvite(ctx, driveID, *driveItemInvite)

	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: []interface{}{permission}})
}

func (api DriveItemPermissionsApi) ListPermissions(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Debug().Err(err).Msg(invalidIdMsg)
		errorcode.RenderError(w, r, err)
		return
	}

	ctx := r.Context()

	permissions, err := api.driveItemPermissionsService.ListPermissions(ctx, itemID)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, permissions)
}

func (api DriveItemPermissionsApi) ListSpaceRootPermissions(w http.ResponseWriter, r *http.Request) {
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		msg := "could not parse driveID"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	ctx := r.Context()
	permissions, err := api.driveItemPermissionsService.ListSpaceRootPermissions(ctx, driveID)

	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, permissions)
}

func (api DriveItemPermissionsApi) DeletePermission(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Debug().Err(err).Msg(invalidIdMsg)
		errorcode.RenderError(w, r, err)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Debug().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	ctx := r.Context()
	err = api.driveItemPermissionsService.DeletePermission(ctx, itemID, permissionID)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func (api DriveItemPermissionsApi) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Debug().Err(err).Msg(invalidIdMsg)
		errorcode.RenderError(w, r, err)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Debug().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	permission := libregraph.Permission{}
	if err = StrictJSONUnmarshal(r.Body, &permission); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	if err = validate.StructCtx(ctx, permission); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	updatedPermission, err := api.driveItemPermissionsService.UpdatePermission(ctx, itemID, permissionID, permission)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &updatedPermission)
}
