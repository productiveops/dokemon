import { ChevronRightIcon } from "@heroicons/react/24/outline"
import { Link } from "react-router-dom"

export function Breadcrumb({ children }: { children: any }) {
  return <div className="text-sm sm:flex-auto">{children}</div>
}

export function BreadcrumbLink({
  to,
  children,
}: {
  to: string
  children: any
}) {
  return (
    <Link to={to} className="text-amber-600">
      {children}
    </Link>
  )
}

export function BreadcrumbSeparator() {
  return <ChevronRightIcon className="mx-1 inline w-3" />
}

export function BreadcrumbCurrent({ children }: { children: any }) {
  return <span>{children}</span>
}
