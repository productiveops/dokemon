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
import useEnvironments from "@/hooks/useEnvironments"
import AddEnvironmentDialog from "./dialogs/add-environment-dialog"
import { useState } from "react"
import { IEnvironmentHead } from "@/lib/api-models"
import EditEnvironmentDialog from "./dialogs/edit-environment-dialog"
import {
  CLASSES_CLICKABLE_TABLE_ROW,
  toastFailed,
  toastSuccess,
} from "@/lib/utils"
import TableButtonDelete from "@/components/widgets/table-button-delete"
import { TableNoData } from "@/components/widgets/table-no-data"
import apiBaseUrl from "@/lib/api-base-url"
import DeleteDialog from "@/components/delete-dialog"

export default function Environments() {
  const { isLoading, environments, mutateEnvironments } = useEnvironments()
  const [environmentHead, setEnvironmentHead] =
    useState<IEnvironmentHead | null>(null)
  const [editEnvironmentOpen, setEditEnvironmentOpen] = useState(false)
  const [
    deleteEnvironmentConfirmationOpen,
    setDeleteEnvironmentConfirmationOpen,
  ] = useState(false)
  const [deleteInProgress, setDeleteInProgress] = useState(false)

  if (isLoading) return <Loading />

  const handleEditEnvironment = (environmentHead: IEnvironmentHead) => {
    setEnvironmentHead({ ...environmentHead })
    setEditEnvironmentOpen(true)
  }

  const handleDeleteEnvironmentConfirmation = (
    environmentHead: IEnvironmentHead
  ) => {
    setEnvironmentHead({ ...environmentHead })
    setDeleteEnvironmentConfirmationOpen(true)
  }

  const handleDelete = async () => {
    setDeleteInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/environments/${environmentHead?.id}`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setDeleteEnvironmentConfirmationOpen(false)
      toastFailed(r.errors?.body)
    } else {
      mutateEnvironments()
      setTimeout(() => {
        setDeleteEnvironmentConfirmationOpen(false)
        toastSuccess("Environment deleted.")
      }, 500)
    }
    setDeleteInProgress(false)
  }

  return (
    <MainArea>
      {editEnvironmentOpen && (
        <EditEnvironmentDialog
          openState={editEnvironmentOpen}
          setOpenState={setEditEnvironmentOpen}
          environmentHead={environmentHead!}
        />
      )}
      {deleteEnvironmentConfirmationOpen && (
        <DeleteDialog
          openState={deleteEnvironmentConfirmationOpen}
          setOpenState={setDeleteEnvironmentConfirmationOpen}
          deleteCaption=""
          deleteHandler={handleDelete}
          isProcessing={deleteInProgress}
          title="Delete Environment"
          message={`Are you sure you want to delete environment '${environmentHead?.name}?'`}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbCurrent>Environments</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <AddEnvironmentDialog />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Name</TableHead>
              <TableHead scope="col">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {environments?.totalRows === 0 && <TableNoData colSpan={2} />}
            {environments?.items &&
              environments?.items.map((item) => (
                <TableRow
                  key={item.name}
                  className={CLASSES_CLICKABLE_TABLE_ROW}
                  onClick={() => {
                    handleEditEnvironment(item)
                  }}
                >
                  <TableCell>{item.name}</TableCell>
                  <TableCell className="text-right">
                    <TableButtonDelete
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteEnvironmentConfirmation(item)
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
