import { cn } from "@/lib/utils"
import { Icons } from "../icons"

export default function SpinnerIcon() {
  return (
    <Icons.spinner
      className={cn(
        "absolute animate-spin text-slate-100 group-enabled:opacity-0"
      )}
    />
  )
}
