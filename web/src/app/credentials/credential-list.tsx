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
import useCredentials from "@/hooks/useCredentials"
import GitHubPATAddDialog from "./github-pat-add-dialog"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { ICredentialHead } from "@/lib/api-models"
import GithubPATDetailsEditDialog from "./github-pat-details-edit-dialog"
import {
  CLASSES_CLICKABLE_TABLE_ROW,
  toastFailed,
  toastSuccess,
} from "@/lib/utils"
import GithubPATSecretEditDialog from "./github-pat-secret-edit-dialog"
import TableButtonDelete from "@/components/widgets/table-button-delete"
import { TableNoData } from "@/components/widgets/table-no-data"
import apiBaseUrl from "@/lib/api-base-url"
import DeleteDialog from "@/components/delete-dialog"

export default function CredentialList() {
  const { isLoading, credentials, mutateCredentials } = useCredentials()
  const [editCredentialOpen, setEditCredentialOpen] = useState(false)
  const [editCredentialSecretOpen, setEditCredentialSecretOpen] =
    useState(false)
  const [
    deleteCredentialConfirmationOpen,
    setDeleteCredentialConfirmationOpen,
  ] = useState(false)
  const [credentialHead, setCredentialHead] = useState<ICredentialHead | null>(
    null
  )
  const [deleteInProgress, setDeleteInProgress] = useState(false)

  if (isLoading) return <Loading />

  const handleEditCredential = (credentialHead: ICredentialHead) => {
    setCredentialHead({ ...credentialHead })
    setEditCredentialOpen(true)
  }

  const handleDeleteCredentialConfirmation = (
    credentialHead: ICredentialHead
  ) => {
    setCredentialHead({ ...credentialHead })
    setDeleteCredentialConfirmationOpen(true)
  }

  const handleDelete = async () => {
    setDeleteInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/credentials/${credentialHead?.id}`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      const r = await response.json()
      setDeleteCredentialConfirmationOpen(false)
      toastFailed(r.errors?.body)
    } else {
      mutateCredentials()
      setTimeout(() => {
        setDeleteCredentialConfirmationOpen(false)
        toastSuccess("Credential deleted.")
      }, 500)
    }
    setDeleteInProgress(false)
  }

  return (
    <MainArea>
      {editCredentialSecretOpen && (
        <GithubPATSecretEditDialog
          openState={editCredentialSecretOpen}
          setOpenState={setEditCredentialSecretOpen}
          credentialHead={credentialHead!}
        />
      )}
      {editCredentialOpen && (
        <GithubPATDetailsEditDialog
          openState={editCredentialOpen}
          setOpenState={setEditCredentialOpen}
          credentialHead={credentialHead!}
        />
      )}
      {deleteCredentialConfirmationOpen && (
        <DeleteDialog
          openState={deleteCredentialConfirmationOpen}
          setOpenState={setDeleteCredentialConfirmationOpen}
          deleteCaption=""
          deleteHandler={handleDelete}
          isProcessing={deleteInProgress}
          title="Delete Credentials"
          message={`Are you sure you want to delete credentials '${credentialHead?.name}?'`}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbCurrent>Credentials</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <GitHubPATAddDialog buttonCaption="Add GitHub Token" />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Name</TableHead>
              <TableHead scope="col">Type</TableHead>
              <TableHead scope="col">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {credentials?.totalRows === 0 && <TableNoData colSpan={3} />}
            {credentials?.items &&
              credentials?.items.map((item) => (
                <TableRow
                  key={item.name}
                  className={CLASSES_CLICKABLE_TABLE_ROW}
                  onClick={() => {
                    handleEditCredential(item)
                  }}
                >
                  <TableCell>{item.name}</TableCell>
                  <TableCell>
                    {item.service === "github" ? "GitHub" : ""}{" "}
                    {item.type === "pat" ? "Personal Access Token" : ""}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      className="mr-4"
                      size={"sm"}
                      onClick={(e) => {
                        e.stopPropagation()
                        setCredentialHead(item)
                        setEditCredentialSecretOpen(true)
                      }}
                    >
                      Update Token
                    </Button>
                    <TableButtonDelete
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteCredentialConfirmation(item)
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
