import apiBaseUrl from "@/lib/api-base-url"
import { IVariable } from "@/lib/api-models"
import useRequest from "@/lib/useRequest"

export default function useVariable(variableId: string) {
  const url = `${apiBaseUrl()}/variables/${variableId}`

  const { data, error, isLoading, mutate } = useRequest<IVariable>({
    url,
  })

  const mutateVariable = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    variable: data,
    mutateVariable,
  }
}
