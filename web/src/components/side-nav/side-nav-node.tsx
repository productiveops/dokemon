import {
  Bars3BottomLeftIcon,
  ComputerDesktopIcon,
  CubeIcon,
  CubeTransparentIcon,
  FolderIcon,
  GlobeAltIcon,
} from "@heroicons/react/24/outline"
import { useParams } from "react-router-dom"
import { SideBarItem, SideNavBack } from "./side-nav"

export function SideNavNode() {
  const { nodeId } = useParams()

  const baseUrl = `/nodes/${nodeId}`
  const items = [
    {
      title: "Compose",
      link: `${baseUrl}/compose`,
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
      title: "Images",
      link: `${baseUrl}/images`,
      icon: (
        <CubeTransparentIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
      ),
    },
    {
      title: "Volumes",
      link: `${baseUrl}/volumes`,
      icon: <FolderIcon className="h-6 w-6 shrink-0" aria-hidden="true" />,
    },
    {
      title: "Networks",
      link: `${baseUrl}/networks`,
      icon: <GlobeAltIcon className="h-6 w-6 shrink-0" aria-hidden="true" />,
    },
    {
      title: "Node Details",
      link: `${baseUrl}/details`,
      icon: (
        <ComputerDesktopIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
      ),
    },
  ]

  return (
    <>
      <SideNavBack to="/nodes" />
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
