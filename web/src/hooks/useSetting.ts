import { ISetting } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useSetting(settingId: string) {
  const url = `${apiBaseUrl()}/settings/${settingId}`

  const { data, error, isLoading, mutate } = useRequest<ISetting>({
    url,
  })

  const mutateSetting = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    setting: data,
    mutateSetting,
  }
}
