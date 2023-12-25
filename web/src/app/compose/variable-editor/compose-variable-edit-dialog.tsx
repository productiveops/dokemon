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
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { SubmitHandler, useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import {
  cn,
  hasUniqueName,
  toastFailed,
  toastSuccess,
  trimString,
} from "@/lib/utils"
import { Checkbox } from "@/components/ui/checkbox"
import apiBaseUrl from "@/lib/api-base-url"
import { INodeComposeVariable } from "@/lib/api-models"
import useNodeComposeVariables from "@/hooks/useNodeComposeVariables"
import SpinnerIcon from "@/components/widgets/spinner-icon"

export default function ComposeVariableEditDialog({
  openState,
  setOpenState,
  variable,
  nodeId,
  nodeComposeProjectId,
}: {
  openState: boolean
  setOpenState: Dispatch<SetStateAction<boolean>>
  variable: INodeComposeVariable
  nodeId: string
  nodeComposeProjectId: string
}) {
  const [isSaving, setIsSaving] = useState(false)
  const { mutateNodeComposeVariables } = useNodeComposeVariables(
    nodeId,
    nodeComposeProjectId
  )

  const formSchema = z.object({
    name: z.preprocess(
      trimString,
      z
        .string()
        .min(1, "Name is required")
        .max(100)
        .refine(
          async (value) =>
            hasUniqueName(
              `${apiBaseUrl()}/nodes/${nodeId}/compose/${nodeComposeProjectId}/variables/${
                variable.id
              }/uniquename?value=${value}`
            ),
          "Another variable with this name already exists"
        )
    ),
    isSecret: z.boolean(),
    value: z.string(),
    nodeComposeProjectId: z.number(),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      ...variable,
      nodeComposeProjectId: Number(nodeComposeProjectId),
    },
  })

  const handleCloseForm = () => {
    setOpenState(false)
    form.reset()
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/compose/${nodeComposeProjectId}/variables/${
        variable.id
      }`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    )
    if (!response.ok) {
      handleCloseForm()
      toastFailed("There was a problem when saving the variable. Try again!")
    } else {
      mutateNodeComposeVariables()
      setTimeout(() => {
        handleCloseForm()
        toastSuccess("Variable has been saved.")
      }, 500)
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
                <DialogTitle>Edit Variable</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4 group-disabled:opacity-50">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Name</FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="isSecret"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Secret</FormLabel>
                      <div>
                        <FormControl>
                          <Checkbox
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </div>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="value"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Value</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          type={
                            form.getValues()["isSecret"] ? "password" : "text"
                          }
                        />
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
