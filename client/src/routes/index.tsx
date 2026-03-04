import { cn } from "@/lib/utils"
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldDescription
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { toast } from "sonner"
import { registerSchema, type RegisterFormData } from "@/features/auth/schema"
import { useRegisterMutation } from "@/features/auth/hooks"

export function RegisterForm({ className, ...props }: React.ComponentProps<"div">) {
  const { mutate: registerMutation, isPending: isRegistering} = useRegisterMutation()
  const navigate = useNavigate();


  const form = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: { email: "", password: "" },
  })

  function onSubmit(values: RegisterFormData) {
    registerMutation(values, {
      onSuccess: () => {
       toast.success("Login successfully.")
        navigate({
          to: "/app"
        })
      },
      onError: (e) => {
        toast.error(e.message)
      }
    });
  }

  return (
    <main className="h-screen flex flex-col items-center justify-center">
    <div className={cn("flex flex-col gap-6 w-full max-w-md mx-auto", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Create an account</CardTitle>
          <CardDescription>
           Get started with Envoy
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  {...form.register('email')}
                  disabled={isRegistering}
                />
                {form.formState.errors.email && (
                  <FieldError>{form.formState.errors.email.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="password">Password</FieldLabel>
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  {...form.register('password')}
                  disabled={isRegistering}
                />
                {form.formState.errors.password && (
                  <FieldError>{form.formState.errors.password.message}</FieldError>
                )}
              </Field>

              <Field>
                <Button isLoading={isRegistering} type="submit" className="w-full" disabled={isRegistering}>
                  Register
                </Button>
                <FieldDescription className="text-center text-sm">
                  Already have an account? <Link to="/login" className="underline underline-offset-4">Login</Link>
                </FieldDescription>
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
    </main>
  )
}

export const Route = createFileRoute('/')({
  component: () => <RegisterForm />
})