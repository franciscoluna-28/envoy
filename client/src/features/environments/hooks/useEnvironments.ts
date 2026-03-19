import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import client from "@/api/client";
import { toast } from "sonner";
import { ENVIRONMENTS_QUERY_KEYS } from "@/features/keys";
import type {
    CreateEnvironmentInput,
  EnvironmentId,
  ProjectId,
} from "@/features/types";

export function useGetEnvironments(projectId: ProjectId) {
  return useQuery({
    queryKey: ENVIRONMENTS_QUERY_KEYS.list(projectId),
    queryFn: async () => {
      const response = await client.GET("/projects/{id}/environments", {
        params: { path: { id: projectId } },
      });
      return response.data;
    },
    enabled: !!projectId,
  });
}

export function useGetEnvironment(projectId: ProjectId, envId: EnvironmentId) {
  return useQuery({
    queryKey: [...ENVIRONMENTS_QUERY_KEYS.list(projectId), envId],
    queryFn: async () => {
      const response = await client.GET("/environments/{id}", {
        params: { path: { id: envId } },
      });
      return response.data;
    },
    enabled: !!projectId && !!envId,
  });
}

export function useCreateEnvironment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateEnvironmentInput) => {
      const response = await client.POST("/projects/{id}/environments", {
        params: { path: { id: data.project_id } },
        body: {
          name: data.name,
          connection_url: data.connection_url,
          project_id: data.project_id,
          type: data.type,
        },
      });
      return response.data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ENVIRONMENTS_QUERY_KEYS.list(variables.project_id),
      });
      toast.success("Environment created successfully!");
    },
    onError: () => {
      toast.error("Failed to create environment");
    },
  });
}
