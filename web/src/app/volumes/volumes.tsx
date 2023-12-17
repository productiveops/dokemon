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
import { IVolume } from "@/lib/api-models"
import { useState } from "react"
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import useVolumes from "@/hooks/useVolumes"
import DeleteVolumeDialog from "./dialogs/delete-volume-dialog"
import PruneVolumesDialog from "./dialogs/prune-volumes-dialog"
import { useParams } from "react-router-dom"
import useNodeHead from "@/hooks/useNodeHead"
import { CLASSES_TABLE_ACTION_ICON } from "@/lib/utils"

export default function Volumes() {
  const { nodeId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { isLoading, volumes } = useVolumes(nodeId!)

  const [deleteVolumeOpen, setDeleteVolumeOpen] = useState(false)
  const [volume, setVolume] = useState<IVolume | null>(null)

  if (isLoading) return <Loading />

  const handleDeleteVolume = (volume: IVolume) => {
    setVolume({ ...volume })
    setDeleteVolumeOpen(true)
  }

  return (
    <MainArea>
      {deleteVolumeOpen && (
        <DeleteVolumeDialog
          openState={deleteVolumeOpen}
          setOpenState={setDeleteVolumeOpen}
          volume={volume!}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeHead?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Volumes</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <PruneVolumesDialog />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Driver</TableHead>
              <TableHead scope="col">Name</TableHead>
              <TableHead scope="col">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!volumes?.items && (
              <TableRow>
                <TableCell colSpan={5} className="text-center">
                  No data to display
                </TableCell>
              </TableRow>
            )}
            {volumes?.items &&
              volumes?.items.map((item) => (
                <TableRow key={item.name}>
                  <TableCell>{item.driver}</TableCell>
                  <TableCell>{item.name}</TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size={"sm"}
                      title="Delete"
                      onClick={() => handleDeleteVolume(item)}
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
