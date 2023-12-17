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
import useImages from "@/hooks/useImages"
import { toast } from "@/components/ui/use-toast"
import { IImage } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import { useParams } from "react-router-dom"

export default function DeleteImageDialog({
  openState,
  setOpenState,
  image,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  image: IImage
}) {
  const { nodeId } = useParams()
  const { mutateImages } = useImages(nodeId!)
  const [isSaving, setIsSaving] = useState(false)

  const handleDeleteImage = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/images/remove`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: image.id, force: true }),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setOpenState(false)
      toast({
        variant: "destructive",
        title: "Failed",
        description: r.errors?.body,
      })
    } else {
      mutateImages()
      setTimeout(() => {
        setOpenState(false)
        toast({
          title: "Success!",
          description: "Image deleted.",
        })
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Image</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete image '${image.name}:${image.tag}'?`}</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isSaving} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isSaving}
              onClick={() => handleDeleteImage()}
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
