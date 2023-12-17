import { Dispatch, SetStateAction, useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { VERSION } from "@/lib/version"
import useSetting from "@/hooks/useSetting"

export default function RegisterNodeDialog({
  open,
  setOpen,
  token,
  updateAgent,
}: {
  open: boolean
  setOpen: Dispatch<SetStateAction<boolean>>
  token: string
  updateAgent: boolean
}) {
  const { setting } = useSetting("SERVER_URL")
  const [copyCaption, setCopyCaption] = useState("Copy")

  const handleCopy = () => {
    let value = (document.getElementById("command") as HTMLInputElement).value
    value = value.replace("{HIDDEN}", token)
    navigator.clipboard.writeText(value)

    setCopyCaption("Copied!")
    setTimeout(() => {
      setCopyCaption("Copy")
    }, 1000)
  }

  const script = () => {
    const c = `docker rm -f dokemon-agent > /dev/null 2>&1
docker run \\
    -e SERVER_URL=${setting?.value} \\
    -e TOKEN={HIDDEN} \\
    -v /var/run/docker.sock:/var/run/docker.sock \\
    --name dokemon-agent --restart unless-stopped \\
    -d productiveops/dokemon-agent:${VERSION}
  `
    if (navigator?.clipboard) return c
    return c.replace("{HIDDEN}", token)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {updateAgent ? "Update Agent" : "Register Node"}
          </DialogTitle>
        </DialogHeader>
        <p className="-mb-6">Run below script on the target server.</p>
        <div className="grid gap-4 py-4 group-disabled:opacity-50">
          <textarea
            id="command"
            className="border-2 p-1 text-sm dark:bg-slate-900"
            readOnly
            value={script()}
            cols={50}
            rows={6}
          ></textarea>
          {navigator?.clipboard && (
            <p className="-mb-4 text-xs">
              <strong>Note:</strong> Use the <em>Copy</em> button below to copy
              the script. The button adds the value for the TOKEN variable in
              the copied script. Directly copying the script text will not have
              the correct value for the TOKEN variable and will not work.
            </p>
          )}
        </div>
        <DialogFooter>
          <fieldset className="group">
            {navigator?.clipboard && (
              <Button
                variant={"default"}
                className="ml-2 w-24"
                onClick={() => handleCopy()}
              >
                {copyCaption}
              </Button>
            )}
            <Button
              variant={"secondary"}
              className="ml-2 w-24"
              onClick={() => setOpen(false)}
            >
              Close
            </Button>
          </fieldset>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
