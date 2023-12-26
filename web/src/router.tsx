import { createBrowserRouter, Navigate } from "react-router-dom"
import "./index.css"
import Root from "@/app/root"
import ErrorPage from "./error-page"
import ContainerList from "./app/containers/container-list"
import ImageList from "./app/images/image-list"
import ContainerLogs from "./app/containers/container-logs"
import ContainerTerminal from "./app/containers/container-terminal"
import VolumeList from "./app/volumes/volume-list"
import NetworkList from "./app/networks/network-list"
import ChangePassword from "./app/auth/change-password"
import Login from "./app/auth/login"
import Setup from "./app/auth/setup"
import ComposeLibraryList from "./app/compose-library/compose-library-list"
import ComposeContainerList from "./app/compose/compose-container-list"
import ComposeLogs from "./app/compose/compose-logs"
import ComposeLibraryCreateFileSystemProject from "./app/compose-library/compose-library-create-filesystem-project"
import ComposeLibraryEditFileSystemProject from "./app/compose-library/compose-library-edit-filesystem-project"
import ComposeList from "./app/compose/compose-list"
import NodeList from "./app/nodes/node-list"
import NodeDetails from "./app/nodes/node-details"
import EnvironmentList from "./app/environments/environment-list"
import VariableList from "./app/variables/variable-list"
import ComposeLibraryCreateGitHubProject from "./app/compose-library/compose-library-create-github-project"
import CredentialList from "./app/credentials/credential-list"
import ComposeLibraryEditGitHubProject from "./app/compose-library/compose-library-edit-github-project"
import ComposeAddGitHub from "./app/compose/compose-add-github"
import ComposeDefinition from "./app/compose/compose-definition"
import ComposeAddLocal from "./app/compose/compose-add-local"

const router = createBrowserRouter([
  {
    path: "/",
    element: <Navigate to="/nodes" />,
  },
  {
    path: "/changepassword",
    element: <ChangePassword />,
  },
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/setup",
    element: <Setup />,
  },
  {
    path: "/*",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "nodes",
        element: <NodeList />,
      },
      {
        path: "nodes/:nodeId/containers",
        element: <ContainerList />,
      },
      {
        path: "nodes/:nodeId/containers/:containerId/logs",
        element: <ContainerLogs />,
      },
      {
        path: "nodes/:nodeId/containers/:containerId/terminal",
        element: <ContainerTerminal />,
      },
      {
        path: "nodes/:nodeId/images",
        element: <ImageList />,
      },
      {
        path: "nodes/:nodeId/volumes",
        element: <VolumeList />,
      },
      {
        path: "nodes/:nodeId/networks",
        element: <NetworkList />,
      },
      {
        path: "nodes/:nodeId/details",
        element: <NodeDetails />,
      },
      {
        path: "nodes/:nodeId/compose",
        element: <ComposeList />,
      },
      {
        path: "nodes/:nodeId/compose/create/github",
        element: <ComposeAddGitHub />,
      },
      {
        path: "nodes/:nodeId/compose/create/local",
        element: <ComposeAddLocal />,
      },
      {
        path: "nodes/:nodeId/compose/:composeProjectId/definition",
        element: <ComposeDefinition />,
      },
      {
        path: "nodes/:nodeId/compose/:composeProjectId/containers",
        element: <ComposeContainerList />,
      },
      {
        path: "nodes/:nodeId/compose/:composeProjectId/logs",
        element: <ComposeLogs />,
      },
      {
        path: "composelibrary",
        element: <ComposeLibraryList />,
      },
      {
        path: "composelibrary/filesystem/create",
        element: <ComposeLibraryCreateFileSystemProject />,
      },
      {
        path: "composelibrary/github/create",
        element: <ComposeLibraryCreateGitHubProject />,
      },
      {
        path: "composelibrary/filesystem/:composeProjectName/edit",
        element: <ComposeLibraryEditFileSystemProject />,
      },
      {
        path: "composelibrary/github/:composeProjectId/edit",
        element: <ComposeLibraryEditGitHubProject />,
      },
      {
        path: "environments",
        element: <EnvironmentList />,
      },
      {
        path: "variables",
        element: <VariableList />,
      },
      {
        path: "credentials",
        element: <CredentialList />,
      },
    ],
  },
])

export default router
