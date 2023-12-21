import { cn } from "@/lib/utils"
import {
  GlobeAltIcon,
  Bars3BottomLeftIcon,
  Bars3Icon,
  CommandLineIcon,
  CubeIcon,
  CubeTransparentIcon,
  FolderIcon,
  ArrowLeftIcon,
  PlayIcon,
  ComputerDesktopIcon,
  SunIcon,
  VariableIcon,
  LockClosedIcon,
} from "@heroicons/react/24/outline"
import { Link, NavLink, useParams, useRoutes } from "react-router-dom"

export default function SideNav() {
  const element = useRoutes([
    {
      path: "nodes/:nodeId/containers/:containerId/*",
      element: <SideNavContainer />,
    },
    {
      path: "compose/create",
      element: <SideNavTopLevel />,
    },
    {
      path: "nodes/:nodeId/compose/create/*",
      element: <SideNavNode />,
    },
    {
      path: "nodes/:nodeId/compose/:composeProjectId/*",
      element: <SideNavCompose />,
    },
    {
      path: "nodes/:nodeId/*",
      element: <SideNavNode />,
    },
    {
      path: "*",
      element: <SideNavTopLevel />,
    },
  ])

  return element
}

function SideNavTopLevel() {
  return (
    <>
      <ul role="list" className="-mx-2 space-y-1">
        <li>
          <SideBarItem to={`/nodes`}>
            <ComputerDesktopIcon
              className="h-6 w-6 shrink-0"
              aria-hidden="true"
            />
            Nodes
          </SideBarItem>
          <SideBarItem to={`/composelibrary`}>
            <Bars3BottomLeftIcon
              className="h-6 w-6 shrink-0"
              aria-hidden="true"
            />
            Compose Library
          </SideBarItem>
          <SideBarItem to={`/environments`}>
            <SunIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Environments
          </SideBarItem>
          <SideBarItem to={`/variables`}>
            <VariableIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Variables
          </SideBarItem>
          <SideBarItem to={`/credentials`}>
            <LockClosedIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Credentials
          </SideBarItem>
        </li>
      </ul>
    </>
  )
}

function SideNavNode() {
  const { nodeId } = useParams()

  return (
    <>
      <div className="mb-1">
        <Link
          to={`/nodes`}
          className="group flex gap-x-4 rounded-md p-2 text-sm font-semibold leading-6 text-gray-300 underline-offset-4 hover:underline"
        >
          <ArrowLeftIcon className="-ml-1 mt-1 h-4 w-4" aria-hidden="true" />{" "}
          Back
        </Link>
      </div>
      <ul role="list" className="-mx-2 space-y-1">
        <li>
          <SideBarItem to={`/nodes/${nodeId}/compose`}>
            <Bars3BottomLeftIcon
              className="h-6 w-6 shrink-0"
              aria-hidden="true"
            />
            Compose
          </SideBarItem>
        </li>

        <li>
          <SideBarItem to={`/nodes/${nodeId}/containers`}>
            <CubeIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Containers
          </SideBarItem>
        </li>
        <li>
          <SideBarItem to={`/nodes/${nodeId}/images`}>
            <CubeTransparentIcon
              className="h-6 w-6 shrink-0"
              aria-hidden="true"
            />
            Images
          </SideBarItem>
        </li>
        <li>
          <SideBarItem to={`/nodes/${nodeId}/volumes`}>
            <FolderIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Volumes
          </SideBarItem>
        </li>
        <li>
          <SideBarItem to={`/nodes/${nodeId}/networks`}>
            <GlobeAltIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Networks
          </SideBarItem>
        </li>
        <li>
          <SideBarItem to={`/nodes/${nodeId}/details`}>
            <ComputerDesktopIcon
              className="h-6 w-6 shrink-0"
              aria-hidden="true"
            />
            Node Details
          </SideBarItem>
        </li>
      </ul>
    </>
  )
}

function SideNavContainer() {
  const { nodeId, containerId } = useParams()

  return (
    <>
      <div className="mb-1">
        <Link
          to={`/nodes/${nodeId}/containers`}
          className="group flex gap-x-4 rounded-md p-2 text-sm font-semibold leading-6 text-gray-300 underline-offset-4 hover:underline"
        >
          <ArrowLeftIcon className="-ml-1 mt-1 h-4 w-4" aria-hidden="true" />{" "}
          Back
        </Link>
      </div>
      <ul role="list" className="-mx-2 space-y-1">
        <li>
          <SideBarItem to={`/nodes/${nodeId}/containers/${containerId}/logs`}>
            <Bars3Icon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Logs
          </SideBarItem>
        </li>
        <li>
          <SideBarItem
            to={`/nodes/${nodeId}/containers/${containerId}/terminal`}
          >
            <CommandLineIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Terminal
          </SideBarItem>
        </li>
      </ul>
    </>
  )
}

function SideNavCompose() {
  const { nodeId, composeProjectId } = useParams()

  return (
    <>
      <div className="mb-1">
        <Link
          to={`/nodes/${nodeId}/compose`}
          className="group flex gap-x-4 rounded-md p-2 text-sm font-semibold leading-6 text-gray-300 underline-offset-4 hover:underline"
        >
          <ArrowLeftIcon className="-ml-1 mt-1 h-4 w-4" aria-hidden="true" />{" "}
          Back
        </Link>
      </div>
      <ul role="list" className="-mx-2 space-y-1">
        <li>
          <SideBarItem
            to={`/nodes/${nodeId}/compose/${composeProjectId}/actions`}
          >
            <PlayIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Actions
          </SideBarItem>
        </li>
        <li>
          <SideBarItem
            to={`/nodes/${nodeId}/compose/${composeProjectId}/definition`}
          >
            <Bars3BottomLeftIcon
              className="h-6 w-6 shrink-0"
              aria-hidden="true"
            />
            Definition
          </SideBarItem>
        </li>
        <li>
          <SideBarItem
            to={`/nodes/${nodeId}/compose/${composeProjectId}/containers`}
          >
            <CubeIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Containers
          </SideBarItem>
        </li>
        <li>
          <SideBarItem to={`/nodes/${nodeId}/compose/${composeProjectId}/logs`}>
            <Bars3Icon className="h-6 w-6 shrink-0" aria-hidden="true" />
            Logs
          </SideBarItem>
        </li>
      </ul>
    </>
  )
}

export function SideBarItem({ to, children }: { to: string; children: any }) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        cn(
          "group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6",
          isActive
            ? "bg-gray-800 text-white"
            : "text-gray-400 hover:bg-gray-800 hover:text-white"
        )
      }
    >
      {children}
    </NavLink>
  )
}
