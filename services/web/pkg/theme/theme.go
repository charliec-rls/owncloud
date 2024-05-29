package theme

import (
	"path"

	"github.com/owncloud/ocis/v2/ocis-pkg/capabilities"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

var (
	_brandingRoot  = "_branding"
	_themeFileName = "theme.json"
)

// themeDefaults contains the default values for the theme.
// all rendered themes get the default values from here.
var themeDefaults = KV{
	"common": KV{
		"shareRoles": KV{
			unifiedrole.UnifiedRoleViewerID: KV{
				"name":     "UnifiedRoleViewer",
				"iconName": "eye",
			},
			unifiedrole.UnifiedRoleSpaceViewerID: KV{
				"label":    "UnifiedRoleSpaceViewer",
				"iconName": "eye",
			},
			unifiedrole.UnifiedRoleFileEditorID: KV{
				"label":    "UnifiedRoleFileEditor",
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleEditorID: KV{
				"label":    "UnifiedRoleEditor",
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleSpaceEditorID: KV{
				"label":    "UnifiedRoleSpaceEditor",
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleManagerID: KV{
				"label":    "UnifiedRoleManager",
				"iconName": "user-star",
			},
			unifiedrole.UnifiedRoleEditorLiteID: KV{
				"label":    "UnifiedRoleEditorLite",
				"iconName": "upload",
			},
			unifiedrole.UnifiedRoleSecureViewerID: KV{
				"label":    "UnifiedRoleSecureView",
				"iconName": "shield",
			},
		},
	},
}

// isFiletypePermitted checks if the given file extension is allowed.
func isFiletypePermitted(filename string, givenMime string) bool {
	// Check if we allow that extension and if the mediatype matches the extension
	extensionMime, ok := capabilities.Default().Theme.Logo.PermittedFileTypes[path.Ext(filename)]
	return ok && extensionMime == givenMime
}
