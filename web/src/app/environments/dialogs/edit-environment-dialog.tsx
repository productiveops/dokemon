import { Dispatch, SetStateAction, useState } from "react"
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
  FormMessage,
} from "@/components/ui/form"
import { Button } from "@/components/ui/button"
import { Icons } from "@/components/icons"
import { Input } from "@/components/ui/input"
import { SubmitHandler, useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { toast } from "@/components/ui/use-toast"
import {
  REGEX_IDENTIFIER,
  REGEX_IDENTIFIER_MESSAGE,
  cn,
  trimString,
} from "@/lib/utils"
import useEnvironments from "@/hooks/useEnvironments"
import apiBaseUrl from "@/lib/api-base-url"
import { IEnvironmentHead } from "@/lib/api-models"

export default function EditEnvironmentDialog({
  openState,
  setOpenState,
  environmentHead,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  environmentHead: IEnvironmentHead
}) {
  const [isSaving, setIsSaving] = useState(false)
  const { mutateEnvironments } = useEnvironments()

  const formSchema = z.object({
    name: z.preprocess(
      trimString,
      z
        .string()
        .min(1, "Name is required")
        .max(20)
        .regex(REGEX_IDENTIFIER, REGEX_IDENTIFIER_MESSAGE)
        .refine(async (value) => {
          const res = await fetch(
            `${apiBaseUrl()}/environments/${
              environmentHead.id
            }/uniquename?value=${value}`
          )
          return (await res.json()).unique
        }, "Another environment with this name already exists")
    ),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: environmentHead,
  })

  const handleCloseForm = () => {
    setOpenState(false)
    form.reset()
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)
    const response = await fetch(
      `${apiBaseUrl()}/environments/${environmentHead.id}`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    )
    if (!response.ok) {
      handleCloseForm()
      toast({
        variant: "destructive",
        title: "Something went wrong.",
        description: "There was a problem saving the environment. Try again!",
      })
    } else {
      mutateEnvironments()
      handleCloseForm()
      toast({
        title: "Success!",
        description: "Environment has been saved.",
      })
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
                <DialogTitle>Edit Environment</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4 group-disabled:opacity-50">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
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
                  <Icons.spinner
                    className={cn(
                      "absolute animate-spin text-slate-100 group-enabled:opacity-0"
                    )}
                  />
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
