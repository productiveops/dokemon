import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { DialogTrigger } from "@radix-ui/react-dialog"
import SpinnerIcon from "./widgets/spinner-icon"

export default function DeleteDialog({
  openState,
  setOpenState,
  deleteCaption,
  title,
  message,
  deleteHandler,
  isProcessing,
  widthClass,
  buttonVisible,
}: {
  openState?: boolean
  setOpenState?: React.Dispatch<React.SetStateAction<boolean>>
  deleteCaption: string
  title: string
  message: string
  deleteHandler: React.MouseEventHandler
  isProcessing: boolean
  widthClass?: string
  buttonVisible?: boolean
}) {
  const [open, setOpen] = useState(false)

  return (
    <Dialog
      open={openState !== undefined ? openState : open}
      onOpenChange={setOpenState !== undefined ? setOpenState : setOpen}
    >
      {openState === undefined && (
        <DialogTrigger asChild>
          <Button
            variant={"destructive"}
            className={cn(
              "ml-1",
              widthClass ? widthClass : "w-24",
              buttonVisible === undefined || buttonVisible ? "" : "hidden"
            )}
          >
            {deleteCaption}
          </Button>
        </DialogTrigger>
      )}
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
              onClick={async (e) => {
                await deleteHandler(e)
                setOpenState !== undefined
                  ? setOpenState(false)
                  : setOpen(false)
              }}
            >
              <SpinnerIcon />
              <span className={cn("group-disabled:opacity-0")}>Delete</span>
            </Button>
            <Button
              variant={"secondary"}
              className="ml-2 w-24"
              onClick={() =>
                setOpenState !== undefined
                  ? setOpenState(false)
                  : setOpen(false)
              }
            >
              Cancel
            </Button>
          </fieldset>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
