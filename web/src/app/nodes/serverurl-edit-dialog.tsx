import { useMemo, useState } from "react"
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
  FormMessage,
} from "@/components/ui/form"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { SubmitHandler, useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import {
  cn,
  toastSomethingWentWrong,
  toastSuccess,
  trimString,
} from "@/lib/utils"
import useNodes from "@/hooks/useNodes"
import useSetting from "@/hooks/useSetting"
import { ISettingUpdateRequest } from "@/lib/api-models"
import { apiSettingsUpdate } from "@/lib/api"
import SpinnerIcon from "@/components/widgets/spinner-icon"

export default function ServerUrlEditDialog() {
  const { setting, mutateSetting } = useSetting("SERVER_URL")
  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const { mutateNodes } = useNodes()

  const formSchema = z.object({
    value: z.preprocess(
      trimString,
      z.string().url("Invalid URL format").min(1, "Value is required").max(200)
    ),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: useMemo(() => {
      if (setting?.value) {
        return { value: setting?.value }
      }

      return { value: `${location.protocol}//${location.host}` }
    }, [setting]),
  })

  const handleCloseForm = () => {
    setOpen(false)
    form.reset()
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (
    data: ISettingUpdateRequest
  ) => {
    setIsSaving(true)
    const response = await apiSettingsUpdate("SERVER_URL", data)
    if (!response.ok) {
      handleCloseForm()
      toastSomethingWentWrong("There was a problem saving the URL. Try again!")
    } else {
      mutateNodes()
      mutateSetting()
      handleCloseForm()
      toastSuccess("Server URL has been saved.")
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>{setting?.value ? "Edit" : "Set"} Server URL</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <fieldset className={cn("group")} disabled={isSaving}>
              <DialogHeader>
                <DialogTitle>Server URL</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4 group-disabled:opacity-50">
                <FormField
                  control={form.control}
                  name="value"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <Input {...field} autoFocus />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <p className="text-sm">
                  This URL should be accessible from machines which you would
                  like to manage from Dokemon.
                </p>
              </div>
              <DialogFooter>
                <Button
                  className={cn(
                    "relative w-24 group-disabled:pointer-events-none"
                  )}
                  type="submit"
                >
                  <SpinnerIcon />
                  <span className={cn("group-disabled:opacity-0")}>
                    {setting?.value ? "Save" : "Set"}
                  </span>
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
