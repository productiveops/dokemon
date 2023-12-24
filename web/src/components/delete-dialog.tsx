import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Icons } from "@/components/icons"
import { cn } from "@/lib/utils"
import { DialogTrigger } from "@radix-ui/react-dialog"

export default function DeleteDialog({
  deleteCaption,
  title,
  message,
  deleteHandler,
  isProcessing,
}: {
  deleteCaption: string
  title: string
  message: string
  deleteHandler: React.MouseEventHandler
  isProcessing: boolean
}) {
  const [open, setOpen] = useState(false)

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant={"destructive"} className="ml-auto w-24">
          {deleteCaption}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{message}</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isProcessing} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isProcessing}
              onClick={deleteHandler}
            >
              <Icons.spinner
                className={cn(
                  "absolute animate-spin text-slate-100 group-enabled:opacity-0"
                )}
              />
              <span className={cn("group-disabled:opacity-0")}>Delete</span>
            </Button>
            <Button
              variant={"secondary"}
              className="ml-2 w-24"
              onClick={() => setOpen(false)}
            >
              Cancel
            </Button>
          </fieldset>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
