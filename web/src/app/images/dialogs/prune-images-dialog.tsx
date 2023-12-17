import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Icons } from "@/components/icons"
import { cn, convertByteToMb } from "@/lib/utils"
import useImages from "@/hooks/useImages"
import { toast } from "@/components/ui/use-toast"
import apiBaseUrl from "@/lib/api-base-url"
import { useParams } from "react-router-dom"

export default function PruneImagesDialog() {
  const { nodeId } = useParams()
  const { mutateImages } = useImages(nodeId!)

  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)

  const handlePruneImages = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/images/prune`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ all: true }),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setOpen(false)
      toast({
        variant: "destructive",
        title: "Failed",
        description: r.errors?.body,
      })
    } else {
      mutateImages()
      const r = await response.json()
      let description = "Nothing found to delete"
      if (r.imagesDeleted?.length > 0) {
        description = `Unused images deleted. Space reclaimed: ${convertByteToMb(
          r.spaceReclaimed
        )}`
      }
      setTimeout(async () => {
        setOpen(false)
        toast({
          title: "Success!",
          description: description,
        })
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant={"default"}>Delete Unused (Prune All)</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Unused</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>Are you sure you want to delete all unused images?</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isSaving} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isSaving}
              onClick={() => handlePruneImages()}
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
