import { useMutation } from "@tanstack/react-query";
import type { LoginInputDto, RegisterInputDto } from "../types";
import client from "@/api/client";
import { useAuthStore } from "@/store/auth";
import { toast } from "sonner";
import { useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";

export const useRegisterMutation = () => {
  return useMutation({
    mutationFn: async (data: RegisterInputDto) => {
      const response = await client.POST("/auth/register", {
        body: data,
      });

      if (response.error) {
        const errorDetail =
          (response.error as any).detail || "Registration failed";
        throw new Error(errorDetail);
      }
      return response.data;
    },
  });
};

export const useLoginMutation = () => {
  const authStore = useAuthStore();

  return useMutation({
    mutationFn: async (data: LoginInputDto) => {
      const response = await client.POST("/auth/login", {
        body: data,
      });

      if (response.error) {
        const errorDetail = (response.error as any).detail || "Login failed";
        throw new Error(errorDetail);
      }
      return response.data;
    },

    onSuccess: (response) => {
      if (response) {
        const { id, email, created_at } = response;

      if(id && email && created_at) { 
        authStore.setAuth({
          id,
          email,
          created_at,
        });
      }
    }
    },
  });
};

export const useLogoutMutation = () => {
  const authStore = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async () => {
      const response = await client.POST("/auth/logout");
      if (response.error) {
        throw new Error("Logout failed");
      }
    },

    onSuccess: (_response) => {
      authStore.logout();
      navigate({
        to: "/"
      })
    },

    onError: () => {
      toast.error("Error while logging out");
    },
  });
};

export const useMe = () => {
  const authStore = useAuthStore();

  return useQuery({
    queryKey: ["auth", "me"],
    queryFn: async () => {
      const response = await client.GET("/auth/me");
      
      if (response.error) {
       if (response.response?.status === 401) {
          authStore.logout();
        } 
        const errorDetail = (response.error as any).detail || "Failed to fetch user";
        throw new Error(errorDetail);
      }

      return response.data;
    },
    enabled: authStore.isAuthenticated,
  });
};