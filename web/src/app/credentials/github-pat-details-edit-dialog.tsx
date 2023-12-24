import { Dispatch, SetStateAction, useEffect, useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { SubmitHandler, useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import {
  REGEX_IDENTIFIER,
  REGEX_IDENTIFIER_MESSAGE,
  cn,
  hasUniqueName,
  toastSomethingWentWrong,
  toastSuccess,
  trimString,
} from "@/lib/utils"
import useCredentials from "@/hooks/useCredentials"
import apiBaseUrl from "@/lib/api-base-url"
import { ICredentialHead } from "@/lib/api-models"
import useCredential from "@/hooks/useCredential"
import SpinnerIcon from "@/components/widgets/spinner-icon"

export default function GithubPATDetailsEditDialog({
  openState,
  setOpenState,
  credentialHead,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  credentialHead: ICredentialHead
}) {
  const [isSaving, setIsSaving] = useState(false)
  const { mutateCredentials } = useCredentials()
  const { credential, mutateCredential } = useCredential(credentialHead.id)

  const formSchema = z.object({
    name: z.preprocess(
      trimString,
      z
        .string()
        .min(1, "Required")
        .max(20)
        .regex(REGEX_IDENTIFIER, REGEX_IDENTIFIER_MESSAGE)
        .refine(
          async (value) =>
            hasUniqueName(
              `${apiBaseUrl()}/credentials/${
                credentialHead.id
              }/uniquename?value=${value}`
            ),
          "Another credential with this name already exists"
        )
    ),
    service: z.string().default("github"),
    type: z.string().default("pat"),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: { name: "", service: "github", type: "pat" },
  })

  useEffect(() => {
    if (credential) form.reset(credential)
  }, [credential])

  const handleCloseForm = () => {
    setOpenState(false)
    form.reset()
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)
    const response = await fetch(
      `${apiBaseUrl()}/credentials/${credentialHead.id}`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    )
    if (!response.ok) {
      handleCloseForm()
      toastSomethingWentWrong(
        "There was a problem saving the credential. Try again!"
      )
    } else {
      mutateCredential()
      mutateCredentials()
      handleCloseForm()
      toastSuccess("Credential has been saved.")
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={openState} onOpenChange={setOpenState}>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <fieldset className={cn("group")} disabled={isSaving}>
              <DialogHeader>
                <DialogTitle>Edit GitHub Credential Details</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4 group-disabled:opacity-50">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Credential Name</FormLabel>
                      <FormControl>
                        <Input {...field} autoFocus />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <DialogFooter>
                <Button
                  className={cn(
                    "relative w-24 group-disabled:pointer-events-none"
                  )}
                  type="submit"
                >
                  <SpinnerIcon />
                  <span className={cn("group-disabled:opacity-0")}>Save</span>
                </Button>
                <Button
                  type="button"
                  className="w-24"
                  variant={"secondary"}
                  onClick={handleCloseForm}
                >
                  Cancel
                </Button>
              </DialogFooter>
            </fieldset>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
