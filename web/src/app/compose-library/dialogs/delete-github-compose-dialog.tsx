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
import { cn, toastFailed, toastSuccess } from "@/lib/utils"
import apiBaseUrl from "@/lib/api-base-url"
import { useNavigate, useParams } from "react-router-dom"
import { DialogTrigger } from "@radix-ui/react-dialog"
import useComposeLibraryItemList from "@/hooks/useComposeLibraryItemList"
import useGitHubComposeLibraryItem from "@/hooks/useGitHubComposeLibraryItem"

export default function DeleteGitHubComposeDialog() {
  const { composeProjectId } = useParams()
  const { gitHubComposeLibraryItem } = useGitHubComposeLibraryItem(
    composeProjectId!
  )
  const navigate = useNavigate()

  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const { mutateComposeLibraryItemList } = useComposeLibraryItemList()

  const handleDelete = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/composelibrary/github/${composeProjectId}`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      const r = await response.json()
      toastFailed(r.errors?.body)
    } else {
      mutateComposeLibraryItemList()
      setTimeout(() => {
        toastSuccess("Compose project deleted.")
        navigate("/composelibrary")
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant={"destructive"} className="ml-auto w-24">
          Delete
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Compose Project</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete project '${gitHubComposeLibraryItem?.projectName}'?`}</p>
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
