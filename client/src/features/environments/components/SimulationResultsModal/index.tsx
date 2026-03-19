import { Database, Plus, Hash, Lock } from "lucide-react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableBody, TableRow, TableCell } from "@/components/ui/table";
import type { DatabaseSchemaItem, SchemaColumn } from "@/features/types";

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  simulationResults: DatabaseSchemaItem[];
  currentSchema: DatabaseSchemaItem[];
  onClose: () => void;
};

export function SimulationResultsModal({ 
  open, 
  onOpenChange, 
  simulationResults, 
  currentSchema = [], 
  onClose 
}: Props) {
  const groupedByTable = simulationResults.reduce((acc, column) => {
    const tableName = column.table_name || 'unknown_table';
    if (!acc[tableName]) {
      acc[tableName] = [];
    }
    acc[tableName].push(column);
    return acc;
  }, {} as Record<string, SchemaColumn[]>);

  const isNewTable = (tableName: string): boolean => {
    if (!currentSchema || currentSchema.length === 0) return true;
    return !currentSchema.some(col => col.table_name === tableName);
  };

  const isNewColumn = (tableName: string, columnName: string): boolean => {
    if (!currentSchema || currentSchema.length === 0) return true;
    return !currentSchema.some(col => 
      col.table_name === tableName && col.column_name === columnName
    );
  };

  // Sort tables: new tables first, then alphabetically
  const sortedTables = Object.entries(groupedByTable).sort(([a], [b]) => {
    const aIsNew = isNewTable(a);
    const bIsNew = isNewTable(b);
    
    if (aIsNew && !bIsNew) return -1;
    if (!aIsNew && bIsNew) return 1;
    
    return a.localeCompare(b);
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-[95vw] lg:max-w-6xl h-[85vh] flex flex-col p-0 overflow-hidden border-stone-200 rounded-3xl shadow-2xl">
        <DialogHeader className="p-6 border-b bg-background backdrop-blur-md sticky top-0 z-20">
          <div className="flex items-center justify-between">
            <div>
              <DialogTitle className="text-xl font-bold tracking-tight text-stone-900 flex items-center gap-2">
                Migration Simulation Results
              </DialogTitle>
              <p className="text-[10px] text-stone-400 font-bold uppercase">
                Target Schema Preview • {Object.keys(groupedByTable).length} Tables
              </p>
            </div>
          </div>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto overflow-x-hidden p-0 relative group scrollbar-thin scrollbar-track-stone-100 scrollbar-thumb-stone-300 hover:scrollbar-thumb-stone-400">
          <div className="absolute inset-0">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 p-8 bg-stone-50/50">
              {sortedTables.map(([tableName, columns]) => {
                const isTableNew = isNewTable(tableName);
                
                return (
                  <Card
                    key={tableName}
                    className={`shadow-sm overflow-hidden space-y-0! m-0! p-0 border group/card hover:shadow-md transition-all duration-300 flex gap-0 ${isTableNew ? 'border-green-400' : ''}`}
                  >
                    <CardHeader 
                      className="border-b px-4 py-3 group-hover/card:bg-stone-100/50 transition-colors"
                    >
                      <CardTitle className="flex items-center justify-between">
                        <div className="flex items-center gap-2.5">
                          <div className="p-1.5 rounded-md border">
                            <Database className="h-3.5 w-3.5 text-blue-600" />
                          </div>
                          <span className="font-bold text-sm text-stone-900 tracking-tight">
                            {tableName}
                          </span>
                          {isTableNew && (
                            <Badge variant="secondary" className="text-[10px] bg-green-100 text-green-700 border-green-200">

                              NEW
                            </Badge>
                          )}
                        </div>
                        <Badge variant="secondary" className="text-[10px]">
                          {columns.length} Columns
                        </Badge>
                      </CardTitle>
                    </CardHeader>

                    <CardContent className="p-0! m-0! space-y-0">
                      <Table>
                        <TableBody>
                          {columns.map((col, index) => {
                            const isPK = col.column_name === "id" || col.column_name?.endsWith("_id");
                            const isColumnNew = isNewColumn(tableName, col.column_name || '');

                            return (
                              <TableRow
                                key={index}
                                className={`border-stone-100 last:border-0 transition-colors group/row`}
                              >
                                <TableCell className="py-2.5 px-4 flex items-center gap-3">
                                  {isPK ? (
                                    <Hash className="h-3.5 w-3.5 text-amber-500 shrink-0 stroke-[2.5px]" />
                                  ) : (
                                    <div className="w-3.5 h-3.5 flex items-center justify-center shrink-0">
                                      {isColumnNew ? (
                                        <Plus className="h-2.5 w-2.5 text-green-600" />
                                      ) : (
                                        <div className="w-1.5 h-1.5 rounded-full bg-stone-200 group-hover/row:bg-blue-400 transition-colors" />
                                      )}
                                    </div>
                                  )}
                                  <span
                                    className={`text-[12px] font-medium tracking-tight ${
                                      isPK ? "text-stone-900 font-bold" : "text-stone-800"
                                    }`}
                                  >
                                    {col.column_name || 'Unknown'}
                                  </span>
                                  {isColumnNew && (
                                    <Badge variant="secondary" className="text-[9px] bg-green-100 text-green-700 border-green-200">
                                      NEW
                                    </Badge>
                                  )}
                                </TableCell>

                                <TableCell className="py-2.5 px-4 text-right">
                                  <div className="flex items-center justify-end gap-2">
                                    <code className="text-[10px] font-bold text-blue-700 bg-blue-50/50 px-1.5 py-0.5 rounded border border-blue-100/50 uppercase tracking-tighter">
                                      {col.data_type
                                        ?.replace("without time zone", "")
                                        ?.replace("with time zone", "")
                                        ?.trim() || 'unknown'}
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
                );
              })}
            </div>
          </div>
        </div>

        <div className="p-4 border-t bg-white/80 backdrop-blur-md flex items-center justify-end gap-3 sticky bottom-0 z-20">
          <div className="flex items-center gap-4 mr-auto">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded bg-green-100 border border-green-200"></div>
              <span className="text-[10px] text-stone-600 font-medium">New items</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded bg-stone-100 border border-stone-200"></div>
              <span className="text-[10px] text-stone-600 font-medium">Existing items</span>
            </div>
          </div>
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
