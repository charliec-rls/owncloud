package svc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

// GetRoleDefinitions a list of permission roles than can be used when sharing with users or groups
func (g Graph) GetRoleDefinitions(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, getRoleDefinitionList(g.config.FilesSharing.EnableResharing))
}

// GetRoleDefinition a permission role than can be used when sharing with users or groups
func (g Graph) GetRoleDefinition(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	roleID, err := url.PathUnescape(chi.URLParam(r, "roleID"))
	if err != nil {
		logger.Debug().Err(err).Str("roleID", chi.URLParam(r, "roleID")).Msg("could not get roleID: unescaping is failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping role id failed")
		return
	}
	role, err := getRoleDefinition(roleID, g.config.FilesSharing.EnableResharing)
	if err != nil {
		logger.Debug().Str("roleID", roleID).Msg("could not get role: not found")
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, role)
}

func getRoleDefinitionList(resharing bool) []*libregraph.UnifiedRoleDefinition {
	return []*libregraph.UnifiedRoleDefinition{
		unifiedrole.NewViewerUnifiedRole(resharing),
		unifiedrole.NewSpaceViewerUnifiedRole(),
		unifiedrole.NewEditorUnifiedRole(resharing),
		unifiedrole.NewSpaceEditorUnifiedRole(),
		unifiedrole.NewFileEditorUnifiedRole(resharing),
		unifiedrole.NewCoownerUnifiedRole(),
		unifiedrole.NewUploaderUnifiedRole(),
		unifiedrole.NewManagerUnifiedRole(),
	}
}

func getRoleDefinition(roleID string, resharing bool) (*libregraph.UnifiedRoleDefinition, error) {
	roleList := getRoleDefinitionList(resharing)
	for _, role := range roleList {
		if role != nil && role.Id != nil && *role.Id == roleID {
			return role, nil
		}
	}
	return nil, fmt.Errorf("role not found")
}
