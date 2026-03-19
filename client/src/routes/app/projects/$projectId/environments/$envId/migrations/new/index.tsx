import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import {
  CheckCircle2,
  AlertTriangle,
  Terminal,
  Loader2,
  XCircle,
} from "lucide-react";
import { createFileRoute } from "@tanstack/react-router";
import { Label } from "@/components/ui/label";
import { SQLEditor } from "@/features/environments/components/SQLEditor";
import { useState } from "react";
import { Separator } from "@/components/ui/separator";
import { useGetEnvironment } from "@/features/environments/hooks/useEnvironments";
import { useGetProject } from "@/features/projects/hooks/useProjects";
import {
  usePreviewSchemaChanges,
  useRunMigration,
  useValidateEnvironmentConnection,
} from "@/features/environments/hooks/useMigrations";
import { toast } from "sonner";
import { ScrollArea } from "@/components/ui/scroll-area";
import { LoadingState } from "@/components/shared/LoadingState";
import type { PreviewError } from "@/features/types";

export const Route = createFileRoute(
  "/app/projects/$projectId/environments/$envId/migrations/new/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId, envId } = Route.useParams();
  const [sqlValue, setSqlValue] = useState("");
  const [migrationName, setMigrationName] = useState("");
  const [description, setDescription] = useState("");
  const [permissionTestEnabled, setPermissionTestEnabled] = useState(false);
  const { isLoading: projectLoading } = useGetProject(projectId);
  const { isLoading: envLoading } = useGetEnvironment(projectId, envId);
  const previewSchema = usePreviewSchemaChanges();
  const runMigration = useRunMigration();
  const validateConnection = useValidateEnvironmentConnection();
  const isLoading = projectLoading || envLoading;

  if (isLoading) {
    return <LoadingState />;
  }

  return (
    <div className="flex flex-col h-full bg-background text-foreground antialiased">
      <div className="flex justify-between items-center py-6">
        <h1 className="text-xl font-bold tracking-tight">New Migration</h1>
      </div>

      <div className="flex-1 grid grid-cols-12 gap-6">
        <div className="col-span-8 flex flex-col h-[calc(100vh-180px)]">
          <div className="flex flex-col h-full rounded-lg border shadow-sm overflow-hidden bg-background">
            <div className="bg-zinc-50/80 border-b px-3 py-2 flex items-center justify-between shrink-0">
              <div className="flex items-center gap-4">
                <div className="flex gap-1.5">
                  <div className="h-2.5 w-2.5 rounded-full bg-red-400/20 border border-red-400/40" />
                  <div className="h-2.5 w-2.5 rounded-full bg-amber-400/20 border border-amber-400/40" />
                  <div className="h-2.5 w-2.5 rounded-full bg-emerald-400/20 border border-emerald-400/40" />
                </div>
                <div className="flex items-center gap-2">
                  <Terminal className="h-3.5 w-3.5 text-zinc-400" />
                  <span className="text-[11px] font-semibold tracking-wide uppercase text-zinc-500">
                    migration.sql
                  </span>
                </div>
              </div>
              <Badge
                variant="outline"
                className="text-[10px] font-mono font-medium px-1.5 py-0 h-5 bg-white"
              >
                {sqlValue.split("-->").length} Statements
              </Badge>
            </div>
            <div className="flex-1 min-h-[300px] relative">
              <div className="absolute inset-0">
                <SQLEditor value={sqlValue} onChange={setSqlValue} />
              </div>
            </div>
            <div className="px-3 py-1.5 bg-zinc-50/50 border-t flex items-center justify-between shrink-0">
              <div className="flex items-center gap-3 text-[10px] text-zinc-400 font-mono">
                <span className="flex items-center gap-1">
                  <span className="text-zinc-300">L:</span>
                  {sqlValue.split("\n").length}
                </span>
                <span className="flex items-center gap-1">
                  <span className="text-zinc-300">S:</span>
                  {(new Blob([sqlValue]).size / 1024).toFixed(2)} KB
                </span>
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 px-2 text-[11px] text-zinc-500 hover:text-destructive hover:bg-destructive/5"
                  onClick={() => setSqlValue("")}
                >
                  Clear
                </Button>
                <div className="h-3 w-px bg-zinc-200" />
              </div>
            </div>
          </div>
        </div>

        <aside className="col-span-4 flex flex-col space-y-4">
          <Card className="p-4 border">
           
            <div className="space-y-3">
              <div className="grid gap-1.5">
                <Label className="text-xs">Migration Name</Label>
                <Input
                  className="text-sm"
                  placeholder="add_performance_indexes"
                  value={migrationName}
                  onChange={(e) => setMigrationName(e.target.value)}
                />
              </div>

              <div className="grid gap-1.5">
                <Label className="text-xs">Description</Label>
                <Textarea
                  className="text-sm min-h-[80px]"
                  placeholder="Adding composite index to optimize feed query..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </div>
            </div>
          </Card>

          <Card className="p-4 border">
         
            <div className="space-y-3">
              <Button
                className="w-full text-sm"
                onClick={() => {
                  if (!sqlValue.trim()) {
                    toast.error("Please enter SQL content");
                    return;
                  }
                  previewSchema.mutate({ envId, sqlContent: sqlValue });
                }}
                disabled={previewSchema.isPending || !sqlValue.trim()}
              >
                {previewSchema.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Preview Schema Changes
              </Button>

              <div className="flex items-center justify-between p-3 rounded border border-dashed bg-muted/30 hover:bg-muted/50 transition-colors cursor-pointer"
                onClick={() => setPermissionTestEnabled(!permissionTestEnabled)}
              >
                <div className="space-y-0.5">
                  <p className="text-xs font-semibold">Permission Test</p>
                  <p className="text-[10px] text-muted-foreground">
                    Test against least privilege users
                  </p>
                </div>
                <Switch
                  checked={permissionTestEnabled}
                  onCheckedChange={setPermissionTestEnabled}
                  onClick={(e) => e.stopPropagation()}
                />
              </div>

              {permissionTestEnabled && (
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full text-xs"
                  onClick={() => {
                    validateConnection.mutate(envId);
                  }}
                  disabled={validateConnection.isPending}
                >
                  {validateConnection.isPending && (
                    <Loader2 className="mr-2 h-3 w-3 animate-spin" />
                  )}
                  Validate Connection
                </Button>
              )}
            </div>
          </Card>

          {/* Results Section - Simulation Output */}
          <Card className="p-4 bg-muted/30 border-dashed flex-1 flex flex-col">
            <div className="flex items-center justify-between mb-3">
              <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
                Simulation Results
              </span>
              {previewSchema.data ? (
                <CheckCircle2 className="h-3 w-3 text-emerald-600" />
              ) : previewSchema.error ? (
                <XCircle className="h-3 w-3 text-destructive" />
              ) : previewSchema.isPending ? (
                <Loader2 className="h-3 w-3 animate-spin text-muted-foreground" />
              ) : (
                <AlertTriangle className="h-3 w-3 text-muted-foreground" />
              )}
            </div>

            <ScrollArea className="flex-1 w-full rounded-md">
              <div className="space-y-2 pr-4">
                {previewSchema.error && (
                  <div className="mb-4 space-y-2">
                    <div className="flex items-center gap-2 text-[10px] font-bold text-red-600 uppercase tracking-tight">
                      <Terminal className="h-3 w-3" />
                      Simulation Error
                    </div>
                    {(previewSchema.error as PreviewError).errors?.map(
                      (err: string, i: number) => (
                        <div
                          key={i}
                          className="p-2 rounded bg-red-50 border border-red-100 font-mono text-[10px] text-red-800 wrap-break-word"
                        >
                          {err}
                        </div>
                      ),
                    )}
                  </div>
                )}
                {previewSchema.data && !previewSchema.error
                  ? previewSchema.data.map((column, index) => (
                      <div key={index}>
                        <ResultItem
                          label={`${column.table_name}.${column.column_name}`}
                          subLabel={column.data_type}
                        />
                        {index < previewSchema.data.length - 1 && (
                          <Separator className="my-2 opacity-30" />
                        )}
                      </div>
                    ))
                  : !previewSchema.isPending &&
                    !previewSchema.error && (
                      <div className="space-y-3 opacity-40">
                        <ResultItem
                          label="Enter SQL and click Preview"
                          isPending
                        />
                        <ResultItem label="Syntax Check Pending" isPending />
                      </div>
                    )}
                {previewSchema.isPending && (
                  <div className="flex items-center justify-center p-8 opacity-50">
                    <Loader2 className="h-4 w-4 animate-spin" />
                  </div>
                )}
              </div>
            </ScrollArea>
          </Card>

          {/* Execute Section - Final Action */}
          <Button
            size="lg"
            className="w-full"
            onClick={() => {
              if (!migrationName.trim()) {
                toast.error("Please enter a migration name");
                return;
              }
              if (!description.trim()) {
                toast.error("Please enter a description");
                return;
              }
              if (!sqlValue.trim()) {
                toast.error("Please enter SQL content");
                return;
              }
              runMigration.mutate({
                envId,
                name: migrationName,
                description: description,
                sqlContent: sqlValue,
                clientId: `migration_${Math.random().toString(36).substr(2, 9)}`,
              });
            }}
            disabled={
              runMigration.isPending ||
              !migrationName.trim() ||
              !description.trim() ||
              !sqlValue.trim()
            }
          >
            {runMigration.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            Apply Migration
          </Button>
        </aside>
      </div>
    </div>
  );
}

function ResultItem({
  label,
  subLabel,
  isPending,
}: {
  label: string;
  subLabel?: string;
  isPending?: boolean;
}) {
  return (
    <div className="flex items-center gap-3 text-[11px] font-mono py-0.5">
      {isPending ? (
        <div className="h-1.5 w-1.5 rounded-full bg-stone-300 animate-pulse" />
      ) : (
        <CheckCircle2 className="h-3 w-3 text-emerald-500 shrink-0" />
      )}
      <div className="flex flex-col truncate">
        <span className="text-stone-800 font-bold truncate">{label}</span>
        {subLabel && (
          <span className="text-[9px] text-muted-foreground uppercase">
            {subLabel}
          </span>
        )}
      </div>
    </div>
  );
}
