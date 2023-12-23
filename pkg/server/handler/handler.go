package handler

import (
	"net/http"
	"time"

	"github.com/productiveops/dokemon/web"

	"github.com/productiveops/dokemon/pkg/server/store"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	composeProjectsPath string
	composeLibraryStore store.ComposeLibraryStore
	credentialStore store.CredentialStore
	environmentStore store.EnvironmentStore
	userStore store.UserStore
	nodeStore store.NodeStore
	nodeComposeProjectStore store.NodeComposeProjectStore
	nodeComposeProjectVariableStore store.NodeComposeProjectVariableStore
	settingStore store.SettingStore
	variableStore store.VariableStore
	variableValueStore store.VariableValueStore
	fileSystemComposeLibraryStore store.FileSystemComposeLibraryStore
}

var (
	defaultTimeout = 30 * time.Second
)

func NewHandler(
	composeProjectsPath string,
	composeLibraryStore store.ComposeLibraryStore,
	credentialStore store.CredentialStore,
	environmentStore store.EnvironmentStore,
	userStore store.UserStore,
	nodeStore store.NodeStore,
	nodeComposeProjectStore store.NodeComposeProjectStore,
	nodeComposeProjectVariableStore store.NodeComposeProjectVariableStore,
	settingStore store.SettingStore,
	variableStore store.VariableStore,
	variableValueStore store.VariableValueStore,
	fileSystemComposeLibraryStore store.FileSystemComposeLibraryStore,
	) *Handler {
		return &Handler{
		composeProjectsPath: composeProjectsPath,
		composeLibraryStore: composeLibraryStore,
		credentialStore: credentialStore,
		environmentStore: environmentStore,
		userStore: userStore,
		nodeStore: nodeStore,
		nodeComposeProjectStore: nodeComposeProjectStore,
		nodeComposeProjectVariableStore: nodeComposeProjectVariableStore,
		settingStore: settingStore,
		variableStore: variableStore,
		variableValueStore: variableValueStore,
		fileSystemComposeLibraryStore: fileSystemComposeLibraryStore,
		}
}

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/ws", h.HandleWebSocket)	// This won't have versioning as we can't expect customers to update agents for new versions

	v1 := e.Group("/api/v1")

	github := v1.Group("/github/filecontent/load")
	github.POST("", h.RetrieveGitHubFileContent)

	settings := v1.Group("/settings")
	settings.GET("/:id", h.GetSettingById)
	settings.PUT("/:id", h.UpdateSetting)

	credentials := v1.Group("/credentials")
	credentials.POST("", h.CreateCredential)
	credentials.PUT("/:id", h.UpdateCredentialDetails)
	credentials.PUT("/:id/secret", h.UpdateCredentialSecret)
	credentials.GET("", h.GetCredentialList)
	credentials.GET("/:id", h.GetCredentialById)
	credentials.DELETE("/:id", h.DeleteCredentialById)
	credentials.GET("/uniquename", h.IsUniqueCredentialName)
	credentials.GET("/:id/uniquename", h.IsUniqueCredentialNameExcludeItself)

	environments := v1.Group("/environments")
	environments.POST("", h.CreateEnvironment)
	environments.PUT("/:id", h.UpdateEnvironment)
	environments.GET("", h.GetEnvironmentList)
	environments.GET("/map", h.GetEnvironmentMap)
	environments.GET("/:id", h.GetEnvironmentById)
	environments.DELETE("/:id", h.DeleteEnvironmentById)
	environments.GET("/uniquename", h.IsUniqueEnvironmentName)
	environments.GET("/:id/uniquename", h.IsUniqueEnvironmentNameExcludeItself)

	nodes := v1.Group("/nodes")
	nodes.POST("", h.CreateNode)
	nodes.PUT("/:id", h.UpdateNode)
	nodes.PATCH("/:id", h.UpdateNodeContainerBaseUrl)
	nodes.GET("", h.GetNodeList)
	nodes.GET("/:id", h.GetNodeById)
	nodes.GET("/head/:id", h.GetNodeHeadById)
	nodes.DELETE("/:id", h.DeleteNodeById)
	nodes.GET("/uniquename", h.IsUniqueNodeName)
	nodes.GET("/:id/uniquename", h.IsUniqueNodeNameExcludeItself)
	nodes.POST("/:id/generatetoken", h.GenerateRegistrationToken)

	containers := nodes.Group("/:nodeId/containers")
	containers.GET("", h.GetContainerList)
	containers.POST("/start", h.StartContainer)
	containers.POST("/stop", h.StopContainer)
	containers.POST("/restart", h.RestartContainer)
	containers.POST("/remove", h.RemoveContainer)
	containers.GET("/:id/logs", h.ViewContainerLogs)
	containers.GET("/:id/terminal", h.OpenContainerTerminal)

	images := nodes.Group("/:nodeId/images")
	images.GET("", h.GetImageList)
	images.POST("/remove", h.RemoveImage)
	images.POST("/prune", h.PruneImages)

	volumes := nodes.Group("/:nodeId/volumes")
	volumes.GET("", h.GetVolumeList)
	volumes.POST("/remove", h.RemoveVolume)
	volumes.POST("/prune", h.PruneVolumes)

	networks := nodes.Group("/:nodeId/networks")
	networks.GET("", h.GetNetworkList)
	networks.POST("/remove", h.RemoveNetwork)
	networks.POST("/prune", h.PruneNetworks)

	composelibrary := v1.Group("/composelibrary")
	composelibrary.GET("", h.GetComposeProjectList)
	composelibrary.GET("/uniquename", h.IsUniqueComposeProjectName)
	composelibrary.GET("/uniquenameexcludeitself", h.IsUniqueComposeProjectNameExcludeItself)

	filesystemcomposelibrary := composelibrary.Group("/filesystem")
	filesystemcomposelibrary.POST("", h.CreateFileSystemComposeProject)
	filesystemcomposelibrary.PUT("/:projectName", h.UpdateFileSystemComposeProject)
	filesystemcomposelibrary.DELETE("/:projectName", h.DeleteFileSystemComposeProject)
	filesystemcomposelibrary.GET("/:projectName", h.GetFileSystemComposeProject)

	githubcomposelibrary := composelibrary.Group("/github")
	githubcomposelibrary.POST("", h.CreateGitHubComposeProject)
	githubcomposelibrary.PUT("/:id", h.UpdateGitHubComposeProject)
	githubcomposelibrary.DELETE("/:id", h.DeleteGitHubComposeProject)
	githubcomposelibrary.GET("/:id", h.GetGitHubComposeProjectById)

	node_compose_project := nodes.Group("/:nodeId/compose")
	node_compose_project.GET("", h.GetNodeComposeProjectList)
	node_compose_project.POST("/create/github", h.CreateGitHubNodeComposeProject)
	node_compose_project.POST("/create/local", h.CreateLocalNodeComposeProject)
	node_compose_project.POST("/create/library", h.AddNodeComposeProjectFromLibrary)
	node_compose_project.GET("/uniquename", h.IsUniqueNodeComposeProjectName)
	node_compose_project.GET("/:id/uniquename", h.IsUniqueNodeComposeProjectNameExcludeItself)
	node_compose_project.PUT("/:id/github", h.UpdateGitHubNodeComposeProject)
	node_compose_project.PUT("/:id/local", h.UpdateLocalNodeComposeProject)
	node_compose_project.GET("/:id", h.GetNodeComposeProject)
	node_compose_project.DELETE("/:id", h.DeleteNodeComposeProject)
	node_compose_project.GET("/:id/containers", h.GetNodeComposeContainerList)
	node_compose_project.GET("/:id/logs", h.GetNodeComposeLogs)
	node_compose_project.GET("/:id/pull", h.GetNodeComposePull)
	node_compose_project.GET("/:id/up", h.GetNodeComposeUp)
	node_compose_project.GET("/:id/down", h.GetNodeComposeDown)

	node_compose_project_variables := node_compose_project.Group("/:node_compose_project_id/variables")
	node_compose_project_variables.POST("", h.CreateVariable)
	node_compose_project_variables.PUT("/:id", h.UpdateVariable)
	node_compose_project_variables.GET("", h.GetVariableList)
	node_compose_project_variables.GET("/:id", h.GetVariableById)
	node_compose_project_variables.DELETE("/:id", h.DeleteVariableById)
	node_compose_project_variables.GET("/uniquename", h.IsUniqueVariableName)
	node_compose_project_variables.GET("/:id/uniquename", h.IsUniqueVariableNameExcludeItself)

	variables := v1.Group("/variables")
	variables.POST("", h.CreateVariable)
	variables.PUT("/:id", h.UpdateVariable)
	variables.GET("", h.GetVariableList)
	variables.GET("/:id", h.GetVariableById)
	variables.DELETE("/:id", h.DeleteVariableById)
	variables.GET("/uniquename", h.IsUniqueVariableName)
	variables.GET("/:id/uniquename", h.IsUniqueVariableNameExcludeItself)

	variablevalues := variables.Group("/:variableId/values/:environmentId")
	variablevalues.PUT("", h.CreateOrUpdateVariableValue)
	variablevalues.GET("", h.GetVariableValue)

	changepassword := v1.Group("/changepassword")
	changepassword.POST("", h.ChangeUserPassword)

	// No auth required
	users := v1.Group("/users")
	users.POST("", h.CreateUser)
	users.POST("/login", h.UserLogin)
	users.POST("/logout", h.UserLogout)
	users.POST("/count", h.UserCount)

	// React App
	web.RegisterHandlers(e)
}

