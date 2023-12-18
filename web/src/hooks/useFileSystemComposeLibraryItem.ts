import apiBaseUrl from "@/lib/api-base-url"
import { ILocalComposeLibraryItem } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useFileSystemComposeLibraryItem(
  composeProjectName: string
) {
  const url = `${apiBaseUrl()}/composelibrary/filesystem/${composeProjectName}`

  const { data, error, isLoading, mutate } =
    useRequest<ILocalComposeLibraryItem>(
      {
        url,
      },
      { revalidateOnFocus: false }
    )

  const mutateLocalComposeLibraryItem = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    localComposeLibraryItem: data,
    mutateLocalComposeLibraryItem,
  }
}
