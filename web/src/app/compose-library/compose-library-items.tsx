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
import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import useComposeLibraryItemList from "@/hooks/useComposeLibraryItemList"
import { CLASSES_CLICKABLE_TABLE_ROW } from "@/lib/utils"
import { TableNoData } from "@/components/widgets/table-no-data"

export default function ComposeLibraryItems() {
  const navigate = useNavigate()
  const { isLoading, composeLibraryItems } = useComposeLibraryItemList()

  if (isLoading) return <Loading />

  return (
    <MainArea>
      <TopBar>
        <Breadcrumb>
          <BreadcrumbCurrent>Compose Library</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions>
          <Button
            className="w-24"
            onClick={() => navigate("/composelibrary/filesystem/create")}
          >
            Create
          </Button>
          <Button
            className="w-36"
            onClick={() => navigate("/composelibrary/github/create")}
          >
            Add from GitHub
          </Button>
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Library Project Name</TableHead>
              <TableHead scope="col">Type</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {composeLibraryItems?.totalRows === -1 && (
              <TableNoData colSpan={3} />
            )}
            {composeLibraryItems?.items &&
              composeLibraryItems?.items.map((item) => (
                <TableRow
                  key={item.projectName}
                  className={CLASSES_CLICKABLE_TABLE_ROW}
                  onClick={() => {
                    if (item.type === "filesystem") {
                      navigate(
                        `/composelibrary/${item.type}/${item.projectName}/edit`
                      )
                    }
                    if (item.type === "github") {
                      navigate(`/composelibrary/${item.type}/${item.id}/edit`)
                    }
                  }}
                >
                  <TableCell>{item.projectName}</TableCell>
                  <TableCell>
                    {item.type === "filesystem" ? "File System" : ""}
                    {item.type === "github" ? "GitHub" : ""}
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </MainContent>
    </MainArea>
  )
}
