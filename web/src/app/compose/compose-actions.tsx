import {
  Breadcrumb,
  BreadcrumbCurrent,
  BreadcrumbLink,
  BreadcrumbSeparator,
} from "@/components/widgets/breadcrumb"
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import { useNavigate, useParams } from "react-router-dom"
import { useEffect, useState } from "react"
import apiBaseUrl, { wsApiBaseUrl } from "@/lib/api-base-url"
import { AttachAddon } from "@xterm/addon-attach"
import { FitAddon } from "@xterm/addon-fit"
import { Button } from "@/components/ui/button"
import {
  newTerminal,
  recreateTerminalElement,
  toastFailed,
  toastSuccess,
} from "@/lib/utils"
import useNodeHead from "@/hooks/useNodeHead"
import useNodeComposeItem from "@/hooks/useNodeComposeItem"
import useComposeLibraryItemList from "@/hooks/useComposeLibraryItemList"
import DeleteDialog from "@/components/delete-dialog"

export default function ComposeActions() {
  const { nodeId, composeProjectId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { nodeComposeItem, mutateNodeComposeItem } = useNodeComposeItem(
    nodeId!,
    composeProjectId!
  )
  const { mutateComposeLibraryItemList } = useComposeLibraryItemList()
  const navigate = useNavigate()

  let terminal = newTerminal()
  let fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)

  useEffect(() => {
    const el = document.getElementById("terminal")
    if (el) {
      terminal.open(el)
      fitAddon.fit()
      addEventListener("resize", () => {
        fitAddon?.fit()
      })
    }
  }, [])

  function handleAction(action: string) {
    const url = `${wsApiBaseUrl()}/nodes/${nodeId}/compose/${composeProjectId}/${action}`
    const socket = new WebSocket(url)

    socket.onclose = function () {
      mutateNodeComposeItem()
    }

    terminal = newTerminal()
    fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)
    terminal.loadAddon(new AttachAddon(socket))

    const terminalEl = recreateTerminalElement("terminalContainer", "terminal")
    terminal.open(terminalEl!)
    fitAddon.fit()
    addEventListener("resize", () => {
      fitAddon?.fit()
    })
  }

  const [deleteInProgress, setDeleteInProgress] = useState(false)
  const handleDelete = async () => {
    setDeleteInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/compose/${composeProjectId}`,
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
        navigate(`/nodes/${nodeId}/compose`)
      }, 500)
    }
    setDeleteInProgress(false)
  }

  return (
    <MainArea>
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeHead?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbLink to={`/nodes/${nodeId}/compose`}>
            Compose
          </BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>
            {nodeComposeItem?.projectName}{" "}
            {nodeComposeItem?.status?.startsWith("running")
              ? "(Running)"
              : "(Not running)"}
          </BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Actions</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions></TopBarActions>
      </TopBar>
      <MainContent>
        <div className="mb-4 flex gap-2">
          <Button
            className="w-24"
            variant={"default"}
            onClick={() => handleAction("pull")}
          >
            Pull
          </Button>
          <Button
            className="w-24"
            variant={"default"}
            onClick={() => handleAction("up")}
          >
            Up
          </Button>
          <Button
            className="w-24"
            variant={"destructive"}
            onClick={() => handleAction("down")}
          >
            Down
          </Button>
          <DeleteDialog
            deleteCaption="Delete"
            title="Delete Compose Project"
            message={`Are you sure you want to delete project '${nodeComposeItem?.projectName}'?`}
            deleteHandler={handleDelete}
            isProcessing={deleteInProgress}
          />
        </div>
        <div id="terminalContainer">
          <h2 className="mb-2 font-bold">Action Logs</h2>
          <div id="terminal"></div>
        </div>
      </MainContent>
    </MainArea>
  )
}
