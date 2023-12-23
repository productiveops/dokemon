import { useParams } from "react-router-dom"
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
import { REGEX_IDENTIFIER, REGEX_IDENTIFIER_MESSAGE, cn } from "@/lib/utils"
import { SubmitHandler, useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import {
  MainContainer,
  Section,
  SectionBody,
} from "@/components/widgets/main-container"
import Editor, { loader, OnMount } from "@monaco-editor/react"
import type monaco from "monaco-editor"
import { Input } from "@/components/ui/input"
import { toast } from "@/components/ui/use-toast"
import { useTheme } from "@/components/ui/theme-provider"
import useNodeHead from "@/hooks/useNodeHead"
import useNodeComposeItem from "@/hooks/useNodeComposeItem"
import ComposeVariableEditor from "./variable-editor/compose-variable-editor"

export default function ComposeDefinitionLocal() {
  const { nodeId } = useParams()
  const { nodeHead } = useNodeHead(nodeId!)
  const { composeProjectId } = useParams()
  const { nodeComposeItem, mutateNodeComposeItem } = useNodeComposeItem(
    nodeId!,
    composeProjectId!
  )
  const [isSaving, setIsSaving] = useState(false)
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor>()
  const [editorMounted, setEditorMounted] = useState(1)
  const { theme } = useTheme()

  const isLibraryProject = () =>
    !!(nodeComposeItem?.libraryProjectId || nodeComposeItem?.libraryProjectName)

  const formSchema = z.object({
    projectName: z
      .string()
      .min(1, "Name is required")
      .max(20)
      .regex(REGEX_IDENTIFIER, REGEX_IDENTIFIER_MESSAGE)
      .refine(async (value) => {
        const res = await fetch(
          `${apiBaseUrl()}/nodes/${nodeId}/compose/${composeProjectId}/uniquename?value=${value}`
        )
        return (await res.json()).unique
      }, "Another project with this name already exists"),
    definition: z.string().optional(),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: useMemo(() => {
      if (nodeComposeItem) {
        editorRef.current?.setValue(nodeComposeItem?.definition!)
        return {
          projectName: nodeComposeItem.projectName,
          definition: nodeComposeItem.definition,
        }
      }

      return { projectName: "", definition: "" }
    }, []),
  })

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    data.definition = editorRef.current?.getValue()
    setIsSaving(true)
    const response = await fetch(
      `${apiBaseUrl()}/nodes/${nodeId}/compose/${composeProjectId}/local`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    )
    if (!response.ok) {
      const r = await response.json()
      toast({
        variant: "destructive",
        title: "Error saving project",
        description: r.errors?.body,
      })
    } else {
      mutateNodeComposeItem()
      toast({
        title: "Success!",
        description: "Project has been saved.",
      })
    }
    setIsSaving(false)
  }

  useEffect(() => {
    form.reset({
      projectName: nodeComposeItem?.projectName,
      definition: nodeComposeItem?.definition,
    })
    if (nodeComposeItem?.definition && editorRef.current) {
      editorRef.current.setValue(nodeComposeItem.definition)
    }
  }, [nodeComposeItem, editorMounted])

  const handleEditorDidMount: OnMount = (editor, _monaco) => {
    editorRef.current = editor
    setEditorMounted(editorMounted + 1)
  }

  loader.init().then((monaco) => {
    monaco.editor.defineTheme("dark", {
      base: "vs-dark",
      inherit: true,
      rules: [],
      colors: {
        "editor.background": "#000000",
      },
    })
  })

  return (
    <MainArea>
      <TopBar>
        <Breadcrumb>
          <BreadcrumbLink to="/nodes">Nodes</BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeHead?.name}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbLink to={`/nodes/${nodeId}/compose`}>
            Compose
          </BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>{nodeComposeItem?.projectName}</BreadcrumbCurrent>
          <BreadcrumbSeparator />
          <BreadcrumbCurrent>Definition</BreadcrumbCurrent>
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
                    <div className="mb-4 flex max-w-2xl flex-col gap-4 pb-4">
                      <FormField
                        control={form.control}
                        name="projectName"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>
                              Project Name{" "}
                              {nodeComposeItem?.status?.startsWith("running")
                                ? "(Running)"
                                : ""}
                            </FormLabel>
                            <FormControl>
                              <Input
                                {...field}
                                autoFocus
                                disabled={nodeComposeItem?.status?.startsWith(
                                  "running"
                                )}
                                title={
                                  nodeComposeItem?.status?.startsWith("running")
                                    ? "Names of running projects cannot be edited"
                                    : ""
                                }
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      {isLibraryProject() && (
                        <FormItem>
                          <FormLabel>Library Project</FormLabel>
                          <FormControl>
                            <Input
                              value={nodeComposeItem?.libraryProjectName}
                              disabled={true}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    </div>
                    <div>
                      <FormLabel className="block pb-4">
                        Definition {isLibraryProject() ? "(Read only)" : ""}
                      </FormLabel>
                      <Editor
                        theme={theme}
                        height="50vh"
                        defaultLanguage="yaml"
                        defaultValue={""}
                        options={{
                          minimap: { enabled: false },
                          readOnly: isLibraryProject(),
                        }}
                        onMount={handleEditorDidMount}
                      />
                    </div>
                  </fieldset>
                </form>
              </Form>
            </SectionBody>
          </Section>
          <Section>
            <SectionBody>
              <ComposeVariableEditor
                nodeId={nodeId!}
                nodeComposeProjectId={composeProjectId!}
              />
            </SectionBody>
          </Section>
        </MainContainer>
      </MainContent>
    </MainArea>
  )
}
