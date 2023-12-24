import { Dispatch, SetStateAction, useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Icons } from "@/components/icons"
import { cn, toastFailed, toastSuccess } from "@/lib/utils"
import { INodeHead } from "@/lib/api-models"
import useNodes from "@/hooks/useNodes"
import { apiNodesDelete } from "@/lib/api"

export default function DeleteNodeDialog({
  openState,
  setOpenState,
  nodeHead,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  nodeHead: INodeHead
}) {
  const { mutateNodes } = useNodes()
  const [isSaving, setIsSaving] = useState(false)

  const handleDelete = async () => {
    setIsSaving(true)
    const response = await apiNodesDelete(nodeHead.id)
    if (!response.ok) {
      const r = await response.json()
      setOpenState(false)
      toastFailed(r.errors?.body)
    } else {
      mutateNodes()
      setTimeout(() => {
        setOpenState(false)
        toastSuccess("Node deleted.")
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Node</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete node '${nodeHead.name}'?`}</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isSaving} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isSaving}
              onClick={() => handleDelete()}
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
              onClick={() => setOpenState(false)}
            >
              Cancel
            </Button>
          </fieldset>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
