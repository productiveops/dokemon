import { IContainer, IPageResponse } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useContainers(nodeId: string) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/containers`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<IContainer>
  >({
    url,
  })

  const mutateContainers = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    containers: data,
    mutateContainers,
  }
}
