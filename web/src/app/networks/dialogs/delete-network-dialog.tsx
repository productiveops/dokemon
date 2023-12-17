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
import useNetworks from "@/hooks/useNetworks"
import { toast } from "@/components/ui/use-toast"
import { INetwork } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import { useParams } from "react-router-dom"

export default function DeleteNetworkDialog({
  openState,
  setOpenState,
  network,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  network: INetwork
}) {
  const { nodeId } = useParams()
  const { mutateNetworks } = useNetworks(nodeId!)

  const [isSaving, setIsSaving] = useState(false)

  const handleDeleteNetwork = async () => {
    setIsSaving(true)

    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/networks/remove`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: network.id, force: true }),
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
      mutateNetworks()
      setTimeout(() => {
        setOpenState(false)
        toast({
          title: "Success!",
          description: "Network deleted.",
        })
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Network</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <p>{`Are you sure you want to delete network '${network.name}'?`}</p>
        </div>
        <DialogFooter>
          <fieldset disabled={isSaving} className="group">
            <Button
              variant={"destructive"}
              className={cn("relative w-24 group-disabled:pointer-events-none")}
              disabled={isSaving}
              onClick={() => handleDeleteNetwork()}
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
