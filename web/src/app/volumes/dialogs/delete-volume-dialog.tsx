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
import useVolumes from "@/hooks/useVolumes"
import { toast } from "@/components/ui/use-toast"
import { IVolume } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import { useParams } from "react-router-dom"

export default function DeleteVolumeDialog({
  openState,
  setOpenState,
  volume,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  volume: IVolume
}) {
  const { nodeId } = useParams()
  const { mutateVolumes } = useVolumes(nodeId!)

  const [isSaving, setIsSaving] = useState(false)

  const handleDeleteVolume = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/volumes/remove`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: volume.name }),
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
      mutateVolumes()
      setTimeout(() => {
        setOpenState(false)
        toast({
          title: "Success!",
          description: "Volume deleted.",
        })
      }, 500)
    }
    setIsSaving(false)
  }

  function shortName(name: string) {
    if (name?.length > 20) {
      return name.substring(0, 30) + "..."
    }
    return name
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Volume</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete volume '${shortName(
            volume.name
          )}'?`}</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isSaving} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isSaving}
              onClick={() => handleDeleteVolume()}
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
