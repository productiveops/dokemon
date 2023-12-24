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
import { IEnvironmentHead } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useEnvironments from "@/hooks/useEnvironments"

export default function DeleteEnvironmentDialog({
  openState,
  setOpenState,
  environmentHead,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  environmentHead: IEnvironmentHead
}) {
  const { mutateEnvironments } = useEnvironments()
  const [isSaving, setIsSaving] = useState(false)

  const handleDelete = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/environments/${environmentHead.id}`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setOpenState(false)
      toastFailed(r.errors?.body)
    } else {
      mutateEnvironments()
      setTimeout(() => {
        setOpenState(false)
        toastSuccess("Environment deleted.")
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Environment</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete environment '${environmentHead.name}'?`}</p>
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
