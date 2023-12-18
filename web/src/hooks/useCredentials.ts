import { IPageResponse, ICredentialHead } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useCredentials() {
  const url = `${apiBaseUrl()}/credentials?p=1&s=1000`

  const { data, error, isLoading, mutate } = useRequest<
    IPageResponse<ICredentialHead>
  >({
    url,
  })

  const mutateCredentials = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    credentials: data,
    mutateCredentials,
  }
}
