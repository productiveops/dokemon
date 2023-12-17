import apiBaseUrl from "@/lib/api-base-url"
import { IPageResponse, IVariableHead } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useVariables() {
  const url = `${apiBaseUrl()}/variables?p=1&s=1000`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<IVariableHead>
  >({
    url,
  })

  const mutateVariables = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    variables: data,
    mutateVariables,
  }
}
