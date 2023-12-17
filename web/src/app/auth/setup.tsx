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
import { toast } from "@/components/ui/use-toast"
import Loading from "@/components/widgets/loading"
import apiBaseUrl from "@/lib/api-base-url"
import { cn, trimString } from "@/lib/utils"
import { zodResolver } from "@hookform/resolvers/zod"
import axios from "axios"
import { useEffect, useMemo, useState } from "react"
import { SubmitHandler, useForm } from "react-hook-form"
import { useNavigate } from "react-router-dom"
import { z } from "zod"

export default function Setup() {
  const navigate = useNavigate()
  const [isSaving, setIsSaving] = useState(false)
  const [loaded, setLoaded] = useState(false)
  const { theme } = useTheme()

  useEffect(() => {
    async function setupComplete() {
      try {
        const response = await axios(`${apiBaseUrl()}/users/count`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
        })

        if (response.data.count > 0) {
          navigate("/login")
        } else {
          setLoaded(true)
        }
      } catch (e) {
        if (axios.isAxiosError(e)) {
          toast({
            variant: "destructive",
            title: "Failed",
            description: e.response?.data.errors?.body,
          })
        }
      }
    }
    setupComplete()
  }, [])

  const formSchema = z.object({
    userName: z.preprocess(
      trimString,
      z
        .string()
        .min(1, "Username is required")
        .max(20, "Username should not be over 20 characters in length")
    ),
    password: z.preprocess(
      trimString,
      z.string().min(8, "Should be at least 8 characters long")
    ),
  })

  type FormSchemaType = z.infer<typeof formSchema>

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    defaultValues: useMemo(() => {
      return { userName: "", password: "" }
    }, []),
  })

  const onSubmit: SubmitHandler<FormSchemaType> = async (data) => {
    setIsSaving(true)

    try {
      await axios(`${apiBaseUrl()}/users`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        data: JSON.stringify(data),
      })
      localStorage.setItem("userName", data.userName)
      navigate("/nodes")
      window.location.reload()
    } catch (e) {
      if (axios.isAxiosError(e)) {
        toast({
          variant: "destructive",
          title: "Failed",
          description: e.response?.data,
        })
      }
    }

    setIsSaving(false)
  }

  if (!loaded) return <Loading />

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
            Create New User
          </h3>
        </div>

        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <fieldset className={cn("group")} disabled={isSaving}>
                <div className="max-w-2xl pb-4">
                  <FormField
                    control={form.control}
                    name="userName"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Username</FormLabel>
                        <FormControl>
                          <Input autoFocus {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="max-w-2xl pb-4">
                  <FormField
                    control={form.control}
                    name="password"
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
                    "flex w-full justify-center rounded-md bg-amber-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-amber-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600"
                  )}
                  disabled={isSaving}
                >
                  <Icons.spinner
                    className={cn(
                      "absolute animate-spin text-slate-100 group-enabled:opacity-0"
                    )}
                  />
                  <span className={cn("group-disabled:opacity-0")}>
                    Create User
                  </span>
                </Button>
              </fieldset>
            </form>
          </Form>
        </div>
      </div>
    </>
  )
}
