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
} from "@/lib/utils"
import apiBaseUrl from "@/lib/api-base-url"
import { useParams } from "react-router-dom"
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
import useComposeLibraryItemList from "@/hooks/useComposeLibraryItemList"
import useNodeComposeList from "@/hooks/useNodeComposeList"
import SpinnerIcon from "@/components/widgets/spinner-icon"

export default function ComposeAddProjectDialog() {
  const { nodeId } = useParams()
  const { composeLibraryItems } = useComposeLibraryItemList()
  const { mutateNodeComposeList } = useNodeComposeList(nodeId!)

  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [libraryProjectComboOpen, setLibraryProjectComboOpen] = useState(false)
  const [projectNameModified, setProjectNameModified] = useState(false)

  const formSchema = z.object({
    libraryProjectId: z.number().nullable(),
    libraryProjectName: z.string().min(1, "Please select library project"),
    projectName: z
      .string()
      .min(1, "Project Name is required")
      .max(20)
      .regex(REGEX_IDENTIFIER, REGEX_IDENTIFIER_MESSAGE)
      .refine(
        async (value) =>
          hasUniqueName(
            `${apiBaseUrl()}/nodes/${nodeId}/compose/uniquename?value=${value}`
          ),
        "Another project with this name already exists"
      ),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      libraryProjectId: null,
      libraryProjectName: "",
      projectName: "",
    },
  })

  const handleCloseForm = () => {
    setOpen(false)
    form.reset()
    setProjectNameModified(false)
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/compose/create/library`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    )
    if (!response.ok) {
      handleCloseForm()
      toastFailed("There was a problem when adding project. Try again!")
    } else {
      mutateNodeComposeList()
      handleCloseForm()
      toastSuccess("Project has been added.")
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>Add From Library</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <fieldset className={cn("group")} disabled={isSaving}>
              <DialogHeader>
                <DialogTitle>Add From Library</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4 group-disabled:opacity-50">
                <FormField
                  control={form.control}
                  name="libraryProjectName"
                  render={({ field }) => (
                    <FormItem className="flex flex-col">
                      <FormLabel>Library Project</FormLabel>
                      <FormControl>
                        <Popover
                          open={libraryProjectComboOpen}
                          onOpenChange={setLibraryProjectComboOpen}
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
                                  ? composeLibraryItems?.items.find(
                                      (item) => item.projectName === field.value
                                    )?.projectName
                                  : "Select project"}
                                <ChevronsUpDown className="ml-2 mt-1 h-4 w-4 shrink-0 opacity-50" />
                              </Button>
                            </FormControl>
                          </PopoverTrigger>
                          <PopoverContent className="w-[200px] p-0">
                            <Command>
                              <CommandInput placeholder="Select project..." />
                              <CommandEmpty>No project found.</CommandEmpty>
                              <CommandGroup>
                                {composeLibraryItems?.items.map((item) => (
                                  <CommandItem
                                    value={item.projectName}
                                    key={item.projectName}
                                    onSelect={() => {
                                      form.setValue(
                                        "libraryProjectId",
                                        item.id!
                                      )
                                      form.setValue(
                                        "libraryProjectName",
                                        item.projectName
                                      )
                                      form.trigger("libraryProjectName")

                                      setLibraryProjectComboOpen(false)
                                      if (!projectNameModified) {
                                        form.setValue(
                                          "projectName",
                                          item.projectName
                                        )
                                      }
                                    }}
                                  >
                                    <Check
                                      className={cn(
                                        "mr-2 h-4 w-4",
                                        item.projectName === field.value
                                          ? "opacity-100"
                                          : "opacity-0"
                                      )}
                                    />
                                    {item.projectName}
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
                <FormField
                  control={form.control}
                  name="projectName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Project Name</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          onKeyUp={() => setProjectNameModified(true)}
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
