import {
  Breadcrumb,
  BreadcrumbCurrent,
  BreadcrumbLink,
  BreadcrumbSeparator,
} from "@/components/widgets/breadcrumb"
import MainArea from "@/components/widgets/main-area"
import MainContent from "@/components/widgets/main-content"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import { FitAddon } from "@xterm/addon-fit"
import { useParams } from "react-router-dom"
import "/node_modules/xterm/css/xterm.css"
import { useEffect, useState } from "react"
import { wsApiBaseUrl } from "@/lib/api-base-url"
import {
  download,
  downloadTerminalTextAsFile,
  newTerminal,
  recreateTerminalElement,
} from "@/lib/utils"
import { AttachAddon } from "@xterm/addon-attach"
import useNodeHead from "@/hooks/useNodeHead"
import { Button } from "@/components/ui/button"
import { Terminal } from "@xterm/xterm"

export default function ContainerLogs() {
  const { nodeId, containerId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const [socket, setSocket] = useState<WebSocket>(null!)
  const [terminal, setTerminal] = useState<Terminal>(null!)

  useEffect(() => {
    const t = newTerminal()
    setTerminal(t)

    if (socket) socket.close()
    const s = new WebSocket(
      `${wsApiBaseUrl()}/nodes/${nodeId}/containers/${containerId}/logs`
    )
    setSocket(s)

    t.loadAddon(new AttachAddon(s))
    const fitAddon = new FitAddon()
    t.loadAddon(fitAddon)

    const terminalEl = recreateTerminalElement("terminalContainer", "terminal")
    t.open(terminalEl!)
    fitAddon.fit()
    addEventListener("resize", () => {
      fitAddon?.fit()
    })
  }, [containerId])

  const handleDownload = () => {
    downloadTerminalTextAsFile(terminal, `logs_${containerId}.txt`)
  }

  return (
    <MainArea>
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeHead?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbLink to={`/nodes/${nodeId}/containers`}>
            Containers
          </BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>
            Logs for <span className="font-semibold">{containerId}</span>
          </BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <Button onClick={handleDownload}>Download to File</Button>
        </TopBarActions>
      </TopBar>
      <MainContent>
        <div id="terminalContainer">
          <div id="terminal"></div>
        </div>
      </MainContent>
    </MainArea>
  )
}
