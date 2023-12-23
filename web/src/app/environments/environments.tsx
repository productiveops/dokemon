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
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { TrashIcon } from "@heroicons/react/24/solid"
import { IEnvironmentHead } from "@/lib/api-models"
import DeleteEnvironmentDialog from "./dialogs/delete-environment-dialog"
import EditEnvironmentDialog from "./dialogs/edit-environment-dialog"
import {
  CLASSES_CLICKABLE_TABLE_ROW,
  CLASSES_TABLE_ACTION_ICON,
} from "@/lib/utils"

export default function Environments() {
  const { isLoading, environments } = useEnvironments()
  const [editEnvironmentOpen, setEditEnvironmentOpen] = useState(false)
  const [deleteEnvironmentOpen, setDeleteEnvironmentOpen] = useState(false)
  const [environmentHead, setEnvironmentHead] =
    useState<IEnvironmentHead | null>(null)

  if (isLoading) return <Loading />

  const handleEditEnvironment = (environmentHead: IEnvironmentHead) => {
    setEnvironmentHead({ ...environmentHead })
    setEditEnvironmentOpen(true)
  }

  const handleDeleteEnvironment = (environmentHead: IEnvironmentHead) => {
    setEnvironmentHead({ ...environmentHead })
    setDeleteEnvironmentOpen(true)
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
      {deleteEnvironmentOpen && (
        <DeleteEnvironmentDialog
          openState={deleteEnvironmentOpen}
          setOpenState={setDeleteEnvironmentOpen}
          environmentHead={environmentHead!}
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
            {environments?.totalRows === 0 && (
              <TableRow>
                <TableCell colSpan={2} className="text-center">
                  No data to display
                </TableCell>
              </TableRow>
            )}
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
                    <Button
                      variant="ghost"
                      size={"sm"}
                      title="Delete"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteEnvironment(item)
                      }}
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
