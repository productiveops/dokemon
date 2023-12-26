import { cn } from "@/lib/utils"

export function MainContainer({
  visible = true,
  children,
}: {
  visible?: boolean
  children: any
}) {
  return (
    <div
      className={cn(
        "bg-white shadow-sm ring-1 ring-gray-900/5 dark:bg-gray-900 dark:ring-gray-800 sm:rounded-xl md:col-span-2",
        visible ? "" : "hidden"
      )}
    >
      {children}
    </div>
  )
}

export function Section({ children }: { children: any }) {
  return <div className="px-4 py-6 sm:p-8">{children}</div>
}

export function SectionHead({ children }: { children: any }) {
  return (
    <div className="mb-4 sm:divide-y sm:divide-gray-900/10 sm:border-b sm:pb-4">
      <h2 className="text-base font-semibold leading-7 text-gray-900">
        {children}
      </h2>
    </div>
  )
}

export function SectionBody({
  className,
  children,
}: {
  className?: string
  children: any
}) {
  return (
    <div
      className={cn(
        "grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-1",
        className
      )}
    >
      {children}
    </div>
  )
}
