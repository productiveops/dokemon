import { Button } from "../ui/button"
import { CLASSES_TABLE_ACTION_ICON } from "@/lib/utils"
import { PencilIcon } from "@heroicons/react/24/outline"

export default function TableButtonEdit({
  onClick,
}: {
  onClick: (e: React.MouseEvent<HTMLButtonElement>) => void
}) {
  return (
    <Button variant="ghost" size={"sm"} title="Edit" onClick={onClick}>
      <PencilIcon className={CLASSES_TABLE_ACTION_ICON} />
    </Button>
  )
}
