import { TrashIcon } from "@heroicons/react/24/solid"
import Loading from "@/components/widgets/loading"
import {
  Breadcrumb,
  BreadcrumbCurrent,
  BreadcrumbLink,
  BreadcrumbSeparator,
} from "@/components/widgets/breadcrumb"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { INetwork } from "@/lib/api-models"
import { useState } from "react"
import useNetworks from "@/hooks/useNetworks"
import DeleteNetworkDialog from "./dialogs/delete-network-dialog"
import PruneNetworksDialog from "./dialogs/prune-networks-dialog"
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import { useParams } from "react-router-dom"
import useNodeHead from "@/hooks/useNodeHead"
import { CLASSES_TABLE_ACTION_ICON } from "@/lib/utils"

export default function Networks() {
  const { nodeId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { isLoading, networks } = useNetworks(nodeId!)

  const [deleteNetworkOpen, setDeleteNetworkOpen] = useState(false)
  const [network, setNetwork] = useState<INetwork | null>(null)

  if (isLoading) return <Loading />

  const handleDeleteNetwork = (network: INetwork) => {
    setNetwork({ ...network })
    setDeleteNetworkOpen(true)
  }

  return (
    <MainArea>
      {deleteNetworkOpen && (
        <DeleteNetworkDialog
          openState={deleteNetworkOpen}
          setOpenState={setDeleteNetworkOpen}
          network={network!}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeHead?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Networks</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <PruneNetworksDialog />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Id</TableHead>
              <TableHead scope="col">Name</TableHead>
              <TableHead scope="col">Driver</TableHead>
              <TableHead scope="col">Scope</TableHead>
              <TableHead scope="col">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!networks?.items && (
              <TableRow>
                <TableCell colSpan={5} className="text-center">
                  No data to display
                </TableCell>
              </TableRow>
            )}
            {networks?.items &&
              networks?.items.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.id.substring(0, 12)}</TableCell>
                  <TableCell>{item.name}</TableCell>
                  <TableCell>{item.driver}</TableCell>
                  <TableCell>{item.scope}</TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size={"sm"}
                      title="Delete"
                      onClick={() => handleDeleteNetwork(item)}
                    >
                      <TrashIcon className={CLASSES_TABLE_ACTION_ICON} />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </MainContent>
    </MainArea>
  )
}
