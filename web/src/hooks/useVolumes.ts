import { IPageResponse, IVolume } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useVolumes(nodeId: string) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/volumes`

  const { data, error, isLoading, mutate } = useRequest<IPageResponse<IVolume>>(
    { url }
  )

  const mutateVolumes = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    volumes: data,
    mutateVolumes,
  }
}
