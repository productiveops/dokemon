import apiBaseUrl from "@/lib/api-base-url"
import { INodeComposeVariable, IPageResponse } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useNodeComposeVariables(
  nodeId: string,
  nodeComposeProjectId: string
) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/compose/${nodeComposeProjectId}/variables?p=1&s=1000`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<INodeComposeVariable>
  >({
    url,
  })

  const mutateNodeComposeVariables = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    nodeComposeVariables: data,
    mutateNodeComposeVariables,
  }
}
