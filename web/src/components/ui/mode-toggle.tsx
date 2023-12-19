import { Moon, Sun } from "lucide-react"

import { Button } from "@/components/ui/button"
import { useTheme } from "@/components/ui/theme-provider"

export function ModeToggle() {
  const { setTheme, theme } = useTheme()

  return (
    <Button
      variant={"outline"}
      size="icon"
      className="flex items-center justify-center border-0 bg-gray-900 text-slate-50 hover:bg-gray-800 hover:text-slate-50"
      title="Toggle Theme"
      onClick={() => {
        theme === "dark" ? setTheme("light") : setTheme("dark")
      }}
    >
      {theme === "light" && (
        <Sun className="h-[1.2rem] w-[1.2rem] transition-all" />
      )}
      {theme === "dark" && (
        <Moon className="h-[1.2rem] w-[1.2rem] transition-all" />
      )}
      <span className="sr-only">Toggle theme</span>
    </Button>
  )
}
