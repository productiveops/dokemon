import { INodeHead } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useNodeHead(nodeId: string) {
  const url = `${apiBaseUrl()}/nodes/head/${nodeId}`

  const { data, error, isLoading, mutate } = useRequest<INodeHead>({
    url,
  })

  const mutateNodeHead = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    nodeHead: data,
    mutateNodeHead,
  }
}
