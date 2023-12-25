import { TrashIcon } from "@heroicons/react/24/solid"
import { Button } from "../ui/button"
import { CLASSES_TABLE_ACTION_ICON } from "@/lib/utils"

export default function TableButtonDelete({
  onClick,
}: {
  onClick: (e: React.MouseEvent<HTMLButtonElement>) => void
}) {
  return (
    <Button variant="ghost" size={"sm"} title="Delete" onClick={onClick}>
      <TrashIcon className={CLASSES_TABLE_ACTION_ICON} />
    </Button>
  )
}
