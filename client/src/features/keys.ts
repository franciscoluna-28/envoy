export const ENVIRONMENTS_QUERY_KEYS = {
  all: ['environments'] as const,
  lists: () => [...ENVIRONMENTS_QUERY_KEYS.all, 'list'] as const,
  list: (projectId: string) => [...ENVIRONMENTS_QUERY_KEYS.lists(), { projectId }] as const,
}

export const MIGRATIONS_QUERY_KEYS = {
  all: ['migrations'] as const,
  lists: () => [...MIGRATIONS_QUERY_KEYS.all, 'list'] as const,
  list: (envId: string) => [...MIGRATIONS_QUERY_KEYS.lists(), { envId }] as const,
  preview: (envId: string, sql: string) => [...MIGRATIONS_QUERY_KEYS.all, 'preview', envId, sql] as const,
}

export const PROJECTS_QUERY_KEYS = {
  all: ['projects'] as const,
  lists: () => [...PROJECTS_QUERY_KEYS.all, 'list'] as const,
  list: (filter: string) => [...PROJECTS_QUERY_KEYS.lists(), { filter }] as const,
  details: () => [...PROJECTS_QUERY_KEYS.all, 'detail'] as const,
  detail: (id: string) => [...PROJECTS_QUERY_KEYS.details(), id] as const,
}

export const ENVIRONMENT_SCHEMA_QUERY_KEYS = {
  all: ['environment-schema'] as const,
  details: () => [...ENVIRONMENT_SCHEMA_QUERY_KEYS.all, 'detail'] as const,
  detail: (envId: string) => [...ENVIRONMENT_SCHEMA_QUERY_KEYS.details(), { envId }] as const,
}