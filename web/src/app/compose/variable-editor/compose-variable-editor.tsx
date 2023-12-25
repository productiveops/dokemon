import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import useNodeComposeVariables from "@/hooks/useNodeComposeVariables"
import ComposeVariableAddDialog from "./compose-variable-add-dialog"
import { Checkbox } from "@/components/ui/checkbox"
import { Button } from "@/components/ui/button"
import { PencilIcon, TrashIcon } from "@heroicons/react/24/outline"
import { CLASSES_TABLE_ACTION_ICON } from "@/lib/utils"
import { INodeComposeVariable } from "@/lib/api-models"
import { useState } from "react"
import ComposeVariableEditDialog from "./compose-variable-edit-dialog"
import apiBaseUrl from "@/lib/api-base-url"
import DeleteDialog from "@/components/delete-dialog"

export default function ComposeVariableEditor({
  nodeId,
  nodeComposeProjectId,
}: {
  nodeId: string
  nodeComposeProjectId: string
}) {
  const { nodeComposeVariables } = useNodeComposeVariables(
    nodeId!,
    nodeComposeProjectId
  )
  const { mutateNodeComposeVariables } = useNodeComposeVariables(
    nodeId,
    nodeComposeProjectId
  )

  const [editVariableOpen, setEditVariableOpen] = useState(false)
  const [variable, setVariable] = useState<INodeComposeVariable | null>(null)

  function handleEditVariable(item: INodeComposeVariable) {
    setVariable({ ...item })
    setEditVariableOpen(true)
  }

  const [deleteVariableConfirmationOpen, setDeleteVariableConfirmationOpen] =
    useState(false)
  function handleDeleteVariableConfirmation(item: INodeComposeVariable) {
    setVariable({ ...item })
    setDeleteVariableConfirmationOpen(true)
  }

  const [deleteInProgress, setDeleteInProgress] = useState(false)
  const handleDeleteVariable = async () => {
    setDeleteInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/compose/${nodeComposeProjectId}/variables/${variable?.id}`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      setDeleteVariableConfirmationOpen(false)
    } else {
      mutateNodeComposeVariables()
      setTimeout(() => {
        setDeleteVariableConfirmationOpen(false)
      }, 500)
    }
    setDeleteInProgress(false)
  }

  return (
    <div>
      {editVariableOpen && (
        <ComposeVariableEditDialog
          openState={editVariableOpen}
          setOpenState={setEditVariableOpen}
          variable={variable!}
          nodeId={nodeId}
          nodeComposeProjectId={nodeComposeProjectId}
        />
      )}
      {deleteVariableConfirmationOpen && (
        <DeleteDialog
          openState={deleteVariableConfirmationOpen}
          setOpenState={setDeleteVariableConfirmationOpen}
          deleteCaption=""
          deleteHandler={handleDeleteVariable}
          isProcessing={deleteInProgress}
          title="Delete Variable"
          message={`Are you sure you want to delete variable '${variable?.name}?'`}
        />
      )}
      <h2>
        Variables{" "}
        <span className="float-right text-xs">
          * Changes made to variables are saved immediately
        </span>
      </h2>
      <Table className="mb-2">
        <TableHeader>
          <TableRow>
            <TableHead className="sm:pl-0" scope="col">
              Name
            </TableHead>
            <TableHead scope="col">Secret</TableHead>
            <TableHead scope="col">Value</TableHead>
            <TableHead scope="col">
              <span className="sr-only">Actions</span>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {nodeComposeVariables?.items?.map((item) => (
            <TableRow key={item.id}>
              <TableCell className="sm:pl-0">{item.name}</TableCell>
              <TableCell>
                <Checkbox checked={item.isSecret} disabled />
              </TableCell>
              <TableCell>{item.isSecret ? "*****" : item.value}</TableCell>
              <TableCell className="text-right">
                <Button
                  variant="ghost"
                  size={"sm"}
                  title="Delete"
                  onClick={(e) => {
                    e.stopPropagation()
                    handleEditVariable(item)
                  }}
                >
                  <PencilIcon className={CLASSES_TABLE_ACTION_ICON} />
                </Button>
                <Button
                  variant="ghost"
                  size={"sm"}
                  title="Delete"
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDeleteVariableConfirmation(item)
                  }}
                >
                  <TrashIcon className={CLASSES_TABLE_ACTION_ICON} />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <ComposeVariableAddDialog
        nodeId={nodeId}
        nodeComposeProjectId={nodeComposeProjectId}
      />
    </div>
  )
}
