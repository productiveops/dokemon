import {
  Bars3BottomLeftIcon,
  Bars3Icon,
  CubeIcon,
} from "@heroicons/react/24/outline"
import { useParams } from "react-router-dom"
import { SideBarItem, SideNavBack } from "./side-nav"

export function SideNavCompose() {
  const { nodeId, composeProjectId } = useParams()

  const baseUrl = `/nodes/${nodeId}/compose/${composeProjectId}`
  const items = [
    {
      title: "Definition",
      link: `${baseUrl}/definition`,
      icon: (
        <Bars3BottomLeftIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
      ),
    },
    {
      title: "Containers",
      link: `${baseUrl}/containers`,
      icon: <CubeIcon className="h-6 w-6 shrink-0" aria-hidden="true" />,
    },
    {
      title: "Logs",
      link: `${baseUrl}/logs`,
      icon: <Bars3Icon className="h-6 w-6 shrink-0" aria-hidden="true" />,
    },
  ]

  return (
    <>
      <SideNavBack to={`/nodes/${nodeId}/compose`} />
      <ul role="list" className="-mx-2 space-y-1">
        {items.map((item) => (
          <li key={item.title}>
            <SideBarItem to={item.link}>
              {item.icon}
              {item.title}
            </SideBarItem>
          </li>
        ))}
      </ul>
    </>
  )
}
