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
	credentialStore store.CredentialStore
	environmentStore store.EnvironmentStore
	userStore store.UserStore
	nodeStore store.NodeStore
	nodeComposeProjectStore store.NodeComposeProjectStore
	settingStore store.SettingStore
	variableStore store.VariableStore
	variableValueStore store.VariableValueStore
	composeLibraryStore store.ComposeLibraryStore
}

var (
	defaultTimeout = 30 * time.Second
)

func NewHandler(
	composeProjectsPath string,
	credentialStore store.CredentialStore,
	environmentStore store.EnvironmentStore,
	userStore store.UserStore,
	nodeStore store.NodeStore,
	nodeComposeProjectStore store.NodeComposeProjectStore,
	settingStore store.SettingStore,
	variableStore store.VariableStore,
	variableValueStore store.VariableValueStore,
	composeLibraryStore store.ComposeLibraryStore,
	) *Handler {
		return &Handler{
		composeProjectsPath: composeProjectsPath,
		credentialStore: credentialStore,
		environmentStore: environmentStore,
		userStore: userStore,
		nodeStore: nodeStore,
		nodeComposeProjectStore: nodeComposeProjectStore,
		settingStore: settingStore,
		variableStore: variableStore,
		variableValueStore: variableValueStore,
		composeLibraryStore: composeLibraryStore,
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
	composelibrary.POST("", h.CreateComposeProject)
	composelibrary.PUT("/:projectName", h.UpdateComposeProject)
	composelibrary.DELETE("/:projectName", h.DeleteComposeProject)
	composelibrary.GET("/:projectName", h.GetComposeProject)

	node_compose := nodes.Group("/:nodeId/compose")
	node_compose.GET("", h.GetNodeComposeProjectList)
	node_compose.POST("", h.CreateNodeComposeProject)
	node_compose.GET("/uniquename", h.IsUniqueNodeComposeProjectName)
	node_compose.GET("/:id", h.GetNodeComposeProject)
	node_compose.DELETE("/:id", h.DeleteNodeComposeProject)
	node_compose.GET("/:id/containers", h.GetNodeComposeContainerList)
	node_compose.GET("/:id/logs", h.GetNodeComposeLogs)
	node_compose.GET("/:id/pull", h.GetNodeComposePull)
	node_compose.GET("/:id/up", h.GetNodeComposeUp)
	node_compose.GET("/:id/down", h.GetNodeComposeDown)

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

