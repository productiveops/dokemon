import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"
import { mutate } from "swr"

export default function useEnvironmentsMap() {
  const url = `${apiBaseUrl()}/environments/map`

  const { data, error, isLoading } = useRequest<{ [id: string]: string }>({
    url,
  })

  const mutateEnvironmentsMap = async () => {
    mutate(url)
  }

  return {
    isLoading,
    isError: error,
    environmentsMap: data,
    mutateEnvironmentsMap,
  }
}
