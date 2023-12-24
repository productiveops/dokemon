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
import { IImage } from "@/lib/api-models"
import { useState } from "react"
import useImages from "@/hooks/useImages"
import { convertByteToMb, toastFailed, toastSuccess } from "@/lib/utils"
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import { useParams } from "react-router-dom"
import useNodeHead from "@/hooks/useNodeHead"
import TableButtonDelete from "@/components/widgets/table-button-delete"
import { TableNoData } from "@/components/widgets/table-no-data"
import apiBaseUrl from "@/lib/api-base-url"
import DeleteDialog from "@/components/delete-dialog"

export default function Images() {
  const { nodeId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { isLoading, images, mutateImages } = useImages(nodeId!)
  const [image, setImage] = useState<IImage | null>(null)
  const [deleteImageConfirmationOpen, setDeleteImageConfirmationOpen] =
    useState(false)
  const [deleteInProgress, setDeleteInProgress] = useState(false)
  const [pruneInProgress, setPruneInProgress] = useState(false)

  if (isLoading) return <Loading />

  const handleDeleteImageConfirmation = (image: IImage) => {
    setImage({ ...image })
    setDeleteImageConfirmationOpen(true)
  }

  const handleDelete = async () => {
    setDeleteInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/images/remove`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: image?.id, force: true }),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setDeleteImageConfirmationOpen(false)
      toastFailed(r.errors?.body)
    } else {
      mutateImages()
      setTimeout(() => {
        setDeleteImageConfirmationOpen(false)
        toastSuccess("Image deleted.")
      }, 500)
    }
    setDeleteInProgress(false)
  }

  const handlePrune = async () => {
    setPruneInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/images/prune`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ all: true }),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      toastFailed(r.errors?.body)
    } else {
      mutateImages()
      const r = await response.json()
      let description = "Nothing found to delete"
      if (r.imagesDeleted?.length > 0) {
        description = `Unused images deleted. Space reclaimed: ${convertByteToMb(
          r.spaceReclaimed
        )}`
      }
      setTimeout(async () => {
        toastSuccess(description)
      }, 500)
    }
    setPruneInProgress(false)
  }

  return (
    <MainArea>
      {deleteImageConfirmationOpen && (
        <DeleteDialog
          openState={deleteImageConfirmationOpen}
          setOpenState={setDeleteImageConfirmationOpen}
          deleteCaption=""
          deleteHandler={handleDelete}
          isProcessing={deleteInProgress}
          title="Delete Image"
          message={`Are you sure you want to delete image '${image?.name}?'`}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeHead?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Images</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <DeleteDialog
            widthClass="w-42"
            deleteCaption="Delete Unused (Prune All)"
            deleteHandler={handlePrune}
            isProcessing={pruneInProgress}
            title="Delete Unused"
            message={`Are you sure you want to delete all unused images?`}
          />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Id</TableHead>
              <TableHead scope="col">Name</TableHead>
              <TableHead scope="col">Tag</TableHead>
              <TableHead scope="col">Size</TableHead>
              <TableHead scope="col">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {images?.items?.length === 0 && <TableNoData colSpan={5} />}
            {images?.items &&
              images?.items.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.id.substring(7, 19)}</TableCell>
                  <TableCell>{item.name}</TableCell>
                  <TableCell>
                    {item.tag}{" "}
                    {item.dangling ? (
                      <span className="text-xs text-red-400"> (Dangling)</span>
                    ) : (
                      ""
                    )}
                  </TableCell>
                  <TableCell>{convertByteToMb(item.size)}</TableCell>
                  <TableCell className="text-right">
                    <TableButtonDelete
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteImageConfirmation(item)
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
