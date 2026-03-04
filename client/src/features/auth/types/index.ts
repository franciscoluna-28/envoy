import type { paths } from "@/api/api";

export type RegisterInputDto = paths["/api/v1/auth/register"]["post"]["requestBody"]["content"]["application/json"]
export type LoginInputDto = paths["/api/v1/auth/login"]["post"]["requestBody"]["content"]["application/json"]
