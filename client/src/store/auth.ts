import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

interface UserData {
  id: string
  email: string
  created_at: string
}

interface AuthState {
  user: UserData | null
  isAuthenticated: boolean
  setAuth: (data: UserData) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,

      setAuth: (data) => 
        set({ 
          user: data, 
          isAuthenticated: true 
        }),

      logout: () => 
        set({ 
          user: null, 
          isAuthenticated: false 
        }),
    }),
    {
      name: 'envoy-auth-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
)