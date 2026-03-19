import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import client from "@/api/client";
import { PROJECTS_QUERY_KEYS } from "@/features/keys";
import type { ProjectId } from "@/features/types";

export function useGetAllProjects() {
  return useQuery({
    queryKey: PROJECTS_QUERY_KEYS.lists(),
    queryFn: async () => {
      const response = await client.GET("/projects");
      return response.data;
    },
  });
}

export function useGetProject(projectId: ProjectId) {
  return useQuery({
    queryKey: PROJECTS_QUERY_KEYS.detail(projectId),
    queryFn: async () => {
      const response = await client.GET("/projects/{id}", {
        params: { path: { id: projectId } },
      });
      return response.data;
    },
    enabled: !!projectId,
  });
}

export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: { name: string }) => {
      const response = await client.POST("/projects", {
        body: data,
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() });
    },
  });
}

export function useUpdateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: ProjectId;
      data: { name: string };
    }) => {
      const response = await client.PUT("/projects/{id}", {
        params: { path: { id } },
        body: { id, ...data },
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() });
      queryClient.invalidateQueries({
        queryKey: PROJECTS_QUERY_KEYS.details(),
      });
    },
  });
}

export function useDeleteProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: ProjectId) => {
      const response = await client.DELETE("/projects/{id}", {
        params: { path: { id } },
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() });
      queryClient.invalidateQueries({
        queryKey: PROJECTS_QUERY_KEYS.details(),
      });
    },
  });
}
