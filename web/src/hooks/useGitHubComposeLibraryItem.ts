import apiBaseUrl from "@/lib/api-base-url"
import { IGitHubComposeLibraryItem } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useGitHubComposeLibraryItem(id: string) {
  const url = `${apiBaseUrl()}/composelibrary/github/${id}`

  const { data, error, isLoading, mutate } =
    useRequest<IGitHubComposeLibraryItem>(
      {
        url,
      },
      { revalidateOnFocus: false }
    )

  const mutateGitHubComposeLibraryItem = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    gitHubComposeLibraryItem: data,
    mutateGitHubComposeLibraryItem,
  }
}
