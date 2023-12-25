import {
  Breadcrumb,
  BreadcrumbCurrent,
  BreadcrumbLink,
  BreadcrumbSeparator,
} from "@/components/widgets/breadcrumb"
import Loading from "@/components/widgets/loading"
import MainArea from "@/components/widgets/main-area"
import MainContent from "@/components/widgets/main-content"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import useContainers from "@/hooks/useContainers"
import useNodeHead from "@/hooks/useNodeHead"
import { wsApiBaseUrl } from "@/lib/api-base-url"
import { IContainer } from "@/lib/api-models"
import { newTerminal, recreateTerminalElement } from "@/lib/utils"
import { AttachAddon } from "@xterm/addon-attach"
import { FitAddon } from "@xterm/addon-fit"
import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"

export default function ContainerTerminal() {
  const { nodeId, containerId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { isLoading, containers } = useContainers(nodeId!)
  const [container, setContainer] = useState<IContainer>(null!)
  const [socket, setSocket] = useState<WebSocket>(null!)
  const [lastContainerId, setLastContainerId] = useState<string>(null!)

  useEffect(() => {
    const c = containers?.items.filter((x) => x.name === containerId)[0]!
    setContainer(c)
  }, [containers])

  useEffect(() => {
    // Currently we use useContainers() hook which refreshes on focus resulting
    // this effect to be triggered and the terminal is reloaded.
    // As a workaround we check below if ID is last opened container is same and this one
    // and do not recreate terminal if its the same. It is better implement
    // useContainer() to get a single container and with refresh on focus disabled
    if (container?.id === lastContainerId) return

    if (container?.state === "running") {
      setLastContainerId(container.id)
      const terminal = newTerminal(false)

      if (socket) socket.close()
      const s = new WebSocket(
        `${wsApiBaseUrl()}/nodes/${nodeId}/containers/${containerId}/terminal`
      )
      setSocket(s)

      terminal.loadAddon(new AttachAddon(s))
      const fitAddon = new FitAddon()
      terminal.loadAddon(fitAddon)

      const terminalEl = recreateTerminalElement(
        "terminalContainer",
        "terminal"
      )
      terminal.open(terminalEl!)
      fitAddon.fit()
      addEventListener("resize", () => {
        fitAddon?.fit()
      })
    }
  }, [container])

  if (isLoading) return <Loading />

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
            Terminal for <span className="font-semibold">{containerId}</span>
          </BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions></TopBarActions>
      </TopBar>
      <MainContent>
        {container?.state === "running" && (
          <div id="terminalContainer">
            <div id="terminal" className=""></div>
          </div>
        )}
        {!container ||
          (container?.state !== "running" && (
            <span>Container is not running</span>
          ))}
      </MainContent>
    </MainArea>
  )
}
