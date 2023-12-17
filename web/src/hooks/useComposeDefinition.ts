import { IComposeDefinition } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useComposeDefinition(projectName: string) {
  const url = `${apiBaseUrl()}/compose/${projectName}/definition`

  const { data, error, isLoading, mutate } = useRequest<IComposeDefinition>({
    url,
  })

  const mutateComposeDefinition = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    composeDefinition: data,
    mutateComposeDefinition,
  }
}
