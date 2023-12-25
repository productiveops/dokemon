import { TableCell, TableRow } from "../ui/table"

export function TableNoData({ colSpan }: { colSpan: number }) {
  return (
    <TableRow>
      <TableCell colSpan={colSpan} className="text-center">
        No data to display
      </TableCell>
    </TableRow>
  )
}
