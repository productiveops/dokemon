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
import TableButtonDelete from "@/components/widgets/table-button-delete"
import { TableNoData } from "@/components/widgets/table-no-data"

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
            {volumes?.items?.length === 0 && <TableNoData colSpan={3} />}
            {volumes?.items &&
              volumes?.items.map((item) => (
                <TableRow key={item.name}>
                  <TableCell>{item.driver}</TableCell>
                  <TableCell>{item.name}</TableCell>
                  <TableCell className="text-right">
                    <TableButtonDelete
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteVolume(item)
                      }}
                    />
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </MainContent>
    </MainArea>
  )
}
