import { 
  DatabaseSchemaNode, 
  DatabaseSchemaNodeHeader, 
  DatabaseSchemaNodeBody, 
  DatabaseSchemaTableRow, 
  DatabaseSchemaTableCell 
} from '@/components/database-schema-node'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'

interface DatabaseSchemaProps {
  schema: any
  isLoading?: boolean
}

export function DatabaseSchema({ schema, isLoading }: DatabaseSchemaProps) {
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
        <DatabaseSchemaNode key={tableName}>
          <DatabaseSchemaNodeHeader>
            {tableName}
          </DatabaseSchemaNodeHeader>
          <DatabaseSchemaNodeBody>
            {(columns as any[]).map((column: any, index: number) => (
              <DatabaseSchemaTableRow key={`${column.column_name}-${index}`}>
                <DatabaseSchemaTableCell className="font-medium text-sm">
                  {column.column_name}
                </DatabaseSchemaTableCell>
                <DatabaseSchemaTableCell>
                  <Badge variant="outline" className="text-xs">
                    {column.data_type}
                  </Badge>
                </DatabaseSchemaTableCell>
                <DatabaseSchemaTableCell>
                  <Badge 
                    variant={column.is_nullable === 'YES' ? 'secondary' : 'default'}
                    className="text-xs"
                  >
                    {column.is_nullable === 'YES' ? 'NULL' : 'NOT NULL'}
                  </Badge>
                </DatabaseSchemaTableCell>
              </DatabaseSchemaTableRow>
            ))}
          </DatabaseSchemaNodeBody>
        </DatabaseSchemaNode>
      ))}
    </div>
  )
}
