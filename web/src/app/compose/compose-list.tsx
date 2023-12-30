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
import MainArea from "@/components/widgets/main-area"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainContent from "@/components/widgets/main-content"
import useNodeComposeList from "@/hooks/useNodeComposeList"
import { Badge } from "@/components/ui/badge"
import { useNavigate, useParams } from "react-router-dom"
import useNodeHead from "@/hooks/useNodeHead"
import ComposeAddProjectDialog from "./compose-add-project-dialog"
import { ArrowUpRight } from "lucide-react"
import { CLASSES_CLICKABLE_TABLE_ROW } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { TableNoData } from "@/components/widgets/table-no-data"
import { StaleStatusIcon } from "../containers/container-list"

export default function ComposeList() {
  const { nodeId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const navigate = useNavigate()
  const { isLoading, nodeComposeItems } = useNodeComposeList(nodeId!)

  if (isLoading) return <Loading />

  return (
    <MainArea>
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeHead?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Compose</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <Button
            variant={"default"}
            onClick={() => {
              navigate(`/nodes/${nodeId}/compose/create/local`)
            }}
          >
            Create
          </Button>
          <Button
            variant={"default"}
            onClick={() => {
              navigate(`/nodes/${nodeId}/compose/create/github`)
            }}
          >
            Add from GitHub
          </Button>
          <ComposeAddProjectDialog />
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Project Name</TableHead>
              <TableHead scope="col">Type</TableHead>
              <TableHead scope="col">Library Project Name</TableHead>
              <TableHead scope="col">Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {nodeComposeItems?.totalRows === 0 && <TableNoData colSpan={5} />}
            {nodeComposeItems?.items &&
              nodeComposeItems?.items.map((item) => (
                <TableRow
                  key={item.projectName}
                  className={CLASSES_CLICKABLE_TABLE_ROW}
                  onClick={() => {
                    navigate(`/nodes/${nodeId}/compose/${item.id}/definition`)
                  }}
                >
                  <TableCell>
                    <span>
                      {item.status.startsWith("running") && (
                        <StaleStatusIcon status={item.stale} />
                      )}
                      {item.projectName}
                    </span>
                  </TableCell>
                  <TableCell>
                    {item.type === "local" ? "Local" : "GitHub"}
                  </TableCell>
                  <TableCell>
                    {!item.libraryProjectName ? (
                      "-"
                    ) : (
                      <span
                        className="p-1 text-amber-600 hover:underline"
                        onClick={(e) => {
                          e.stopPropagation()
                          if (item.libraryProjectId) {
                            window.open(
                              `${location.protocol}//${location.host}/composelibrary/github/${item.libraryProjectId}/edit`
                            )
                          } else {
                            window.open(
                              `${location.protocol}//${location.host}/composelibrary/filesystem/${item.libraryProjectName}/edit`
                            )
                          }
                        }}
                      >
                        {item.libraryProjectName}
                        <ArrowUpRight className="ml-1 inline w-4" />
                      </span>
                    )}
                  </TableCell>
                  <TableCell>
                    {item.status && item.status.startsWith("running") && (
                      <Badge variant="default">{item.status}</Badge>
                    )}
                    {item.status && !item.status.startsWith("running") && (
                      <Badge variant="destructive">{item.status}</Badge>
                    )}
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </MainContent>
    </MainArea>
  )
}
