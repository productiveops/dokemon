import apiBaseUrl from "@/lib/api-base-url"
import { IComposeLibraryItem } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useComposeLibraryItem(composeProjectName: string) {
  const url = `${apiBaseUrl()}/composelibrary/${composeProjectName}`

  const { data, error, isLoading, mutate } = useRequest<IComposeLibraryItem>(
    {
      url,
    },
    { revalidateOnFocus: false }
  )

  const mutateComposeLibraryItem = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    composeLibraryItem: data,
    mutateComposeLibraryItem,
  }
}
