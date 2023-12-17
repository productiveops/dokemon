import { INodeComposeItemHead } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useComposeItem(projectName: string) {
  const url = `${apiBaseUrl()}/compose/${projectName}`

  const { data, error, isLoading, mutate } = useRequest<INodeComposeItemHead>({
    url,
  })

  const mutateComposeItem = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    composeItem: data,
    mutateComposeItem,
  }
}
