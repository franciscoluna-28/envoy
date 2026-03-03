import { createFileRoute } from '@tanstack/react-router'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuth } from '@/hooks/useAuth'
import { Button } from "@/components/ui/button"
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useMutation } from '@tanstack/react-query'
import client from '../api/client'

const registerSchema = z.object({
  email: z.string().email('Please enter a valid email address').min(1, 'Email is required'),
  password: z.string().min(8, 'Password must be at least 8 characters long').min(1, 'Password is required'),
})

type RegisterFormData = z.infer<typeof registerSchema>



function RegisterForm() {
  const { mutate: registerMutation, isPending: isRegistering, error: mutationError } = useMutation({
    mutationFn: async (data: RegisterFormData) => {
      const response = await client.POST('/api/v1/auth/register', {
        body: data,
      });
      
      if (response.error) {
        const errorDetail = (response.error as any).detail || 'Registration failed';
        throw new Error(errorDetail);
      }
      
      return response.data;
    },
  });

  const form = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: { email: "", password: "" },
  })

  function onSubmit(values: RegisterFormData) {
    registerMutation(values, {
      onSuccess: () => {
        console.log('Account created!');
      },
    });
  }

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Create Account</CardTitle>
      </CardHeader>
      <CardContent>
        <FieldSet>
          <FieldLegend>Create your account</FieldLegend>
          

        

          <FieldGroup>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <Field>
                <FieldLabel>Email</FieldLabel>
                <Input
                  type="email"
                  placeholder="Enter your email"
                  {...form.register('email')}
              
                  disabled={isRegistering} 
                />
                {form.formState.errors.email && (
                  <FieldError>{form.formState.errors.email.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel>Password</FieldLabel>
                <Input
                  type="password"
                  placeholder="Create a password"
                  {...form.register('password')}
                  disabled={isRegistering}
                />
                {form.formState.errors.password && (
                  <FieldError>{form.formState.errors.password.message}</FieldError>
                )}
              </Field>

              <Button type="submit" className="w-full" disabled={isRegistering}>
                {isRegistering ? (
                  <span className="flex items-center gap-2">
                    Creating account...
                  </span>
                ) : (
                  "Register"
                )}
              </Button>
            </form>
          </FieldGroup>
        </FieldSet>
          {mutationError && (
            <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded text-sm">
              {mutationError.message}
            </div>
          )}
      </CardContent>
    </Card>
  )
}

export const Route = createFileRoute('/')({
  component: () => <RegisterForm />
})

