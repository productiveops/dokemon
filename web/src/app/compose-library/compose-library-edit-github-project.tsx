import { useNavigate, useParams } from "react-router-dom"
import {
  Breadcrumb,
  BreadcrumbCurrent,
  BreadcrumbLink,
  BreadcrumbSeparator,
} from "@/components/widgets/breadcrumb"
import apiBaseUrl from "@/lib/api-base-url"
import { useEffect, useMemo, useRef, useState } from "react"
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
  REGEX_IDENTIFIER,
  REGEX_IDENTIFIER_MESSAGE,
  cn,
  hasUniqueName,
  initMonaco,
  toastFailed,
  toastSuccess,
} from "@/lib/utils"
import { SubmitHandler, useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import {
  MainContainer,
  Section,
  SectionBody,
} from "@/components/widgets/main-container"
import Editor, { OnMount } from "@monaco-editor/react"
import type monaco from "monaco-editor"
import { Input } from "@/components/ui/input"
import { useTheme } from "@/components/ui/theme-provider"
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
import useCredentials from "@/hooks/useCredentials"
import GitHubPATAddDialog from "../credentials/github-pat-add-dialog"
import useGitHubComposeLibraryItem from "@/hooks/useGitHubComposeLibraryItem"
import DeleteDialog from "@/components/delete-dialog"
import useComposeLibraryItemList from "@/hooks/useComposeLibraryItemList"

export default function ComposeLibraryEditGitHubProject() {
  const { composeProjectId } = useParams()
  const { gitHubComposeLibraryItem } = useGitHubComposeLibraryItem(
    composeProjectId!
  )
  const { mutateComposeLibraryItemList } = useComposeLibraryItemList()
  const [isSaving, setIsSaving] = useState(false)
  const { credentials } = useCredentials()
  const [credentialsComboOpen, setCredentialsComboOpen] = useState(false)
  const { theme } = useTheme()
  const navigate = useNavigate()

  initMonaco()

  const formSchema = z.object({
    projectName: z
      .string()
      .min(1, "Name is required")
      .max(20)
      .regex(REGEX_IDENTIFIER, REGEX_IDENTIFIER_MESSAGE)
      .refine(
        async (value) =>
          hasUniqueName(
            `${apiBaseUrl()}/composelibrary/uniquenameexcludeitself?newvalue=${value}&currentvalue=${gitHubComposeLibraryItem?.projectName}&id=${gitHubComposeLibraryItem?.id}`
          ),
        "Another project with this name already exists"
      ),
    url: z.string().url(),
    credentialId: z.number().nullable(),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: useMemo(() => {
      return { projectName: "", url: "", credentialId: 0 }
    }, []),
  })

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)
    const response = await fetch(
      `${apiBaseUrl()}/composelibrary/github/${composeProjectId}`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      toastFailed(r.errors?.body)
    } else {
      toastSuccess("Project has been saved.")
    }
    setIsSaving(false)
  }

  useEffect(() => {
    form.reset({
      credentialId: gitHubComposeLibraryItem?.credentialId,
      projectName: gitHubComposeLibraryItem?.projectName,
      url: gitHubComposeLibraryItem?.url,
    })
    handleLoadFileContent()
  }, [gitHubComposeLibraryItem])

  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor>()

  const handleEditorDidMount: OnMount = (editor, _monaco) => {
    editorRef.current = editor
  }

  const handleLoadFileContent = async (
    e?: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    if (e) e.preventDefault()

    const url = form.getValues("url")
    const credentialId = form.getValues("credentialId")
    if (!url) {
      return
    }

    const response = await fetch(`${apiBaseUrl()}/github/filecontent/load`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        credentialId: credentialId,
        url: url,
      }),
    })
    if (!response.ok) {
      const data: { errors: { body: string } } = await response.json()
      toastFailed(data.errors?.body)
      editorRef.current?.setValue("# Error while loading content")
    } else {
      const data: { content: string } = await response.json()
      editorRef.current?.setValue(data.content)
    }
  }

  const [deleteInProgress, setDeleteInProgress] = useState(false)
  const handleDelete = async () => {
    setDeleteInProgress(true)
    const response = await fetch(
      `${apiBaseUrl()}/composelibrary/github/${composeProjectId}`,
      {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      }
    )
    if (!response.ok) {
      const r = await response.json()
      toastFailed(r.errors?.body)
    } else {
      mutateComposeLibraryItemList()
      setTimeout(() => {
        toastSuccess("Compose project deleted.")
        navigate("/composelibrary")
      }, 500)
    }
    setDeleteInProgress(false)
  }

  return (
    <MainArea>
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/composelibrary">Compose Library</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>
            {gitHubComposeLibraryItem?.projectName}
          </BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Edit</BreadcrumbCurrent>
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
        <DeleteDialog
          deleteCaption="Delete"
          title="Delete Compose Project"
          message={`Are you sure you want to delete project '${gitHubComposeLibraryItem?.projectName}'?`}
          deleteHandler={handleDelete}
          isProcessing={deleteInProgress}
        />
      </div>
      <MainContent>
        <MainContainer>
          <Section>
            <SectionBody>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                  <fieldset className={cn("group")} disabled={isSaving}>
                    <div className="mb-4 flex max-w-2xl flex-col gap-4 pb-4">
                      <FormField
                        control={form.control}
                        name="projectName"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Library Project Name</FormLabel>
                            <FormControl>
                              <Input {...field} autoFocus />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="url"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>GitHub URL of Compose File</FormLabel>
                            <FormControl>
                              <Input {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="credentialId"
                        render={({ field }) => (
                          <FormItem className="flex flex-col">
                            <FormLabel>
                              GitHub Credentials (For private repository)
                            </FormLabel>
                            <div className="flex gap-2">
                              <FormControl>
                                <Popover
                                  open={credentialsComboOpen}
                                  onOpenChange={setCredentialsComboOpen}
                                >
                                  <PopoverTrigger asChild>
                                    <FormControl>
                                      <Button
                                        variant="outline"
                                        role="combobox"
                                        className={cn(
                                          "inline-flex w-[300px] justify-between font-normal text-slate-800 dark:text-slate-50",
                                          !field.value &&
                                            "text-muted-foreground"
                                        )}
                                      >
                                        {field.value
                                          ? credentials?.items.find(
                                              (item) => item.id === field.value
                                            )?.name
                                          : "(None)"}
                                        <ChevronsUpDown className="ml-2 mt-1 h-4 w-4 shrink-0 opacity-50" />
                                      </Button>
                                    </FormControl>
                                  </PopoverTrigger>
                                  <PopoverContent className="w-[300px] p-0">
                                    <Command>
                                      <CommandInput placeholder="Select credential..." />
                                      <CommandEmpty>
                                        No credentials found.
                                      </CommandEmpty>
                                      <CommandGroup>
                                        <CommandItem
                                          onSelect={() => {
                                            form.setValue("credentialId", null)
                                            form.trigger("credentialId")
                                            setCredentialsComboOpen(false)
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
                                        {credentials?.items?.map((item) => (
                                          <CommandItem
                                            value={item.name}
                                            key={item.id}
                                            onSelect={() => {
                                              form.setValue(
                                                "credentialId",
                                                item.id
                                              )
                                              form.trigger("credentialId")
                                              setCredentialsComboOpen(false)
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
                              <GitHubPATAddDialog buttonCaption="Add New" />
                            </div>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>
                    <div className="mb-6">
                      <Button
                        variant={"default"}
                        onClick={handleLoadFileContent}
                      >
                        Refresh File Content
                      </Button>
                    </div>
                    <div>
                      <FormLabel className="block pb-4">
                        Definition (Read only)
                      </FormLabel>
                      <Editor
                        theme={theme}
                        height="50vh"
                        defaultLanguage="yaml"
                        defaultValue={""}
                        options={{
                          minimap: { enabled: false },
                          readOnly: true,
                        }}
                        onMount={handleEditorDidMount}
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
