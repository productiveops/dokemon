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
import { useParams } from "react-router-dom"
import { useEffect, useState } from "react"
import { wsApiBaseUrl } from "@/lib/api-base-url"
import { AttachAddon } from "@xterm/addon-attach"
import { FitAddon } from "@xterm/addon-fit"
import { newTerminal, recreateTerminalElement } from "@/lib/utils"
import useNodeHead from "@/hooks/useNodeHead"
import useNodeComposeItem from "@/hooks/useNodeComposeItem"

export default function ComposeLogs() {
  const { nodeId, composeProjectId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { nodeComposeItem } = useNodeComposeItem(nodeId!, composeProjectId!)
  const [socket, setSocket] = useState<WebSocket>(null!)

  useEffect(() => {
    const terminal = newTerminal()

    if (socket) socket.close()
    const s = new WebSocket(
      `${wsApiBaseUrl()}/nodes/${nodeId}/compose/${composeProjectId}/logs`
    )
    setSocket(s)

    terminal.loadAddon(new AttachAddon(s))
    const fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)

    const terminalEl = recreateTerminalElement("terminalContainer", "terminal")
    terminal.open(terminalEl!)
    fitAddon.fit()
    addEventListener("resize", () => {
      fitAddon?.fit()
    })
  }, [composeProjectId])

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
          <BreadcrumbCurrent>Logs</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions></TopBarActions>
      </TopBar>
      <MainContent>
        <div id="terminalContainer">
          <div id="terminal"></div>
        </div>
      </MainContent>
    </MainArea>
  )
}
