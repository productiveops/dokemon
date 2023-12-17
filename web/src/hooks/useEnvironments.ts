import { IPageResponse, IEnvironmentHead } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useEnvironments() {
  const url = `${apiBaseUrl()}/environments?p=1&s=1000`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<IEnvironmentHead>
  >({
    url,
  })

  const mutateEnvironments = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    environments: data,
    mutateEnvironments,
  }
}
