import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import client from '@/api/client'
import { toast } from 'sonner'
import { useNavigate } from '@tanstack/react-router'

const MIGRATIONS_QUERY_KEYS = {
  all: ['migrations'] as const,
  lists: () => [...MIGRATIONS_QUERY_KEYS.all, 'list'] as const,
  list: (envId: string) => [...MIGRATIONS_QUERY_KEYS.lists(), { envId }] as const,
  preview: (envId: string, sql: string) => [...MIGRATIONS_QUERY_KEYS.all, 'preview', envId, sql] as const,
}

export function useGetEnvironmentMigrations(envId: string) {
  return useQuery({
    queryKey: MIGRATIONS_QUERY_KEYS.list(envId),
    queryFn: async () => {
      const response = await client.GET('/environments/{id}/migrations', {
        params: { path: { id: envId } }
      })

      // Check if the response indicates an error
      if (response.error) {
        throw new Error(response.error.message || 'Failed to fetch migrations')
      }

      return response.data
    },
    enabled: !!envId,
    retry: (failureCount, error) => {
      // Don't retry on 4xx errors, but retry on 5xx and network errors
      if (error?.message?.includes('401') || error?.message?.includes('404') || error?.message?.includes('403')) {
        return false
      }
      return failureCount < 3
    }
  })
}

export function usePreviewSchemaChanges() {
  return useMutation({
    mutationFn: async ({ envId, sqlContent }: { envId: string; sqlContent: string }) => {
      const response = await client.POST('/environments/{id}/migrations/preview', {
        params: { path: { id: envId } },
        body: {
          sql_content: sqlContent
        }
      })
      
      // Type assertion for the new response format
      const responseData = response.data as any
      
      // Handle the new response format
      if (!responseData?.success) {
        const error = new Error(responseData?.message || 'Failed to preview schema changes')
        // Attach detailed errors to the error object
        if (responseData?.errors) {
          (error as any).details = responseData.errors
        }
        throw error
      }
      
      return responseData.data
    },
    onError: (error: any) => {
      // Extract detailed errors if available
      const errorMessage = error.details?.length 
        ? error.details.join('; ') 
        : error.message || 'Failed to preview schema changes'
      
      toast.error(errorMessage)
      console.error('Preview error:', error)
    }
  })
}

export function useRunMigration() {
  const queryClient = useQueryClient()
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({ 
      envId, 
      name, 
      description, 
      sqlContent, 
      clientId 
    }: { 
      envId: string; 
      name: string; 
      description: string; 
      sqlContent: string; 
      clientId: string;
    }) => {
      const response = await client.POST('/environments/{id}/migrations', {
        params: { path: { id: envId } },
        body: {
          name,
          description,
          sql_content: sqlContent,
          client_id: clientId,
          environment_id: envId
        }
      })

      // Check if the response indicates an error
      if (response.error) {
  throw response.error;
      }

    if (response.data && (response.data as any).success === false) {
  throw response.data;
}
      
      return response.data
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: MIGRATIONS_QUERY_KEYS.list(variables.envId) })
      queryClient.invalidateQueries({ queryKey: ['environments', 'list'] }) 
      toast.success("Migration executed successfully!")
      navigate({
        to: "/app/projects/$projectId/environments/$envId",
        params: {
          projectId: variables.clientId,
          envId: variables.envId
        }
      })
    },
    onError: (error: any) => {
      // Extract detailed errors if available
      const errorMessage = error.details?.length 
        ? error.details.join('; ') 
        : error.message || 'Failed to run migration'
      
      toast.error(errorMessage)
      console.error('Migration error:', error)
    }
  })
}
    export function useValidateEnvironmentConnection() {
  return useMutation({
    mutationFn: async (envId: string) => {
      const response = await client.POST('/environments/{id}/validate', {
        params: { path: { id: envId } }
      })

      // Check if the response indicates an error
      if (response.error) {
        const error = new Error(response.error.message || 'Environment validation failed')
        throw error
      }

      return response.data
    },
    onSuccess: () => {
      toast.success("Environment connection validated successfully!")
    },
    onError: (error: any) => {
      const errorMessage = error.message || 'Failed to validate environment connection'
      toast.error(errorMessage)
      console.error('Validation error:', error)
    }
  })
}
