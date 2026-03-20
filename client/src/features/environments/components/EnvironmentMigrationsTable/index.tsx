import { Spinner } from "@/components/ui/spinner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { EnvironmentMigration } from "@/features/types";
import { StatusBadge } from "@/features/environments/components/StatusBadge";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { PreviewSQLModal } from "../PreviewSQLModal";
import { EmptyMigrationsState } from "../EmptyMigrationsState";
import { formatDate } from "@/utils/date";
import {
  AlertCircle,
  DatabaseIcon,
} from "lucide-react";

type Props = {
  migrations: EnvironmentMigration[];
  isLoading: boolean;
  onViewResults?: (migration: EnvironmentMigration) => void;
};

export function EnvironmentMigrationsTable({
  migrations,
  isLoading,
  onViewResults,
}: Props) {
  const [previewSqlOpen, setPreviewSqlOpen] = useState(false);
  const [previewSql, setPreviewSql] = useState("");

  return (
    <div className="rounded-2xl border overflow-hidden shadow-sm bg-background mt-0">
      {migrations?.length === 0 && !isLoading ? (
        <EmptyMigrationsState />
      ) : (
        <>
          <PreviewSQLModal
            open={previewSqlOpen}
            onOpenChange={setPreviewSqlOpen}
            sql={previewSql}
            name={"Preview SQL"}
          />
          <Table>
            <TableHeader className="bg-muted-foreground/5">
              <TableRow className="hover:bg-transparent border-stone-100">
                <TableHead className="text-xs font-medium text-black py-4 pl-6">
                  Node Metadata
                </TableHead>
                <TableHead className="text-xs font-medium text-black py-4">
                  Raw SQL Inspection
                </TableHead>
                <TableHead className="text-xs font-medium text-black py-4">
                  Status
                </TableHead>
                <TableHead className="text-xs font-medium text-black py-4 text-right pr-8">
                  Execution Time
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={4} className="py-10">
                    <Spinner />
                  </TableCell>
                </TableRow>
              ) : (
                migrations?.map((m) => (
                  <TableRow
                    key={m.id}
                    className="border-stone-50 transition-colors group hover:bg-stone-50/30"
                  >
                    <TableCell className="py-5 pl-6">
                      <div className="flex items-center gap-2">
                        <div
                          className={`font-medium text-sm ${m.status === "failed" ? "text-red-900" : "text-stone-900"}`}
                        >
                          {m.name}
                        </div>

                        {m.status === "success" && onViewResults && (
                          <Button
                            variant="ghost"
                            className="h-6 px-2 text-xs font-medium text-green-600 bg-green-50/50 hover:bg-green-100 hover:text-green-700 transition-all rounded-md border border-green-100/50 ml-1"
                            onClick={() => onViewResults(m)}
                          >
                            View Results
                          </Button>
                        )}
                      </div>
                      <div className="text-xs text-muted-foreground mt-0.5 text-wrap max-w-[400px] wrap-break-word">
                        {m?.description || "No description"}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 px-2.5 gap-2 text-zinc-500 hover:text-zinc-900 hover:bg-zinc-100/80 transition-all border border-transparent hover:border-zinc-200"
                        onClick={() => {
                          setPreviewSqlOpen(true);
                          setPreviewSql(m.sql_content || "");
                        }}
                      >
                        <DatabaseIcon
                          size={14}
                          strokeWidth={2.5}
                          className="opacity-70"
                        />
                        <span className="text-[11px] font-mono font-bold tracking-tight">
                          SQL
                        </span>
                      </Button>
                    </TableCell>

                    <TableCell className="py-3">
                      {m.status === "failed" ? (
                        <Popover>
                          <PopoverTrigger asChild>
                            <button className="outline-none focus:ring-2 focus:ring-red-500/20 rounded-full transition-transform active:scale-95">
                              <StatusBadge status={m.status} />
                            </button>
                          </PopoverTrigger>
                          <PopoverContent className="w-[480px] p-0 overflow-hidden shadow-2xl border-red-200">
                            <div className="bg-red-50/50 px-3 py-2 border-b border-red-100 flex items-center justify-between">
                              <div className="flex items-center gap-2">
                                <AlertCircle className="h-3.5 w-3.5 text-red-600" />
                                <span className="text-[10px] font-black uppercase tracking-widest text-red-600">
                                  Engine Exception
                                </span>
                              </div>
                              <span className="text-[9px] font-mono text-red-400">
                                SQL_RUNTIME_ERROR
                              </span>
                            </div>
                            <div className="p-4 bg-white">
                              <div className="bg-slate-950 p-3 rounded-lg overflow-x-auto border border-slate-800">
                                <code className="text-[11px] font-mono text-red-400 leading-relaxed block whitespace-pre-wrap">
                                  {m.error_message ||
                                    "No stack trace available for this exception."}
                                </code>
                              </div>
                            </div>
                          </PopoverContent>
                        </Popover>
                      ) : (
                        <StatusBadge status={m.status} />
                      )}
                    </TableCell>

                    <TableCell className="text-right pr-8">
                      <div className="text-xs font-bold text-foreground tabular-nums">
                        {m.duration ? `${m.duration}ms` : "--"}
                      </div>
                      <div className="text-xs text-muted-foreground mt-0.5 whitespace-nowrap">
                        {m.executed_at
                          ? formatDate(m.executed_at)
                          : "Scheduled"}
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </>
      )}
    </div>
  );
}
