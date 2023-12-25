import { Bars3Icon, CommandLineIcon } from "@heroicons/react/24/outline"
import { useParams } from "react-router-dom"
import { SideBarItem, SideNavBack } from "./side-nav"

export function SideNavContainer() {
  const { nodeId, containerId } = useParams()

  return (
    <>
      <SideNavBack to={`/nodes/${nodeId}/containers`} />
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
