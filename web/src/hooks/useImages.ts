import { IImage, IPageResponse } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useImages(nodeId: string) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/images`

  const { data, error, isLoading, mutate } = useRequest<IPageResponse<IImage>>({
    url,
  })

  const mutateImages = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    images: data,
    mutateImages,
  }
}
