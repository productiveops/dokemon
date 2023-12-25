import { cn } from "@/lib/utils"
import { Link, NavLink, useRoutes } from "react-router-dom"
import { SideNavTopLevel } from "./side-nav-top-level"
import { SideNavNode } from "./side-nav-node"
import { SideNavContainer } from "./side-nav-container"
import { SideNavCompose } from "./side-nav-compose"
import { ArrowLeftIcon } from "@heroicons/react/24/outline"

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

export function SideNavBack({ to }: { to: string }) {
  return (
    <div className="mb-1">
      <Link
        to={to}
        className="group flex gap-x-4 rounded-md p-2 text-sm font-semibold leading-6 text-gray-300 underline-offset-4 hover:underline"
      >
        <ArrowLeftIcon className="-ml-1 mt-1 h-4 w-4" aria-hidden="true" /> Back
      </Link>
    </div>
  )
}
