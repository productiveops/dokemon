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
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import { useNavigate } from "react-router-dom"
import useNodes from "@/hooks/useNodes"
import AddNodeDialog from "./dialogs/add-node-dialog"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import apiBaseUrl from "@/lib/api-base-url"
import { toast } from "@/components/ui/use-toast"
import RegisterNodeDialog from "./dialogs/register-node-dialog"
import { TrashIcon } from "@heroicons/react/24/solid"
import { INodeHead } from "@/lib/api-models"
import DeleteNodeDialog from "./dialogs/delete-node-dialog"
import EditServerUrlDialog from "./dialogs/edit-serverurl-dialog"
import useSetting from "@/hooks/useSetting"
import { VERSION } from "@/lib/version"
import { CLASSES_TABLE_ACTION_ICON } from "@/lib/utils"

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
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/generatetoken`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      toast({
        variant: "destructive",
        title: "Something went wrong.",
        description:
          "There was a problem when generating the registration token. Try again!",
      })
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
              <TableHead scope="col">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!nodes?.items && (
              <TableRow>
                <TableCell colSpan={5} className="text-center">
                  No data to display
                </TableCell>
              </TableRow>
            )}
            {nodes?.items &&
              nodes?.items.map((item) => (
                <TableRow
                  key={item.name}
                  className="cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800"
                  onClick={() => {
                    if (item.registered) {
                      navigate(`/nodes/${item.id}/compose`)
                    }
                  }}
                >
                  <TableCell>
                    {item.online ? (
                      <span
                        className="-ml-2 mr-3 text-lg text-green-600"
                        title="Online"
                      >
                        ●
                      </span>
                    ) : (
                      <span
                        className="-ml-2 mr-3 text-lg text-slate-300"
                        title="Offline"
                      >
                        ●
                      </span>
                    )}
                    {item.name}
                    {item.id === 1 ? (
                      <span title="Dokemon Server"> *</span>
                    ) : (
                      ""
                    )}
                  </TableCell>
                  <TableCell>
                    {item.environment ? item.environment : "-"}
                  </TableCell>
                  <TableCell className="text-right">
                    {item.id !== 1 &&
                      item.registered &&
                      item.agentVersion !== VERSION && (
                        <Button
                          className="mr-4"
                          size={"sm"}
                          onClick={(e) => {
                            e.stopPropagation()
                            handleRegister(item.id, true)
                          }}
                        >
                          Update Agent
                        </Button>
                      )}
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
                    {item.id !== 1 && (
                      <Button
                        variant="ghost"
                        size={"sm"}
                        title="Delete"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDeleteNode(item)
                        }}
                      >
                        <TrashIcon className={CLASSES_TABLE_ACTION_ICON} />
                      </Button>
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
