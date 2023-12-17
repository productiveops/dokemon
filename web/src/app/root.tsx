import { Outlet, useNavigate } from "react-router-dom"
import SideNav from "../components/side-nav"
import axios from "axios"
import apiBaseUrl from "@/lib/api-base-url"
import { toast } from "@/components/ui/use-toast"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { LifebuoyIcon, UserCircleIcon } from "@heroicons/react/24/outline"
import { cn } from "@/lib/utils"
import { VERSION } from "@/lib/version"
import { ModeToggle } from "@/components/ui/mode-toggle"

export default function Root() {
  const navigate = useNavigate()

  async function handleSignout() {
    try {
      await axios(`${apiBaseUrl()}/users/logout`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      })

      navigate("/login")
    } catch (e) {
      if (axios.isAxiosError(e)) {
        toast({
          variant: "destructive",
          title: "Failed",
          description: e.response?.data.errors?.body,
        })
      }
    }
  }

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
                  <ul role="list" className="-mx-2 space-y-1">
                    <li>
                      <a
                        href="https://discord.gg/Nfevu4gJVG"
                        target="_blank"
                        className={cn(
                          "group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-gray-400 hover:bg-gray-800 hover:text-white"
                        )}
                      >
                        <LifebuoyIcon
                          className="h-6 w-6 shrink-0"
                          aria-hidden="true"
                        />
                        Support
                      </a>
                    </li>
                    <li>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <a
                            href="#"
                            className="group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-gray-400 hover:text-white focus:outline-none"
                          >
                            <UserCircleIcon
                              className="h-6 w-6 shrink-0"
                              aria-hidden="true"
                            />
                            {localStorage.getItem("userName")}
                          </a>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-56">
                          <DropdownMenuItem
                            onClick={() => {
                              navigate("/changepassword")
                            }}
                          >
                            Change Password
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={handleSignout}>
                            Sign Out
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </li>
                  </ul>
                </li>
              </ul>
            </nav>
          </div>
        </div>

        <div className="lg:pl-72">
          <main className="py-10 ">
            <div className="px-4 sm:px-6 lg:px-8">
              <Outlet />
            </div>
          </main>
        </div>
      </div>
    </>
  )
}
