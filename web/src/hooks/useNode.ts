import { INode } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useNode(nodeId: string) {
  const url = `${apiBaseUrl()}/nodes/${nodeId}`

  const { data, error, isLoading, mutate } = useRequest<INode>({
    url,
  })

  const mutateNode = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    node: data,
    mutateNode,
  }
}
