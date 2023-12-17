import { INetwork, IPageResponse } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useNetworks(nodeId: string) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/networks`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<INetwork>
  >({
    url,
  })

  const mutateNetworks = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    networks: data,
    mutateNetworks,
  }
}
