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
import { ICredentialHead } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useCredentials from "@/hooks/useCredentials"

export default function DeleteCredentialDialog({
  openState,
  setOpenState,
  credentialHead,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  credentialHead: ICredentialHead
}) {
  const { mutateCredentials } = useCredentials()
  const [isSaving, setIsSaving] = useState(false)

  const handleDelete = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/credentials/${credentialHead.id}`,
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
      mutateCredentials()
      setTimeout(() => {
        setOpenState(false)
        toastSuccess("Credential deleted.")
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Credential</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete credential '${credentialHead.name}'?`}</p>
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
