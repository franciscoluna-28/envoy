import { requireAuth } from '@/utils/guard'
import { createFileRoute } from '@tanstack/react-router'
import { useAuthStore } from '@/store/auth'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { LogOut, ShieldCheck } from 'lucide-react'
import { useLogoutMutation, useMe } from '@/features/auth/hooks'

export const Route = createFileRoute('/app/')({
  component: RouteComponent,
  beforeLoad: requireAuth
})

function RouteComponent() {
  const { user } = useAuthStore()
  const { isLoading, error } = useMe();

  const { mutate: logout} = useLogoutMutation()

  if (isLoading) {
    return (
      <main className="h-screen w-full flex items-center justify-center p-4">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </main>
    )
  }

  if (error) {
    return (
      <main className="h-screen w-full flex items-center justify-center p-4">
        <Card className="w-full max-w-sm">
          <CardHeader className="text-center">
            <CardTitle className="text-xl text-destructive">Error</CardTitle>
            <CardDescription>
              Failed to load user data
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button className='w-full' onClick={() => logout()}>
              Try Again
            </Button>
          </CardContent>
        </Card>
      </main>
    )
  }

  return (
    <main className="h-screen w-full flex items-center justify-center p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex size-12 items-center justify-center rounded-full border">
            <ShieldCheck className="size-6" />
          </div>
          <CardTitle className="text-2xl font-bold tracking-tight">Access Granted</CardTitle>
          <CardDescription className="">
            The app works lol.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="rounded-lg border p-6">
            <p className="text-xs uppercase tracking-wider font-semibold mb-1">Authenticated as</p>
            <p className="text-sm font-mono truncate">{user?.email}</p>
          </div>
          
          <Button className='w-full'
            onClick={() => logout()}
          >
            <LogOut className="size-4" /> 
            Logout
          </Button>
        </CardContent>
      </Card>
    </main>
  )
}