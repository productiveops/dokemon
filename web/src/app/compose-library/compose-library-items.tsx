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
            onClick={() => navigate("/composelibrary/create")}
          >
            Create
          </Button>
        </TopBarActions>
      </TopBar>
      <MainContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead scope="col">Name</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!composeLibraryItems?.items && (
              <TableRow>
                <TableCell colSpan={5} className="text-center">
                  No data to display
                </TableCell>
              </TableRow>
            )}
            {composeLibraryItems?.items &&
              composeLibraryItems?.items.map((item) => (
                <TableRow
                  key={item.projectName}
                  className="cursor-pointer hover:bg-slate-50"
                  onClick={() => {
                    let navigateTo = "edit"
                    navigate(
                      `/composelibrary/${item.projectName}/${navigateTo}`
                    )
                  }}
                >
                  <TableCell>{item.projectName}</TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </MainContent>
    </MainArea>
  )
}
