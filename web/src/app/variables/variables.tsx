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
import useVariables from "@/hooks/useVariables"
import AddVariableDialog from "./dialogs/add-variable-dialog"
import { useState } from "react"
import { IVariableHead } from "@/lib/api-models"
import DeleteVariableDialog from "./dialogs/delete-variable-dialog"
import useEnvironmentsMap from "@/hooks/useEnvironmentsMap"
import { Checkbox } from "@/components/ui/checkbox"
import EditVariableDialog from "./dialogs/edit-variable-dialog"
import EditVariableValueDialog from "./dialogs/edit-variable-value-dialog"
import TableButtonDelete from "@/components/widgets/table-button-delete"
import TableButtonEdit from "@/components/widgets/table-button-edit"
import { TableNoData } from "@/components/widgets/table-no-data"

export default function Variables() {
  const { isLoading: mapIsLoading, environmentsMap } = useEnvironmentsMap()
  const { isLoading, variables } = useVariables()
  const [editVariableOpen, setEditVariableOpen] = useState(false)
  const [editVariableValueOpen, setEditVariableValueOpen] = useState(false)
  const [editVariableValueEnvironmentId, setEditVariableValueEnvironmentId] =
    useState<string | null>(null)
  const [deleteVariableOpen, setDeleteVariableOpen] = useState(false)
  const [variableHead, setVariableHead] = useState<IVariableHead | null>(null)

  if (mapIsLoading || isLoading) return <Loading />

  const handleEditVariable = (variableHead: IVariableHead) => {
    setVariableHead({ ...variableHead })
    setEditVariableOpen(true)
  }

  const handleEditVariableValue = (
    variableHead: IVariableHead,
    envId: string
  ) => {
    setVariableHead({ ...variableHead })
    setEditVariableValueOpen(true)
    setEditVariableValueEnvironmentId(envId)
  }

  const handleDeleteVariable = (variableHead: IVariableHead) => {
    setVariableHead({ ...variableHead })
    setDeleteVariableOpen(true)
  }

  return (
    <MainArea>
      {editVariableOpen && (
        <EditVariableDialog
          openState={editVariableOpen}
          setOpenState={setEditVariableOpen}
          variableHead={variableHead!}
        />
      )}
      {editVariableValueOpen && (
        <EditVariableValueDialog
          openState={editVariableValueOpen}
          setOpenState={setEditVariableValueOpen}
          variableHead={variableHead!}
          environmentId={editVariableValueEnvironmentId!}
        />
      )}
      {deleteVariableOpen && (
        <DeleteVariableDialog
          openState={deleteVariableOpen}
          setOpenState={setDeleteVariableOpen}
          variableHead={variableHead!}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbCurrent>Variables</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <AddVariableDialog />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col" className="w-[10px]">
                <span className="sr-only">Actions</span>
              </TableHead>
              <TableHead scope="col">Name</TableHead>
              <TableHead scope="col">Secret</TableHead>
              <TableHead scope="col">Values</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {variables?.totalRows === 0 && <TableNoData colSpan={4} />}
            {variables?.items &&
              variables?.items.map((item) => (
                <TableRow key={item.name}>
                  <TableCell>
                    <TableButtonEdit onClick={() => handleEditVariable(item)} />
                    <TableButtonDelete
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteVariable(item)
                      }}
                    />
                  </TableCell>
                  <TableCell>{item.name}</TableCell>
                  <TableCell>
                    <Checkbox
                      checked={item.isSecret}
                      aria-readonly
                      className="cursor-default"
                    />
                  </TableCell>
                  <TableCell>
                    {Object.keys(item.values!).map((environmentId) => (
                      <VariableValue
                        key={environmentId}
                        envName={environmentsMap![environmentId]}
                        isSecret={item.isSecret!}
                        value={item.values![environmentId]}
                        onClick={() => {
                          handleEditVariableValue(item, environmentId)
                        }}
                      />
                    ))}
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </MainContent>
    </MainArea>
  )
}

function truncateVariableValue(value: string, chars = 20) {
  if (value && value.length > chars) {
    return value.substring(0, chars) + " ..."
  }
  return value
}

function VariableValue({
  envName,
  value,
  isSecret,
  onClick,
}: {
  envName: string
  value: string
  isSecret: boolean
  onClick: React.MouseEventHandler
}) {
  return (
    <div
      key={envName}
      className="mr-4 inline-flex cursor-pointer items-center rounded-md text-xs font-medium text-gray-600 ring-1 ring-inset ring-gray-500/40 hover:shadow-md dark:text-gray-50 dark:ring-gray-500 dark:hover:underline dark:hover:underline-offset-2"
      onClick={onClick}
    >
      <span className="rounded-l-md bg-gray-100 px-3 py-2 font-bold ring-1 ring-inset ring-gray-500/40 dark:bg-slate-700 dark:ring-gray-500">
        {envName}
      </span>
      <span className="px-3 text-slate-900 dark:text-slate-50">
        {value === null ? (
          <i>Unspecified</i>
        ) : isSecret ? (
          "*****"
        ) : (
          truncateVariableValue(value)
        )}
      </span>
    </div>
  )
}
