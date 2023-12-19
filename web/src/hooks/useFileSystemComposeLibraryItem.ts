import apiBaseUrl from "@/lib/api-base-url"
import { IFileSystemComposeLibraryItem } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useFileSystemComposeLibraryItem(
  composeProjectName: string
) {
  const url = `${apiBaseUrl()}/composelibrary/filesystem/${composeProjectName}`

  const { data, error, isLoading, mutate } =
    useRequest<IFileSystemComposeLibraryItem>(
      {
        url,
      },
      { revalidateOnFocus: false }
    )

  const mutateFileSystemComposeLibraryItem = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    fileSystemComposeLibraryItem: data,
    mutateFileSystemComposeLibraryItem,
  }
}
