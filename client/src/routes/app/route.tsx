import { requireAuth } from "@/utils/guard";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { useAuthStore } from "@/store/auth";
import { LogOut, FolderOpen, ChevronDown } from "lucide-react";
import Logo from "@/assets/logo.png";
import { useLogoutMutation, useMe } from "@/features/auth/hooks";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
} from "@/components/ui/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Spinner } from "@/components/ui/spinner";

export const Route = createFileRoute("/app")({
  component: RouteComponent,
  beforeLoad: requireAuth,
});

function RouteComponent() {
  const { user } = useAuthStore();
  const { isLoading, error } = useMe();
  const { mutate: logout } = useLogoutMutation();

  return (
    <SidebarProvider>
      <Sidebar>
        <SidebarHeader>
          <div className="flex flex-col items-center gap-4 px-4 py-6">
            <div className="flex items-center justify-center">
              <img
                src={Logo}
                alt="Envoy"
                className="w-36 h-36 transition-transform hover:scale-105 bg-contain object-contain"
                onError={(e) => {
                  const target = e.target as HTMLImageElement;
                  target.style.display = "none";
                  target.nextElementSibling?.classList.remove("hidden");
                }}
              />
              <div className="hidden w-16 h-16 items-center justify-center bg-primary text-primary-foreground rounded-xl font-bold text-2xl">
                E
              </div>
            </div>
            <div className="text-center hidden">
              <h1 className="text-xl font-bold tracking-tight text-sidebar-foreground">
                Envoy
              </h1>
              <p className="text-sm text-muted-foreground">
                Management Console
              </p>
            </div>
          </div>
        </SidebarHeader>

        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Navigation</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      to="/app"
                      className="transition-all hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                    >
                      <FolderOpen className="w-4 h-4" />
                      <span>Projects</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter className="border-t">
          <SidebarMenu>
            <SidebarMenuItem>
              {error && (
                <p className="text-red-500 text-sm px-2 py-1">
                  {error.message || "Unknown error"}
                </p>
              )}
              {isLoading ? (
                <div className="flex items-center justify-center p-2">
                  <Spinner />
                </div>
              ) : (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <SidebarMenuButton className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
                      <span className="truncate font-medium">
                        {user?.email}
                      </span>
                      <ChevronDown className="ml-auto w-4 h-4" />
                    </SidebarMenuButton>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    className="w-48 rounded-lg"
                    side="bottom"
                    align="end"
                    sideOffset={4}
                  >
                    <DropdownMenuItem onClick={() => logout()}>
                      <LogOut className="w-4 h-4 mr-2" />
                      <span>Log out</span>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sidebar>

      <SidebarInset>
        <div className="flex flex-1 flex-col gap-4 p-4 mt-8">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
