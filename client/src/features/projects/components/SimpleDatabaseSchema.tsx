import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

interface SimpleDatabaseSchemaProps {
  schema: any
  isLoading?: boolean
}

export function SimpleDatabaseSchema({ schema, isLoading }: SimpleDatabaseSchemaProps) {
  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-full" />
        <div className="space-y-2">
          {[1, 2, 3, 4].map((i) => (
            <Skeleton key={i} className="h-6 w-full" />
          ))}
        </div>
      </div>
    )
  }

  if (!schema || schema.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <p>No schema data available</p>
      </div>
    )
  }

  // Group schema data by table
  const groupedByTable = (schema || []).reduce((acc: Record<string, any[]>, item: any) => {
    const tableName = item.table_name || 'unknown'
    if (!acc[tableName]) {
      acc[tableName] = []
    }
    acc[tableName].push(item)
    return acc
  }, {})

  return (
    <div className="space-y-8">
      {Object.entries(groupedByTable).map(([tableName, columns]) => (
        <div key={tableName} className="border rounded-lg">
          <div className="bg-muted/50 px-4 py-3 border-b">
            <h3 className="font-semibold text-lg">{tableName}</h3>
            <p className="text-sm text-muted-foreground">
              {(columns as any[]).length} columns
            </p>
          </div>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[300px]">Column Name</TableHead>
                  <TableHead className="w-[150px]">Data Type</TableHead>
                  <TableHead className="w-[120px]">Nullable</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(columns as any[]).map((column: any, index: number) => (
                  <TableRow key={`${column.column_name}-${index}`}>
                    <TableCell className="font-mono text-sm">
                      {column.column_name}
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className="text-xs">
                        {column.data_type}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <span className={`text-xs px-2 py-1 rounded ${
                        column.is_nullable === 'YES' 
                          ? 'bg-yellow-100 text-yellow-800' 
                          : 'bg-green-100 text-green-800'
                      }`}>
                        {column.is_nullable === 'YES' ? 'YES' : 'NO'}
                      </span>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </div>
      ))}
    </div>
  )
}
