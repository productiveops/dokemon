import { Icons } from "@/components/icons"
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useTheme } from "@/components/ui/theme-provider"
import apiBaseUrl from "@/lib/api-base-url"
import { cn, toastFailed, trimString } from "@/lib/utils"
import { zodResolver } from "@hookform/resolvers/zod"
import axios from "axios"
import { useMemo, useState } from "react"
import { SubmitHandler, useForm } from "react-hook-form"
import { useNavigate } from "react-router-dom"
import { z } from "zod"

export default function ChangePassword() {
  const navigate = useNavigate()
  const [isSaving, setIsSaving] = useState(false)
  const { theme } = useTheme()

  const formSchema = z.object({
    currentPassword: z.preprocess(trimString, z.string()),
    newPassword: z.preprocess(
      trimString,
      z.string().min(8, "Should be at least 8 characters long")
    ),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: useMemo(() => {
      return { currentPassword: "", newPassword: "" }
    }, []),
  })

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)

    try {
      const response = await axios(`${apiBaseUrl()}/changepassword`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        data: JSON.stringify(data),
      })

      if (response?.status === 204) {
        navigate("/nodes")
      } else {
        toastFailed("Password update failed")
      }
    } catch (e) {
      if (axios.isAxiosError(e)) {
        toastFailed(e.response?.data)
      }
    }

    setIsSaving(false)
  }

  return (
    <>
      <div className="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-sm">
          <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
            <img
              className="mx-auto h-8 w-32"
              src={`/assets/images/${
                theme === "light"
                  ? "dokemon-light.svg"
                  : "dokemon-dark-small.svg"
              }`}
              alt="Dokémon"
            />
          </h2>
          <h3 className="mt-10 text-center text-xl font-bold leading-9 tracking-tight text-foreground">
            Change Password
          </h3>
        </div>

        <div className="mt-10  sm:mx-auto sm:w-full sm:max-w-sm sm:rounded-xl">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <fieldset className={cn("group")} disabled={isSaving}>
                <div className="max-w-2xl pb-4">
                  <FormField
                    control={form.control}
                    name="currentPassword"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Current Password</FormLabel>
                        <FormControl>
                          <Input autoFocus {...field} type="password" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="max-w-2xl pb-4">
                  <FormField
                    control={form.control}
                    name="newPassword"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Password</FormLabel>
                        <FormControl>
                          <Input {...field} type="password" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <Button
                  type="submit"
                  className={cn(
                    "mb-3 flex w-full justify-center rounded-md bg-amber-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-amber-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600"
                  )}
                  disabled={isSaving}
                >
                  <Icons.spinner
                    className={cn(
                      "absolute animate-spin text-slate-100 group-enabled:opacity-0"
                    )}
                  />
                  <span className={cn("group-disabled:opacity-0")}>
                    Change Password
                  </span>
                </Button>
                <Button
                  variant={"secondary"}
                  type="button"
                  className={cn(
                    "flex w-full justify-center rounded-md px-3 py-1.5 text-sm font-semibold leading-6 shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                  )}
                  disabled={isSaving}
                  onClick={() => navigate("/nodes")}
                >
                  <Icons.spinner
                    className={cn(
                      "absolute animate-spin text-slate-100 group-enabled:opacity-0"
                    )}
                  />
                  <span className={cn("group-disabled:opacity-0")}>Cancel</span>
                </Button>
              </fieldset>
            </form>
          </Form>
        </div>
      </div>
    </>
  )
}
