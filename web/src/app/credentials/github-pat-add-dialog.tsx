import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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
  toastFailed,
  toastSuccess,
  trimString,
} from "@/lib/utils"
import useCredentials from "@/hooks/useCredentials"
import apiBaseUrl from "@/lib/api-base-url"
import SpinnerIcon from "@/components/widgets/spinner-icon"

export default function GitHubPATAddDialog({
  buttonCaption,
  disabled,
}: {
  buttonCaption: string
  disabled?: boolean
}) {
  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const { mutateCredentials } = useCredentials()

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
              `${apiBaseUrl()}/credentials/uniquename?value=${value}`
            ),
          "Another credential with this name already exists"
        )
    ),
    secret: z.string().min(1, "Required").max(200),
    service: z.string().default("github"),
    type: z.string().default("pat"),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      secret: "",
      service: "github",
      type: "pat",
    },
  })

  const handleCloseForm = () => {
    setOpen(false)
    form.reset()
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)
    const response = await fetch(`${apiBaseUrl()}/credentials`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })
    if (!response.ok) {
      handleCloseForm()
      toastFailed(
        "There was a problem when creating new credential. Try again!"
      )
    } else {
      mutateCredentials()
      handleCloseForm()
      toastSuccess("New credential has been added.")
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button disabled={disabled !== undefined || disabled}>
          {buttonCaption}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <fieldset className={cn("group")} disabled={isSaving}>
              <DialogHeader>
                <DialogTitle>Add GitHub Personal Access Token</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4 group-disabled:opacity-50">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Credential Name</FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="secret"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>GitHub Personal Access Token</FormLabel>
                      <FormControl>
                        <Input {...field} type="password" />
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
