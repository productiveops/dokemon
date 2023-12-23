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
import { cn } from "@/lib/utils"
import apiBaseUrl from "@/lib/api-base-url"
import { INodeComposeVariable } from "@/lib/api-models"
import useNodeComposeVariables from "@/hooks/useNodeComposeVariables"

export default function DeleteComposeVariableDialog({
  openState,
  setOpenState,
  variable,
  nodeId,
  nodeComposeProjectId,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  variable: INodeComposeVariable
  nodeId: string
  nodeComposeProjectId: string
}) {
  const [isSaving, setIsSaving] = useState(false)
  const { mutateNodeComposeVariables } = useNodeComposeVariables(
    nodeId,
    nodeComposeProjectId
  )

  const handleDeleteVariable = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/compose/${nodeComposeProjectId}/variables/${
        variable.id
      }`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      setOpenState(false)
    } else {
      mutateNodeComposeVariables()
      setTimeout(() => {
        setOpenState(false)
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Variable</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete variable '${variable.name}'?`}</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isSaving} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isSaving}
              onClick={() => handleDeleteVariable()}
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
