export default function TopBarActions({ children }: { children?: any }) {
  return (
    <div className="mt-4 flex gap-1 sm:ml-16 sm:mt-0 sm:flex-none">
      {children}
    </div>
  )
}
