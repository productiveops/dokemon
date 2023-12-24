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
import useContainers from "@/hooks/useContainers"
import { IContainer } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import { useParams } from "react-router-dom"

export default function DeleteContainerDialog({
  openState,
  setOpenState,
  container,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  container: IContainer
}) {
  const { nodeId } = useParams()
  const [isSaving, setIsSaving] = useState(false)
  const { mutateContainers } = useContainers(nodeId!)

  const handleDeleteContainer = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/containers/remove`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: container.id, force: true }),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setOpenState(false)
      toastFailed(r.errors?.body)
    } else {
      mutateContainers()
      setTimeout(() => {
        setOpenState(false)
        toastSuccess("Container deleted.")
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Container</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete container '${container.name}'?`}</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isSaving} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isSaving}
              onClick={() => handleDeleteContainer()}
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
