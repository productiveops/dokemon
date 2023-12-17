import apiBaseUrl from "@/lib/api-base-url"
import { IEnvironment } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useEnvironment(environmentId: string) {
  const url = `${apiBaseUrl()}/environments/${environmentId}`

  const { data, error, isLoading, mutate } = useRequest<IEnvironment>({
    url,
  })

  const mutateEnvironment = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    environment: data,
    mutateEnvironment,
  }
}
