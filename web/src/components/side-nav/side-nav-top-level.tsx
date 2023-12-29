import {
  Bars3BottomLeftIcon,
  ComputerDesktopIcon,
  LockClosedIcon,
  SunIcon,
  VariableIcon,
} from "@heroicons/react/24/outline"
import { SideBarItem } from "./side-nav"

const items = [
  {
    title: "Nodes",
    link: "/nodes",
    icon: (
      <ComputerDesktopIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
    ),
  },
  {
    title: "Compose Library",
    link: "/composelibrary",
    icon: (
      <Bars3BottomLeftIcon className="h-6 w-6 shrink-0" aria-hidden="true" />
    ),
  },
  {
    title: "Environments",
    link: "/environments",
    icon: <SunIcon className="h-6 w-6 shrink-0" aria-hidden="true" />,
  },
  {
    title: "Variables",
    link: "/variables",
    icon: <VariableIcon className="h-6 w-6 shrink-0" aria-hidden="true" />,
  },
  {
    title: "Credentials",
    link: "/credentials",
    icon: <LockClosedIcon className="h-6 w-6 shrink-0" aria-hidden="true" />,
  },
]

export function SideNavTopLevel() {
  return (
    <>
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
