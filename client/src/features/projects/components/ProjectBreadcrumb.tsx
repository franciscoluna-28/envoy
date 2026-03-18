import { Link, useParams } from '@tanstack/react-router'
import { useGetProject } from '@/features/projects'
import { useGetEnvironment } from '@/features/projects/hooks/useEnvironments'
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Home, Database, FileText } from 'lucide-react'

interface ProjectBreadcrumbProps {
  currentPage?: 'project' | 'environment' | 'migrations' | 'new-migration'
  currentPageName?: string
}

export function ProjectBreadcrumb({ currentPage, currentPageName }: ProjectBreadcrumbProps) {
  const params = useParams({ strict: false })
  const { projectId, envId } = params as { projectId?: string; envId?: string }

  // Always call hooks, but handle undefined cases
  const { data: project } = useGetProject(projectId || 'placeholder')
  const { data: environment } = useGetEnvironment(projectId || 'placeholder', envId || 'placeholder')

  const getIcon = (type: string) => {
    switch (type) {
      case 'home':
        return <Home className="h-3 w-3" />
      case 'project':
        return <Database className="h-3 w-3" />
      case 'environment':
        return <Database className="h-3 w-3" />
      case 'migrations':
        return <FileText className="h-3 w-3" />
      case 'new-migration':
        return <FileText className="h-3 w-3" />
      default:
        return null
    }
  }

  const getBreadcrumbItems = () => {
    const items = [
      {
        label: 'Projects',
        href: '/app',
        icon: getIcon('home'),
        isCurrent: false,
      },
    ]

    if (projectId && project) {
      items.push({
        label: project?.name || '',
        href: projectId ? `/app/projects/${projectId}` : '',
        icon: getIcon('project'),
        isCurrent: currentPage === 'project' || false,
      })
    }

    if (envId && environment && projectId) {
      items.push({
        label: environment.name,
        href: `/app/projects/${projectId}/environments/${envId}`,
        icon: getIcon('environment'),
        isCurrent: currentPage === 'environment' || false,
      })
    }

    if (currentPageName) {
      items.push({
        label: currentPageName,
        href: '',
        icon: getIcon(currentPage?.toLowerCase() || ''),
        isCurrent: true,
      })
    }

    return items
  }

  const breadcrumbItems = getBreadcrumbItems()

  return (
    <Breadcrumb className="text-xs">
      <BreadcrumbList>
        {breadcrumbItems.map((item, index) => (
          <div key={item.label} className="flex items-center gap-1.5">
            <BreadcrumbItem>
              {item.isCurrent ? (
                <BreadcrumbPage className="flex items-center gap-1.5 text-black font-medium">
                  {item.icon}
                  {item.label}
                </BreadcrumbPage>
              ) : (
                <BreadcrumbLink 
                  asChild
                  className="flex items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors"
                >
                  <Link to={item.href}>
                    {item.icon}
                    {item.label}
                  </Link>
                </BreadcrumbLink>
              )}
            </BreadcrumbItem>
            {index < breadcrumbItems.length - 1 && (
              <BreadcrumbSeparator />
            )}
          </div>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  )
}
