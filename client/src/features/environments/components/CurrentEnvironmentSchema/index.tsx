import { useMemo } from "react";
import { Database, Hash, Lock, Table as TableIcon } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableRow, TableCell } from "@/components/ui/table";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import type { DatabaseSchemaItem } from "@/features/types";
import { LoadingState } from "@/components/shared/LoadingState";

interface CurrentDatabaseSchemaProps {
  schema: DatabaseSchemaItem;
  isLoading?: boolean;
}

export function CurrentDatabaseSchema({
  schema,
  isLoading,
}: CurrentDatabaseSchemaProps) {
  if (isLoading) {
    return <LoadingState />;
  }

  if (!schema || Object.keys(schema).length === 0) {
    return (
      <div className="h-[200px] flex items-center justify-center border border-dashed rounded-lg bg-stone-50/10">
        <div className="text-center">
          <Database className="h-6 w-6 text-stone-300 mx-auto mb-2" />
          <p className="text-xs text-stone-400 font-mono">
            No schema detected in this environment
          </p>
        </div>
      </div>
    );
  }

  const groupedByTable = useMemo(() => {
    // schema is an object where keys are table names and values are arrays of columns
    return schema as Record<string, any[]>;
  }, [schema]);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 p-8 bg-stone-50/50">
      {Object.entries(groupedByTable).map(([tableName, columns]) => (
        <Card
          key={tableName}
          className="shadow-sm overflow-hidden p-0 border-stone-200 bg-white group/card hover:shadow-md transition-all duration-300"
        >
          <CardHeader className="border-b bg-stone-50/50 px-4 py-3 group-hover/card:bg-stone-100/50 transition-colors">
            <CardTitle className="flex items-center justify-between">
              <div className="flex items-center gap-2.5">
                <div className="p-1.5 bg-blue-50 rounded-md border border-blue-100">
                  <TableIcon className="h-3.5 w-3.5 text-blue-600" />
                </div>
                <span className="font-bold text-sm text-stone-900 tracking-tight">
                  {tableName}
                </span>
              </div>
              <Badge variant="secondary" className="text-[10px]">
                {columns.length} Columns
              </Badge>
            </CardTitle>
          </CardHeader>

          <CardContent className="p-0">
            <Table>
              <TableBody>
                {columns.map((col, index) => {
                  const isPK =
                    col.column_name === "id" || col.column_name.endsWith("_id");

                  return (
                    <TableRow
                      key={index}
                      className="hover:bg-blue-50/30 border-stone-100 last:border-0 group/row transition-colors"
                    >
                      <TableCell className="py-2.5 px-4 flex items-center gap-3">
                        {isPK ? (
                          <Hash className="h-3.5 w-3.5 text-amber-500 shrink-0 stroke-[2.5px]" />
                        ) : (
                          <div className="w-3.5 h-3.5 flex items-center justify-center shrink-0">
                            <div className="w-1.5 h-1.5 rounded-full bg-stone-200 group-hover/row:bg-blue-400 transition-colors" />
                          </div>
                        )}
                        <span
                          className={`text-[12px] font-medium tracking-tight ${isPK ? "text-stone-900 font-bold" : "text-stone-800"}`}
                        >
                          {col.column_name}
                        </span>
                      </TableCell>

                      <TableCell className="py-2.5 px-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <code className="text-[10px] font-bold text-blue-700 bg-blue-50/50 px-1.5 py-0.5 rounded border border-blue-100/50 uppercase tracking-tighter">
                            {col.data_type
                              .replace("without time zone", "")
                              .trim()}
                          </code>
                          {col.is_nullable === "NO" && (
                            <Lock className="h-3 w-3 text-stone-400 stroke-[2.5px]" />
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
