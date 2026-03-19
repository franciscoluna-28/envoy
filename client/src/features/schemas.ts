import { z } from "zod";

export const createEnvironmentSchema = z.object({
  name: z
    .string()
    .min(1, "Environment name is required")
    .max(100, "Environment name must be less than 100 characters"),
  connection_string: z.string().min(1, "Connection string is required"),
  type: z.enum(["production", "staging", "development"]),
});

export const createProjectSchema = z.object({
  name: z.string().min(1, 'Project name is required').max(100, 'Project name must be less than 100 characters')
})

export const updateProjectSchema = z.object({
  name: z.string().min(1, 'Project name is required').max(100, 'Project name must be less than 100 characters')
})