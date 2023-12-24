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
import AddGitHubPATDialog from "./dialogs/add-github-pat-dialog"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { ICredentialHead } from "@/lib/api-models"
import DeleteCredentialDialog from "./dialogs/delete-credential-dialog"
import EditGithubPATDetailsDialog from "./dialogs/edit-github-pat-details-dialog"
import { CLASSES_CLICKABLE_TABLE_ROW } from "@/lib/utils"
import EditGithubPATSecretDialog from "./dialogs/edit-github-pat-secret-dialog"
import TableButtonDelete from "@/components/widgets/table-button-delete"
import { TableNoData } from "@/components/widgets/table-no-data"

export default function Credentials() {
  const { isLoading, credentials } = useCredentials()
  const [editCredentialOpen, setEditCredentialOpen] = useState(false)
  const [editCredentialSecretOpen, setEditCredentialSecretOpen] =
    useState(false)
  const [deleteCredentialOpen, setDeleteCredentialOpen] = useState(false)
  const [credentialHead, setCredentialHead] = useState<ICredentialHead | null>(
    null
  )

  if (isLoading) return <Loading />

  const handleEditCredential = (credentialHead: ICredentialHead) => {
    setCredentialHead({ ...credentialHead })
    setEditCredentialOpen(true)
  }

  const handleDeleteCredential = (credentialHead: ICredentialHead) => {
    setCredentialHead({ ...credentialHead })
    setDeleteCredentialOpen(true)
  }

  return (
    <MainArea>
      {editCredentialSecretOpen && (
        <EditGithubPATSecretDialog
          openState={editCredentialSecretOpen}
          setOpenState={setEditCredentialSecretOpen}
          credentialHead={credentialHead!}
        />
      )}
      {editCredentialOpen && (
        <EditGithubPATDetailsDialog
          openState={editCredentialOpen}
          setOpenState={setEditCredentialOpen}
          credentialHead={credentialHead!}
        />
      )}
      {deleteCredentialOpen && (
        <DeleteCredentialDialog
          openState={deleteCredentialOpen}
          setOpenState={setDeleteCredentialOpen}
          credentialHead={credentialHead!}
        />
      )}
      <TopBar>
        <Breadcrumb>
          <BreadcrumbCurrent>Credentials</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <AddGitHubPATDialog buttonCaption="Add GitHub Token" />
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
                        handleDeleteCredential(item)
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
