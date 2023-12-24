import Loading from "@/components/widgets/loading"
import { Breadcrumb, BreadcrumbCurrent } from "@/components/widgets/breadcrumb"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import TableButtonDelete from "@/components/widgets/table-button-delete"
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import AddNodeDialog from "./dialogs/add-node-dialog"
import { Button } from "@/components/ui/button"
import RegisterNodeDialog from "./dialogs/register-node-dialog"
import DeleteNodeDialog from "./dialogs/delete-node-dialog"
import EditServerUrlDialog from "./dialogs/edit-serverurl-dialog"
import { INodeHead } from "@/lib/api-models"
import { apiNodesGenerateToken } from "@/lib/api"
import {
  CLASSES_CLICKABLE_TABLE_ROW,
  cn,
  toastSomethingWentWrong,
} from "@/lib/utils"
import { useNavigate } from "react-router-dom"
import useNodes from "@/hooks/useNodes"
import { useState } from "react"
import useSetting from "@/hooks/useSetting"

export default function Nodes() {
  const navigate = useNavigate()
  const { isLoading, nodes } = useNodes()
  const { setting } = useSetting("SERVER_URL")
  const [token, setToken] = useState("")
  const [updateAgent, setUpdateAgent] = useState(false)
  const [registerNodeOpen, setRegisterNodeOpen] = useState(false)
  const [deleteNodeOpen, setDeleteNodeOpen] = useState(false)
  const [nodeHead, setNodeHead] = useState<INodeHead | null>(null)

  if (isLoading) return <Loading />

  const handleDeleteNode = (nodeHead: INodeHead) => {
    setNodeHead({ ...nodeHead })
    setDeleteNodeOpen(true)
  }

  const handleRegister = async (nodeId: number, update: boolean) => {
    const response = await apiNodesGenerateToken(nodeId)

    if (!response.ok) {
      toastSomethingWentWrong(
        "There was a problem when generating the registration token. Try again!"
      )
    } else {
      const data: { token: string } = await response.json()
      setToken(data.token)
      setUpdateAgent(update)
      setRegisterNodeOpen(true)
    }
  }

  return (
    <MainArea>
      {deleteNodeOpen && (
        <DeleteNodeDialog
          openState={deleteNodeOpen}
          setOpenState={setDeleteNodeOpen}
          nodeHead={nodeHead!}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbCurrent>Nodes</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <RegisterNodeDialog
            open={registerNodeOpen}
            setOpen={setRegisterNodeOpen}
            token={token}
            updateAgent={updateAgent}
          />
          <AddNodeDialog disabled={!setting?.value} />
          <EditServerUrlDialog />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">
                <span className="ml-3">Name</span>
              </TableHead>
              <TableHead scope="col">Environment</TableHead>
              <TableHead scope="col">Agent Version</TableHead>
              <TableHead scope="col">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {nodes?.totalRows === 0 && (
              <TableRow>
                <TableCell colSpan={4} className="text-center">
                  No data to display
                </TableCell>
              </TableRow>
            )}
            {nodes?.items &&
              nodes?.items.map((item) => (
                <TableRow
                  key={item.name}
                  className={CLASSES_CLICKABLE_TABLE_ROW}
                  onClick={() => {
                    if (item.registered) {
                      navigate(`/nodes/${item.id}/compose`)
                    }
                  }}
                >
                  <TableCell>
                    <NodeStatusIcon nodeHead={item} />
                    {item.name}
                  </TableCell>
                  <TableCell>
                    {item.environment ? item.environment : "-"}
                  </TableCell>
                  <TableCell>{getAgentVersion(item)}</TableCell>
                  <TableCell className="text-right">
                    {!item.registered && (
                      <Button
                        className="mr-4"
                        size={"sm"}
                        onClick={(e) => {
                          e.stopPropagation()
                          handleRegister(item.id, false)
                        }}
                      >
                        Register
                      </Button>
                    )}
                    {!isDokemonNode(item) && (
                      <TableButtonDelete
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDeleteNode(item)
                        }}
                      />
                    )}
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </MainContent>
    </MainArea>
  )
}

function isDokemonNode(nodeHead: INodeHead) {
  return nodeHead.id === 1
}

function NodeStatusIcon({ nodeHead }: { nodeHead: INodeHead }) {
  const statusClassName = nodeHead.online ? "text-green-600" : "text-slate-300"
  const title = nodeHead.online ? "Online" : "Offline"

  return (
    <span className={cn("-ml-2 mr-3 text-lg", statusClassName)} title={title}>
      ●
    </span>
  )
}

function getAgentVersion(nodeHead: INodeHead) {
  if (isDokemonNode(nodeHead)) return "N/A (Dokemon Server)"
  if (nodeHead.agentVersion) return nodeHead.agentVersion
  return "-"
}
