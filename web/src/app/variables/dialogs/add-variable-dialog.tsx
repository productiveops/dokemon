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
import { Icons } from "@/components/icons"
import { Input } from "@/components/ui/input"
import { SubmitHandler, useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { cn, trimString } from "@/lib/utils"
import useVariables from "@/hooks/useVariables"
import { Checkbox } from "@/components/ui/checkbox"
import apiBaseUrl from "@/lib/api-base-url"
import { toast } from "@/components/ui/use-toast"

export default function AddVariableDialog() {
  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const { mutateVariables } = useVariables()

  const formSchema = z.object({
    name: z.preprocess(
      trimString,
      z
        .string()
        .min(1, "Name is required")
        .max(100)
        .refine(async (value) => {
          const res = await fetch(
            `${apiBaseUrl()}/variables/uniquename?value=${value}`
          )
          return (await res.json()).unique
        }, "Another variable with this name already exists")
    ),
    isSecret: z.boolean(),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      isSecret: true,
    },
  })

  const handleCloseForm = () => {
    setOpen(false)
    form.reset()
  }

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)
    const response = await fetch(`${apiBaseUrl()}/variables`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })
    if (!response.ok) {
      handleCloseForm()
      toast({
        variant: "destructive",
        title: "Uh oh! Something went wrong.",
        description: "There was a problem when saving new variable. Try again!",
      })
    } else {
      mutateVariables()
      setTimeout(() => {
        handleCloseForm()
        toast({
          title: "Success!",
          description: "New variable has been added.",
        })
      }, 500)
    }
    setIsSaving(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>Add Variable</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <fieldset className={cn("group")} disabled={isSaving}>
              <DialogHeader>
                <DialogTitle>Add Variable</DialogTitle>
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
