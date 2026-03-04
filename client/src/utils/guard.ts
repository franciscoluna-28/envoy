import { redirect } from '@tanstack/react-router'
import { useAuthStore } from '@/store/auth'

export const requireAuth = () => {
  const { isAuthenticated } = useAuthStore.getState()
  
  if (!isAuthenticated) {
    throw redirect({
      to: '/login',
      search: {
        redirect: window.location.pathname, 
      },
    })
  }
}

export const requireGuest = () => {
  const { isAuthenticated } = useAuthStore.getState()
  
  if (isAuthenticated) {
    throw redirect({
      to: '/app',
    })
  }
}