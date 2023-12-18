import { ICredential } from "@/lib/api-models"
import apiBaseUrl from "@/lib/api-base-url"
import useRequest from "@/lib/useRequest"

export default function useCredential(credentialId: number) {
  const url = `${apiBaseUrl()}/credentials/${credentialId}`

  const { data, error, isLoading, mutate } = useRequest<ICredential>({
    url,
  })

  const mutateCredential = async () => {
    mutate()
  }

  return {
    isLoading,
    isError: error,
    credential: data,
    mutateCredential,
  }
}
