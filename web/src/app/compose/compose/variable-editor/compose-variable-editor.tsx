import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import useNodeComposeVariables from "@/hooks/useNodeComposeVariables"
import AddComposeVariableDialog from "./add-compose-variable-dialog"
import { Checkbox } from "@/components/ui/checkbox"
import { Button } from "@/components/ui/button"
import { PencilIcon, TrashIcon } from "@heroicons/react/24/outline"
import { CLASSES_TABLE_ACTION_ICON } from "@/lib/utils"
import { INodeComposeVariable } from "@/lib/api-models"
import { useState } from "react"
import DeleteComposeVariableDialog from "./delete-compose-variable-dialog"
import EditComposeVariableDialog from "./edit-compose-variable-dialog"

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
  const [editVariableOpen, setEditVariableOpen] = useState(false)
  const [deleteVariableOpen, setDeleteVariableOpen] = useState(false)
  const [variable, setVariable] = useState<INodeComposeVariable | null>(null)

  function handleEditVariable(item: INodeComposeVariable) {
    setVariable({ ...item })
    setEditVariableOpen(true)
  }

  function handleDeleteVariable(item: INodeComposeVariable) {
    setVariable({ ...item })
    setDeleteVariableOpen(true)
  }

  return (
    <div>
      {editVariableOpen && (
        <EditComposeVariableDialog
          openState={editVariableOpen}
          setOpenState={setEditVariableOpen}
          variable={variable!}
          nodeId={nodeId}
          nodeComposeProjectId={nodeComposeProjectId}
        />
      )}
      {deleteVariableOpen && (
        <DeleteComposeVariableDialog
          openState={deleteVariableOpen}
          setOpenState={setDeleteVariableOpen}
          variable={variable!}
          nodeId={nodeId}
          nodeComposeProjectId={nodeComposeProjectId}
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
                    handleDeleteVariable(item)
                  }}
                >
                  <TrashIcon className={CLASSES_TABLE_ACTION_ICON} />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <AddComposeVariableDialog
        nodeId={nodeId}
        nodeComposeProjectId={nodeComposeProjectId}
      />
    </div>
  )
}
