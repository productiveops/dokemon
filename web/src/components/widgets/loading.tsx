import { Icons } from "@/components/icons"
import { cn } from "@/lib/utils"

export default function Loading() {
  return (
    <div
      className="flex flex-col items-center justify-center"
      style={{ height: "300px" }}
    >
      <Icons.spinner className={cn("absolute animate-spin text-blue-700")} />
    </div>
  )
}
