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
import useVariables from "@/hooks/useVariables"
import apiBaseUrl from "@/lib/api-base-url"
import { IVariableHead } from "@/lib/api-models"

export default function DeleteVariableDialog({
  openState,
  setOpenState,
  variableHead,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  variableHead: IVariableHead
}) {
  const [isSaving, setIsSaving] = useState(false)
  const { mutateVariables } = useVariables()

  const handleDeleteVariable = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/variables/${variableHead.id}`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      setOpenState(false)
    } else {
      mutateVariables()
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
          <p>{`Are you sure you want to delete variable '${variableHead.name}'?`}</p>
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
