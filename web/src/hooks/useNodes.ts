import { IPageResponse, INodeHead } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useNodes() {
  const url = `${apiBaseUrl()}/nodes?p=1&s=1000`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<INodeHead>
  >({
    url,
  })

  const mutateNodes = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    nodes: data,
    mutateNodes,
  }
}
