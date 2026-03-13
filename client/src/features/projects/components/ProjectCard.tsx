import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { MoreHorizontal, ArrowRight } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'

interface ProjectCardProps {
  project: {
    id?: string
    name?: string
  }
  onUpdate?: () => void
  onDelete?: () => void
}

export function ProjectCard({ project, onUpdate, onDelete }: ProjectCardProps) {
  return (
    <Card className="group relative flex flex-col justify-between border-input hover:border-foreground/30 transition-all duration-200 bg-card shadow-sm hover:shadow-md">
      <CardHeader className="pb-4">
        <div className="flex items-start justify-between">
          <CardTitle className="text-lg font-medium tracking-tight group-hover:underline underline-offset-4">
            {project.name || 'Untitled Project'}
          </CardTitle>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-[160px]">
              <DropdownMenuItem onClick={onUpdate}>
                Update Project
              </DropdownMenuItem>
              <DropdownMenuItem 
                className="text-destructive"
                onClick={onDelete}
              >
                Delete Project
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>

      <CardContent>
        <div className="flex items-center gap-2 mb-6">
          <div className="h-1.5 w-1.5 rounded-full bg-foreground shadow-[0_0_5px_rgba(0,0,0,0.2)]" />
          <span className="text-xs font-medium text-foreground">Production Active</span>
        </div>
        
        <div className="flex items-center justify-between border-t pt-4">
          <div className="text-[10px] font-mono text-muted-foreground uppercase flex gap-3">
            <span>0 ENVS</span>
            <span>•</span>
            <span>Updated 2d ago</span>
          </div>
          
          <Link 
            to="/app/projects/$projectId" 
            params={{ projectId: project.id || '' }} 
            className="text-xs font-semibold flex items-center gap-1 hover:gap-2 transition-all"
          >
            DETAILS
            <ArrowRight className="w-3 h-3" />
          </Link>
        </div>
      </CardContent>
    </Card>
  )
}
