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
import { useParams } from "react-router-dom"
import useNodeHead from "@/hooks/useNodeHead"
import TableButtonDelete from "@/components/widgets/table-button-delete"
import { TableNoData } from "@/components/widgets/table-no-data"
import DeleteDialog from "@/components/delete-dialog"
import { convertByteToMb, toastFailed, toastSuccess } from "@/lib/utils"
import apiBaseUrl from "@/lib/api-base-url"

export default function Volumes() {
  const { nodeId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { isLoading, volumes, mutateVolumes } = useVolumes(nodeId!)

  const [volume, setVolume] = useState<IVolume | null>(null)
  const [deleteVolumeOpenConfirmation, setDeleteVolumeOpenConfirmation] =
    useState(false)
  const [pruneVolumesOpenConfirmation, setPruneVolumesOpenConfirmation] =
    useState(false)
  const [deleteInProgress, setDeleteInProgress] = useState(false)
  const [pruneInProgress, setPruneInProgress] = useState(false)

  if (isLoading) return <Loading />

  const handleDeleteVolumeConfirmation = (volume: IVolume) => {
    setVolume({ ...volume })
    setDeleteVolumeOpenConfirmation(true)
  }

  const handleDelete = async () => {
    setDeleteInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/volumes/remove`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: volume?.name }),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setDeleteVolumeOpenConfirmation(false)
      toastFailed(r.errors?.body)
    } else {
      mutateVolumes()
      setTimeout(() => {
        setDeleteVolumeOpenConfirmation(false)
        toastSuccess("Volume deleted.")
      }, 500)
    }
    setDeleteInProgress(false)
  }

  const handlePrune = async () => {
    setPruneInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/volumes/prune`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ all: true }),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setPruneVolumesOpenConfirmation(false)
      toastFailed(r.errors?.body)
    } else {
      mutateVolumes()
      const r = await response.json()
      let description = "Nothing found to delete"
      if (r.volumesDeleted?.length > 0) {
        description = `${
          r.volumesDeleted.length
        } unused volumes deleted. Space reclaimed: ${convertByteToMb(
          r.spaceReclaimed
        )}`
      }
      setTimeout(async () => {
        setPruneVolumesOpenConfirmation(false)
        toastSuccess(description)
      }, 500)
    }
    setPruneInProgress(false)
  }

  return (
    <MainArea>
      {deleteVolumeOpenConfirmation && (
        <DeleteDialog
          openState={deleteVolumeOpenConfirmation}
          setOpenState={setDeleteVolumeOpenConfirmation}
          deleteCaption=""
          deleteHandler={handleDelete}
          isProcessing={deleteInProgress}
          title="Delete Volume"
          message={`Are you sure you want to delete volume '${volume?.name}?'`}
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
          <DeleteDialog
            widthClass="w-42"
            deleteCaption="Delete Unused (Prune All)"
            deleteHandler={handlePrune}
            isProcessing={pruneInProgress}
            title="Delete Unused"
            message={`Are you sure you want to delete all unused volumes?`}
          />
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
                        handleDeleteVolumeConfirmation(item)
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
