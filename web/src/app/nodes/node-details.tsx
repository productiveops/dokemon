import { useParams } from "react-router-dom"
import {
  Breadcrumb,
  BreadcrumbCurrent,
  BreadcrumbLink,
  BreadcrumbSeparator,
} from "@/components/widgets/breadcrumb"
import apiBaseUrl from "@/lib/api-base-url"
import { useEffect, useMemo, useState } from "react"
import TopBar from "@/components/widgets/top-bar"
import TopBarActions from "@/components/widgets/top-bar-actions"
import MainArea from "@/components/widgets/main-area"
import MainContent from "@/components/widgets/main-content"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { z } from "zod"
import {
  cn,
  hasUniqueName,
  toastSomethingWentWrong,
  toastSuccess,
  trimString,
} from "@/lib/utils"
import { SubmitHandler, useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import {
  MainContainer,
  Section,
  SectionBody,
} from "@/components/widgets/main-container"
import { Input } from "@/components/ui/input"
import useNodes from "@/hooks/useNodes"
import useEnvironments from "@/hooks/useEnvironments"
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
import useNode from "@/hooks/useNode"
import { INodeUpdateRequest } from "@/lib/api-models"
import { apiNodesUpdate } from "@/lib/api"

export default function NodeDetails() {
  const { nodeId } = useParams()
  const { node, mutateNode } = useNode(nodeId!)
  const { environments } = useEnvironments()
  const [environmentsComboOpen, setEnvironmentsComboOpen] = useState(false)
  const { mutateNodes } = useNodes()
  const [isSaving, setIsSaving] = useState(false)

  const formSchema = z.object({
    name: z.preprocess(
      trimString,
      z
        .string()
        .min(1, "Name is required")
        .max(20)
        .refine(
          async (value) =>
            hasUniqueName(
              `${apiBaseUrl()}/nodes/${nodeId}/uniquename?value=${value}`
            ),
          "Another node with this name already exists"
        )
    ),
    environmentId: z.number().nullable(),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: useMemo(() => {
      if (!node) return { name: "", environmentId: null }
      return node
    }, [node]),
  })

  const onSubmit: SubmitHandler<FormSchemaType> = async (
    data: INodeUpdateRequest
  ) => {
    setIsSaving(true)
    const response = await apiNodesUpdate(Number(nodeId), data)
    if (!response.ok) {
      toastSomethingWentWrong(
        "There was a problem when saving the node details. Try again!"
      )
    } else {
      mutateNode()
      mutateNodes()
      toastSuccess("Node details have been saved.")
    }
    setIsSaving(false)
  }

  useEffect(() => {
    if (node) {
      form.reset(node)
    }
  }, [node])

  return (
    <MainArea>
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{node?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Details</BreadcrumbCurrent>
        </Breadcrumb>
        <TopBarActions></TopBarActions>
      </TopBar>
      <div className="-mb-8 pt-4">
        <Button
          className="mb-4 mr-2 w-24"
          onClick={form.handleSubmit(onSubmit)}
        >
          Save
        </Button>
      </div>
      <MainContent>
        <MainContainer>
          <Section>
            <SectionBody>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                  <fieldset className={cn("group")} disabled={isSaving}>
                    <div className="flex max-w-2xl flex-col gap-6 pb-4">
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
                                    <CommandEmpty>
                                      No environment found.
                                    </CommandEmpty>
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
                                            form.setValue(
                                              "environmentId",
                                              item.id
                                            )
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
                  </fieldset>
                </form>
              </Form>
            </SectionBody>
          </Section>
        </MainContainer>
      </MainContent>
    </MainArea>
  )
}
