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
  cn,
  hasUniqueName,
  toastSomethingWentWrong,
  toastSuccess,
  trimString,
} from "@/lib/utils"
import useNodes from "@/hooks/useNodes"
import apiBaseUrl from "@/lib/api-base-url"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { Check, ChevronsUpDown } from "lucide-react"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/components/ui/command"
import useEnvironments from "@/hooks/useEnvironments"
import { INodeCreateRequest } from "@/lib/api-models"
import { apiNodesCreate } from "@/lib/api"
import SpinnerIcon from "@/components/widgets/spinner-icon"

export default function NodeAddDialog({ disabled }: { disabled: boolean }) {
  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const { mutateNodes } = useNodes()
  const { environments } = useEnvironments()
  const [environmentsComboOpen, setEnvironmentsComboOpen] = useState(false)

  const formSchema = z.object({
    name: z.preprocess(
      trimString,
      z
        .string()
        .min(1, "Name is required")
        .max(20)
        .refine(
          async (value) =>
            hasUniqueName(`${apiBaseUrl()}/nodes/uniquename?value=${value}`),
          "Another node with this name already exists"
        )
    ),
    environmentId: z.number().nullable(),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      environmentId: null,
    },
  })

  const handleCloseForm = () => {
    setOpen(false)
    form.reset()
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (
    data: INodeCreateRequest
  ) => {
    setIsSaving(true)
    const response = await apiNodesCreate(data)
    if (!response.ok) {
      handleCloseForm()
      toastSomethingWentWrong(
        "There was a problem when creating new node. Try again!"
      )
    } else {
      mutateNodes()
      handleCloseForm()
      toastSuccess("New node has been added.")
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          className={cn(disabled ? "opacity-70 hover:bg-amber-600" : "")}
          title={disabled ? "Please set Server URL before adding nodes" : ""}
          disabled={disabled}
        >
          Add Node
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <fieldset className={cn("group")} disabled={isSaving}>
              <DialogHeader>
                <DialogTitle>Add Node</DialogTitle>
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
                  name="environmentId"
                  render={({ field }) => (
                    <FormItem className="flex w-[300px] flex-col">
                      <FormLabel>Environment</FormLabel>
                      <FormControl>
                        <Popover
                          open={environmentsComboOpen}
                          onOpenChange={setEnvironmentsComboOpen}
                        >
                          <PopoverTrigger asChild>
                            <FormControl>
                              <Button
                                variant="outline"
                                role="combobox"
                                className={cn(
                                  "inline-flex w-auto justify-between font-normal text-slate-800 dark:text-slate-50",
                                  !field.value && "text-muted-foreground"
                                )}
                              >
                                {field.value
                                  ? environments?.items.find(
                                      (item) => item.id === field.value
                                    )?.name
                                  : "(None)"}
                                <ChevronsUpDown className="ml-2 mt-1 h-4 w-4 shrink-0 opacity-50" />
                              </Button>
                            </FormControl>
                          </PopoverTrigger>
                          <PopoverContent className="w-[300px] p-0">
                            <Command>
                              <CommandInput placeholder="Select environment..." />
                              <CommandEmpty>No environment found.</CommandEmpty>
                              <CommandGroup>
                                <CommandItem
                                  onSelect={() => {
                                    form.setValue("environmentId", null)
                                    form.trigger("environmentId")
                                    setEnvironmentsComboOpen(false)
                                  }}
                                >
                                  <Check
                                    className={cn(
                                      "mr-2 h-4 w-4",
                                      field.value === null
                                        ? "opacity-100"
                                        : "opacity-0"
                                    )}
                                  />
                                  None
                                </CommandItem>
                                {environments?.items.map((item) => (
                                  <CommandItem
                                    value={item.name}
                                    key={item.id}
                                    onSelect={() => {
                                      form.setValue("environmentId", item.id)
                                      form.trigger("environmentId")
                                      setEnvironmentsComboOpen(false)
                                    }}
                                  >
                                    <Check
                                      className={cn(
                                        "mr-2 h-4 w-4",
                                        item.id === field.value
                                          ? "opacity-100"
                                          : "opacity-0"
                                      )}
                                    />
                                    {item.name}
                                  </CommandItem>
                                ))}
                              </CommandGroup>
                            </Command>
                          </PopoverContent>
                        </Popover>
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
