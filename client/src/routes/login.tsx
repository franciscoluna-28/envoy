import { cn } from "@/lib/utils";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldDescription,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { loginSchema, type LoginFormData } from "@/features/auth/schema";
import { useLoginMutation } from "@/features/auth/hooks";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const { mutate: loginMutation, isPending: isLoggingIn } = useLoginMutation();
  const navigate = useNavigate();

  const form = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  function onSubmit(values: LoginFormData) {
    loginMutation(values, {
      onSuccess: () => {
        toast.success("Welcome back.");
        navigate({
          to: "/app",
        });
      },
      onError: (e) => {
        toast.error(e.message);
      },
    });
  }

  return (
    <main className="h-screen flex flex-col items-center justify-center">
      <div
        className={cn("flex flex-col gap-6 w-full max-w-md mx-auto", className)}
        {...props}
      >
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Welcome back</CardTitle>
            <CardDescription>
              Enter your credentials to access Envoy
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
                    {...form.register("email")}
                    disabled={isLoggingIn}
                  />
                  {form.formState.errors.email && (
                    <FieldError>
                      {form.formState.errors.email.message}
                    </FieldError>
                  )}
                </Field>

                <Field>
                  <div className="flex items-center justify-between">
                    <FieldLabel htmlFor="password">Password</FieldLabel>
                    {/*  <Link to="/forgot-password" size="sm" className="text-xs underline-offset-4 hover:underline">
                      Forgot?
                    </Link> */}
                  </div>
                  <Input
                    id="password"
                    type="password"
                    placeholder="••••••••"
                    {...form.register("password")}
                    disabled={isLoggingIn}
                  />
                  {form.formState.errors.password && (
                    <FieldError>
                      {form.formState.errors.password.message}
                    </FieldError>
                  )}
                </Field>

                <Field>
                  <Button
                    isLoading={isLoggingIn}
                    type="submit"
                    className="w-full"
                  >
                    Login
                  </Button>
                  <FieldDescription className="text-center text-sm mt-2">
                    Don&apos;t have an account?{" "}
                    <Link
                      to="/"
                      className="underline underline-offset-4 font-medium"
                    >
                      Sign up
                    </Link>
                  </FieldDescription>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>
      </div>
    </main>
  );
}

export const Route = createFileRoute("/login")({
  component: () => <LoginForm />,
});
