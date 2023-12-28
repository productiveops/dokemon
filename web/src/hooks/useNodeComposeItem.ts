import { INodeComposeItem } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useNodeComposeItem(
  nodeId: string,
  composeProjectId: string
) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}/compose/${composeProjectId}`

  const { data, error, isLoading, mutate } = useRequest<INodeComposeItem>(
    {
      url,
    },
    { revalidateOnFocus: false }
  )

  const mutateNodeComposeItem = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    nodeComposeItem: data,
    mutateNodeComposeItem,
  }
}
