import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Database } from 'lucide-react'

interface SimpleGridSchemaProps {
  schema: any
  isLoading?: boolean
}

export function SimpleGridSchema({ schema, isLoading }: SimpleGridSchemaProps) {
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="text-center">
          <Database className="h-12 w-12 text-muted-foreground animate-pulse mx-auto mb-4" />
          <p className="text-muted-foreground">Loading database schema...</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {[1, 2].map((i) => (
            <div key={i} className="border rounded-lg p-4">
              <Skeleton className="h-6 w-32 mb-4" />
              <div className="space-y-2">
                {[1, 2, 3].map((j) => (
                  <div key={j} className="flex justify-between items-center">
                    <Skeleton className="h-4 w-24" />
                    <Skeleton className="h-4 w-16" />
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (!schema || schema.length === 0) {
    return (
      <div className="text-center py-12">
        <Database className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
        <p className="text-muted-foreground">No schema data available</p>
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

  const tables = Object.entries(groupedByTable)

  return (
    <div className="space-y-6">
      <div className="text-center mb-8">
        <Database className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
        <h2 className="text-2xl font-semibold">Database Schema</h2>
        <p className="text-muted-foreground">
          {tables.length} table{tables.length !== 1 ? 's' : ''} • {schema.length} total column{schema.length !== 1 ? 's' : ''}
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
        {tables.map(([tableName, columns]) => (
          <div key={tableName} className="border rounded-lg bg-white shadow-sm">
            <div className="bg-linear-to-r from-blue-50 to-indigo-50 px-4 py-3 border-b">
              <h3 className="font-semibold text-gray-900 flex items-center gap-2">
                <Database className="h-4 w-4 text-blue-600" />
                {tableName}
              </h3>
              <p className="text-xs text-gray-500 mt-1">
                {(columns as any[]).length} columns
              </p>
            </div>
            
            <div className="p-4 space-y-3">
              {(columns as any[]).map((column: any, index: number) => (
                <div key={index} className="flex justify-between items-center py-2 border-b border-gray-100 last:border-0">
                  <div className="flex-1">
                    <div className="font-mono text-sm text-gray-900">
                      {column.column_name}
                    </div>
                    {column.column_name.toLowerCase().includes('id') && (
                      <div className="text-xs text-amber-600 mt-1">Primary Key</div>
                    )}
                  </div>
                  <div className="flex items-center gap-2 ml-4">
                    <Badge variant="outline" className="text-xs">
                      {column.data_type}
                    </Badge>
                    {column.is_nullable === 'YES' && (
                      <Badge variant="secondary" className="text-xs">
                        nullable
                      </Badge>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
