import { INodeComposeContainer, IPageResponse } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useNodeComposeContainers(
  nodeId: string,
  composeProjectId: string
) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/compose/${composeProjectId}/containers`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<INodeComposeContainer>
  >({
    url,
  })

  const mutateComposeContainers = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    composeContainers: data,
    mutateComposeContainers,
  }
}
