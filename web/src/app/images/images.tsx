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
import { IImage } from "@/lib/api-models"
import { useState } from "react"
import useImages from "@/hooks/useImages"
import DeleteImageDialog from "./dialogs/delete-image-dialog"
import { CLASSES_TABLE_ACTION_ICON, convertByteToMb } from "@/lib/utils"
import PruneImagesDialog from "./dialogs/prune-images-dialog"
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import { useParams } from "react-router-dom"
import useNodeHead from "@/hooks/useNodeHead"

export default function Images() {
  const { nodeId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { isLoading, images } = useImages(nodeId!)
  const [deleteImageOpen, setDeleteImageOpen] = useState(false)
  const [image, setImage] = useState<IImage | null>(null)

  if (isLoading) return <Loading />

  const handleDeleteImage = (image: IImage) => {
    setImage({ ...image })
    setDeleteImageOpen(true)
  }

  return (
    <MainArea>
      {deleteImageOpen && (
        <DeleteImageDialog
          openState={deleteImageOpen}
          setOpenState={setDeleteImageOpen}
          image={image!}
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
          <PruneImagesDialog />
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
            {images?.items?.length === 0 && (
              <TableRow>
                <TableCell colSpan={5} className="text-center">
                  No data to display
                </TableCell>
              </TableRow>
            )}
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
                    <Button
                      variant="ghost"
                      size={"sm"}
                      title="Delete"
                      onClick={() => handleDeleteImage(item)}
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
