import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Field,
  FieldGroup,
  FieldLabel,
  FieldDescription,
  FieldError,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { DialogClose } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { createEnvironmentSchema } from "@/features/schemas";
import { useCreateEnvironment } from "../../hooks/useEnvironments";


type CreateEnvironmentForm = z.infer<typeof createEnvironmentSchema>;

interface CreateEnvironmentModalProps {
  projectId: string;
  onCreated?: () => void;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function CreateEnvironmentModal({
  projectId,
  onCreated,
  open,
  onOpenChange,
}: CreateEnvironmentModalProps) {
  const form = useForm<CreateEnvironmentForm>({
    resolver: zodResolver(createEnvironmentSchema),
    defaultValues: {
      name: "",
      connection_string: "",
      type: "development",
    },
  });

  const { mutate, isPending } = useCreateEnvironment();

  const onSubmit = (data: CreateEnvironmentForm) => {

    mutate(
      { 
        project_id: projectId, 
        name: data.name,
        connection_url: data.connection_string,
        type: data.type
      },
      {
        onSuccess: () => {
          form.reset();
          onCreated?.();
          onOpenChange?.(false);
        },
      },
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          New Environment
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Create New Environment</DialogTitle>
          <DialogDescription>
            Add a new database environment to your project.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="name">Environment Name</FieldLabel>
              <Input
                id="name"
                placeholder="e.g. Production Database"
                {...form.register("name")}
              />
              <FieldDescription>
                This will be used to identify your database environment.
              </FieldDescription>
              <FieldError
                errors={
                  form.formState.errors.name ? [form.formState.errors.name] : []
                }
              />
            </Field>

            <Field>
              <FieldLabel htmlFor="connection_string">
                Connection String
              </FieldLabel>
              <Input
                id="connection_string"
                placeholder="postgresql://envoy_superuser:superuser_password_123@envoy-postgres-dev:5432/envoy_dev"
                {...form.register("connection_string")}
              />
              <FieldDescription>
                The database connection string. Include SSL mode in the URL
                parameters.
              </FieldDescription>
              <FieldError
                errors={
                  form.formState.errors.connection_string
                    ? [form.formState.errors.connection_string]
                    : []
                }
              />
            </Field>

            <Field>
              <FieldLabel htmlFor="type">Environment Type</FieldLabel>
              <div className="flex gap-2 rounded-lg w-full">
                {["development", "staging", "production"].map((option) => (
                  <Button
                    key={option}
                    type="button"
                    size="sm"
                    onClick={() =>
                      form.setValue(
                        "type",
                        option as CreateEnvironmentForm["type"],
                      )
                    }
                    variant={
                      form.watch("type") === option ? "selected" : "unselected"
                    }
                  >
                    {option.charAt(0).toUpperCase() + option.slice(1)}
                  </Button>
                ))}
              </div>
            </Field>
          </FieldGroup>

          <DialogFooter className="gap-2 sm:gap-0">
            <DialogClose asChild>
              <Button type="button" variant="ghost">
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isPending}>
              {isPending ? "Creating..." : "Create Environment"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
