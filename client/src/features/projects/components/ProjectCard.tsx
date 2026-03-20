import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { MoreHorizontal, Calendar, Layout } from "lucide-react";
import { Link } from "@tanstack/react-router";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { Project } from "@/features/types";
import { formatDate } from "@/utils/date";

interface ProjectCardProps {
  project: Project
  onUpdate?: () => void;
  onDelete?: () => void;
}

export function ProjectCard({ project, onUpdate, onDelete }: ProjectCardProps) {
  return (
    <Card className="group relative min-w-[350px] shadow-sm hover:shadow-md transition-all duration-300 border-stone-200 p-0">
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-8">
          <div className="space-y-1.5">
            <h3 className="text-lg font-semibold tracking-tight text-stone-900 transition-colors line-clamp-1 flex items-center">
              <Layout className="w-4 h-4 mr-2 text-stone-400  transition-colors" />
              {project.name || "Untitled Project"}
            </h3>
            <div className="flex items-center gap-2">
              <Badge
                variant="secondary"
              
              >
                PostgreSQL
              </Badge>

              <span className="text-stone-300 text-[10px]">•</span>

              <div className="flex items-center gap-1 text-xs">
                <Calendar className="w-3 h-3 text-muted-foreground" />
                <span>
                  {project.created_at
                    ? formatDate(project.created_at)
                    : "Just now"}
                </span>
              </div>
            </div>
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="h-8 w-8 p-0 hover:bg-stone-100 text-stone-400 data-[state=open]:bg-stone-100 data-[state=open]:text-stone-900"
              >
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent
              align="end"
              className="w-[180px] p-1.5 shadow-xl border-stone-200"
            >
              <DropdownMenuItem
                onClick={onUpdate}
                className="text-xs font-semibold cursor-pointer px-2.5 py-2 rounded-md transition-colors"
              >
                Update Project
              </DropdownMenuItem>

              <Separator className="my-1.5" />

              <DropdownMenuItem
                className="text-destructive text-xs font-semibold cursor-pointer px-2.5 py-2 rounded-md hover:bg-red-50 hover:text-red-600 transition-all"
                onClick={onDelete}
              >
                Delete Project
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="flex items-center">
          <Link
            to="/app/projects/$projectId"
            className="w-full"
            params={{ projectId: project.id || "" }}
          >
            <Button variant="secondary" size="lg" className="w-full">
              Open Project
            </Button>
          </Link>
        </div>
      </CardContent>
    </Card>
  );
}
