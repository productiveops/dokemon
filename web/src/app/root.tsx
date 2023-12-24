import { Outlet } from "react-router-dom"
import SideNav from "../components/side-nav/side-nav"
import { VERSION } from "@/lib/version"
import { ModeToggle } from "@/components/ui/mode-toggle"
import { SideNavLeftBottom } from "@/components/side-nav/side-nav-left-bottom"

export default function Root() {
  return (
    <>
      <div>
        <div className="fixed inset-y-0 z-50 flex w-72 flex-col">
          <div className="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 pb-0 dark:bg-gray-950">
            <div className="flex h-16 shrink-0 items-center text-white">
              <img
                className="ml-9 w-24"
                src="/assets/images/dokemon-dark-small.svg"
                alt="Dokémon"
              />
              <span className="ml-3 mr-5 pt-[3px] text-sm">v{VERSION}</span>
              <ModeToggle />
            </div>
            <nav className="flex flex-1 flex-col">
              <ul role="list" className="flex flex-1 flex-col gap-y-7">
                <li>
                  <SideNav />
                </li>
                <li className="mt-auto pb-2">
                  <SideNavLeftBottom />
                </li>
              </ul>
            </nav>
          </div>
        </div>

        <div className="lg:pl-72">
          <main className="py-10 ">
            <Outlet />
          </main>
        </div>
      </div>
    </>
  )
}
