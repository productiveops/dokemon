import { IComposeLibraryItemHead, IPageResponse } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useComposeLibraryItemList() {
  const url = `${apiBaseUrl()}/composelibrary?p=1&s=1000`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<IComposeLibraryItemHead>
  >({
    url,
  })

  const mutateComposeLibraryItemList = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    composeLibraryItems: data,
    mutateComposeLibraryItemList,
  }
}
