import { useNavigate, useParams } from "react-router-dom"
import useNodeComposeItem from "@/hooks/useNodeComposeItem"
import ComposeDefinitionGitHub from "./compose-definition-github"
import ComposeDefinitionLocal from "./compose-definition-local"
import apiBaseUrl, { wsApiBaseUrl } from "@/lib/api-base-url"
import {
  newTerminal,
  recreateTerminalElement,
  toastFailed,
  toastSuccess,
} from "@/lib/utils"
import { FitAddon } from "@xterm/addon-fit"
import { AttachAddon } from "@xterm/addon-attach"
import { useState } from "react"

export default function ComposeDefinition() {
  const { nodeId } = useParams()
  const { composeProjectId } = useParams()
  const { nodeComposeItem, mutateNodeComposeItem } = useNodeComposeItem(
    nodeId!,
    composeProjectId!
  )
  const navigate = useNavigate()
  const [logsOpen, setLogsOpen] = useState(false)
  const [editing, setEditing] = useState(false)

  let terminal = newTerminal()
  let fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)

  function handleComposeAction(action: string) {
    setLogsOpen(true)

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
      mutateNodeComposeItem()
      setTimeout(() => {
        toastSuccess("Compose project deleted.")
        navigate(`/nodes/${nodeId}/compose`)
      }, 500)
    }
    setDeleteInProgress(false)
  }

  if (nodeComposeItem?.type === "github") {
    return (
      <ComposeDefinitionGitHub
        deleteHandler={handleDelete}
        deleteProcessingStatus={deleteInProgress}
        composeActionHandler={handleComposeAction}
        logsOpen={logsOpen}
        setLogsOpen={setLogsOpen}
        editing={editing}
        setEditing={setEditing}
      />
    )
  } else {
    return (
      <ComposeDefinitionLocal
        deleteHandler={handleDelete}
        deleteProcessingStatus={deleteInProgress}
        composeActionHandler={handleComposeAction}
        logsOpen={logsOpen}
        setLogsOpen={setLogsOpen}
        editing={editing}
        setEditing={setEditing}
      />
    )
  }
}
