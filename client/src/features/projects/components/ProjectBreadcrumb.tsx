import { Link, useLocation, useParams } from '@tanstack/react-router'
import { useGetEnvironment } from '@/features/environments/hooks/useEnvironments'
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Home, Database, FileText, ChevronRight, PlusCircle, LayoutDashboard } from 'lucide-react'
import { useGetProject } from '../hooks/useProjects'
import { cn } from '@/lib/utils'
import React from 'react'

export function ProjectBreadcrumb() {
  const location = useLocation()
  const params = useParams({ strict: false })
  const { projectId, envId } = params as { projectId?: string; envId?: string }

  const { data: project, isLoading: loadingProject } = useGetProject(projectId || '')
  const { data: environment, isLoading: loadingEnv } = useGetEnvironment(
    projectId || '', 
    envId || ''
  )

  const path = location.pathname
  const isMigrationsPage = path.includes('/migrations') && !path.includes('/new')
  const isNewMigrationPage = path.includes('/migrations/new')

  const icons = {
    home: <Home className="h-3.5 w-3.5" />,
    project: <LayoutDashboard className="h-3.5 w-3.5" />,
    environment: <Database className="h-3.5 w-3.5" />,
    migrations: <FileText className="h-3.5 w-3.5" />,
    new: <PlusCircle className="h-3.5 w-3.5" />,
  }

  const getBreadcrumbItems = () => {
    const items = [
      {
        label: 'Projects',
        href: '/app',
        icon: icons.home,
        active: path === '/app',
      },
    ]

    if (projectId) {
      items.push({
        label: loadingProject ? 'Loading...' : (project?.name || 'Project'),
        href: `/app/projects/${projectId}`,
        icon: icons.project,
        active: path === `/app/projects/${projectId}`,
      })
    }

    if (envId && projectId) {
      items.push({
        label: loadingEnv ? 'Loading...' : (environment?.name || 'Environment'),
        href: `/app/projects/${projectId}/environments/${envId}`,
        icon: icons.environment,
        active: path === `/app/projects/${projectId}/environments/${envId}`,
      })
    }

    if (isMigrationsPage) {
      items.push({
        label: 'Migrations',
        href: `/app/projects/${projectId}/environments/${envId}`,
        icon: icons.migrations,
        active: true,
      })
    }

    if (isNewMigrationPage) {
      items.push({
        label: 'New Migration',
        href: '',
        icon: icons.new,
        active: true,
      })
    }

    return items
  }

  const breadcrumbItems = getBreadcrumbItems()

  return (
    <Breadcrumb className="text-[11px] font-medium tracking-tight">
      <BreadcrumbList className="gap-1 sm:gap-1">
        {breadcrumbItems.map((item, index) => (
          <React.Fragment key={`${item.label}-${index}`}>
            <BreadcrumbItem>
              {item.active ? (
                <BreadcrumbPage 
                  className={cn(
                    "flex items-center gap-2 px-2 py-1 rounded-md transition-all",
                    "text-stone-900 font-medium!"
                  )}
                >
                  {item.icon}
                  <span className="truncate max-w-[150px]">{item.label}</span>
                </BreadcrumbPage>
              ) : (
                <BreadcrumbLink 
                  asChild
                  className="flex items-center gap-2 font-normal px-2 py-1 rounded-md text-stone-500 hover:text-stone-900 hover:bg-stone-100/50 transition-colors"
                >
                  <Link to={item.href}>
                    {item.icon}
                    <span className="truncate max-w-[150px]">{item.label}</span>
                  </Link>
                </BreadcrumbLink>
              )}
            </BreadcrumbItem>
            {index < breadcrumbItems.length - 1 && (
              <BreadcrumbSeparator className="opacity-40">
                <ChevronRight className="h-3 w-3" />
              </BreadcrumbSeparator>
            )}
          </React.Fragment>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  )
}