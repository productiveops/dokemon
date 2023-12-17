import { createBrowserRouter, Navigate } from "react-router-dom"
import "./index.css"
import Root from "@/app/root"
import ErrorPage from "./error-page"
import Containers from "./app/containers/containers"
import Images from "./app/images/images"
import ContainerLogs from "./app/containers/container/logs"
import ContainerTerminal from "./app/containers/container/terminal"
import Volumes from "./app/volumes/volumes"
import Networks from "./app/networks/networks"
import ChangePassword from "./app/auth/change-password"
import Login from "./app/auth/login"
import Setup from "./app/auth/setup"
import ComposeLibraryItems from "./app/compose-library/compose-library-items"
import ComposeContainers from "./app/compose/compose/containers"
import ComposeLogs from "./app/compose/compose/logs"
import ComposeActions from "./app/compose/compose/actions"
import CreateComposeProject from "./app/compose-library/create-compose-project"
import EditComposeProject from "./app/compose-library/edit-compose-project"
import NodeCompose from "./app/compose/node-compose"
import Nodes from "./app/nodes/nodes"
import NodeDetails from "./app/nodes/node-details"
import Environments from "./app/environments/environments"
import Variables from "./app/variables/variables"

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
        element: <Nodes />,
      },
      {
        path: "nodes/:nodeId/containers",
        element: <Containers />,
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
        element: <Images />,
      },
      {
        path: "nodes/:nodeId/volumes",
        element: <Volumes />,
      },
      {
        path: "nodes/:nodeId/networks",
        element: <Networks />,
      },
      {
        path: "nodes/:nodeId/compose",
        element: <NodeCompose />,
      },
      {
        path: "nodes/:nodeId/details",
        element: <NodeDetails />,
      },
      {
        path: "composelibrary",
        element: <ComposeLibraryItems />,
      },
      {
        path: "composelibrary/create",
        element: <CreateComposeProject />,
      },
      {
        path: "composelibrary/:composeProjectName/edit",
        element: <EditComposeProject />,
      },
      {
        path: "environments",
        element: <Environments />,
      },
      {
        path: "variables",
        element: <Variables />,
      },
      {
        path: "nodes/:nodeId/compose/:composeProjectId/actions",
        element: <ComposeActions />,
      },
      {
        path: "nodes/:nodeId/compose/:composeProjectId/containers",
        element: <ComposeContainers />,
      },
      {
        path: "nodes/:nodeId/compose/:composeProjectId/logs",
        element: <ComposeLogs />,
      },
    ],
  },
])

export default router
