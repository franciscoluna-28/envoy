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
import { formatDistanceToNow } from "date-fns";
import { useState } from "react";
import { PreviewSQLModal } from "../PreviewSQLModal";
import { EmptyMigrationsState } from "../EmptyMigrationsState";

type Props = {
  migrations: EnvironmentMigration[];
  isLoading: boolean;
  onViewResults?: (migration: EnvironmentMigration) => void;
};


export function EnvironmentMigrationsTable({ migrations, isLoading, onViewResults }: Props) {
  const [previewSqlOpen, setPreviewSqlOpen] = useState(false);
  const [previewSql, setPreviewSql] = useState("");

  return (
    <div className="rounded-2xl border overflow-hidden shadow-sm bg-background mt-0">
      {migrations?.length === 0 && !isLoading ? (
       <EmptyMigrationsState />
      ) : (
        <>
          <PreviewSQLModal open={previewSqlOpen} onOpenChange={setPreviewSqlOpen} sql={previewSql} name={"Preview SQL"} />
        <Table>
          <TableHeader className="bg-muted-foreground/5">
            <TableRow className="hover:bg-transparent border-stone-100">
              <TableHead className="text-xs font-medium text-black py-4 pl-6">
                Node Metadata
              </TableHead>
              <TableHead className="text-xs font-medium text-black py-4">
                View SQL
              </TableHead>
              <TableHead className="text-xs font-medium text-black py-4">
                Status
              </TableHead>
              <TableHead className="text-xs font-medium text-black py-4 text-right pr-8">
                Performance / Time
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
                      {m.status === "failed" && (
                        <Popover>
                          <PopoverTrigger asChild>
                            <Button
                              variant="ghost"
                              className="h-6 px-2 text-xs font-medium text-red-500 bg-red-50/50 hover:bg-red-100 hover:text-red-700 transition-all rounded-md border border-red-100/50 ml-1"
                            >
                              View Logs
                            </Button>
                          </PopoverTrigger>
                          <PopoverContent className="w-[450px] p-0 overflow-hidden border-stone-800 shadow-2xl bg-stone-950">
                            <div className="bg-stone-900 px-3 py-2 border-b border-stone-800 flex items-center justify-between">
                              <span className="text-[10px] font-bold text-red-400 uppercase tracking-widest">
                                Engine Exception
                              </span>
                              <code className="text-[9px] text-stone-500 font-mono">
                                ID: {m?.id?.slice(0, 8)}
                              </code>
                            </div>
                            <div className="p-4">
                              <code className="text-[11px] font-mono text-red-400 leading-relaxed block whitespace-pre-wrap">
                                {m.error_message}
                              </code>
                            </div>
                          </PopoverContent>
                        </Popover>
                      )}
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
                    <div className="text-xs text-stone-400 mt-0.5">
                      {m.description || "System Migration"}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() => {
                        setPreviewSqlOpen(true);
                        setPreviewSql(m.sql_content || "");
                      }}
                    >
                      View SQL
                    </Button>
                  </TableCell>

                  <TableCell>
                    <StatusBadge status={m.status} />
                  </TableCell>

                  <TableCell className="text-right pr-8">
                    <div className="text-xs font-bold text-stone-900 tabular-nums">
                      {m.duration ? `${m.duration}ms` : "--"}
                    </div>
                    <div className="text-[10px] text-stone-400 font-bold uppercase tracking-tighter mt-0.5 whitespace-nowrap">
                      {m.executed_at
                        ? formatDistanceToNow(new Date(m.executed_at), {
                            addSuffix: true,
                          })
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

