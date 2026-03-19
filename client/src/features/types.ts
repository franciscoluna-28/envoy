import type { paths, components } from "@/api/api";

export type Project = paths["/projects"]["get"]["responses"]["200"]["content"]["application/json"][number];
export type Environment = paths["/projects/{id}/environments"]["get"]["responses"]["200"]["content"]["application/json"][number];
export type DatabaseSchemaItem = components["schemas"]["environments.SchemaColumn"][];
export type EnvironmentId = NonNullable<Environment["id"]>;
export type ProjectId = NonNullable<Project["id"]>;
export type CreateEnvironmentInput = paths["/projects/{id}/environments"]["post"]["requestBody"]["content"]["application/json"];
export type SchemaColumn = components["schemas"]["environments.SchemaColumn"];

export interface PreviewError extends Error {
  errors?: string[];
}