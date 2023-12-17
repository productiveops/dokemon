import { INodeComposeItemHead, IPageResponse } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useNodeComposeList(nodeId: string) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/compose?p=1&s=1000`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<INodeComposeItemHead>
  >({
    url,
  })

  const mutateNodeComposeList = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    nodeComposeItems: data,
    mutateNodeComposeList,
  }
}
