import { Database } from 'lucide-react'

interface CompactGridSchemaProps {
  schema: any
  isLoading?: boolean
}

export function CompactGridSchema({ schema, isLoading }: CompactGridSchemaProps) {
  if (isLoading) {
    return (
      <div className="text-center py-12">
        <Database className="h-8 w-8 text-gray-400 animate-pulse mx-auto mb-2" />
        <p className="text-gray-500">Loading...</p>
      </div>
    )
  }

  if (!schema || schema.length === 0) {
    return (
      <div className="text-center py-12">
        <Database className="h-8 w-8 text-gray-400 mx-auto mb-2" />
        <p className="text-gray-500">No schema data</p>
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
    <div className="space-y-4">
      <div className="text-center mb-6">
        <Database className="h-6 w-6 text-gray-500 mx-auto mb-1" />
        <h2 className="text-lg font-medium text-gray-800">Database Schema</h2>
        <p className="text-xs text-gray-500">
          {tables.length} table{tables.length !== 1 ? 's' : ''} • {schema.length} column{schema.length !== 1 ? 's' : ''}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {tables.map(([tableName, columns]) => (
          <div key={tableName} className="border border-gray-200 bg-white">
            <div className="bg-gray-50 px-3 py-2 border-b border-gray-200">
              <h3 className="font-medium text-sm text-gray-900 flex items-center gap-1">
                <Database className="h-3 w-3 text-gray-500" />
                {tableName}
              </h3>
              <p className="text-xs text-gray-500 mt-0.5">
                {(columns as any[]).length} cols
              </p>
            </div>
            
            <div className="p-2 space-y-1 max-h-48 overflow-y-auto">
              {(columns as any[]).map((column: any, index: number) => (
                <div key={index} className="flex justify-between items-center py-1 border-b border-gray-100 last:border-0">
                  <div className="flex-1 min-w-0">
                    <div className="font-mono text-xs text-gray-800 truncate">
                      {column.column_name}
                    </div>
                    {column.column_name.toLowerCase().includes('id') && (
                      <div className="text-xs text-amber-600">PK</div>
                    )}
                  </div>
                  <div className="flex items-center gap-1 ml-2 shrink-0">
                    <span className="text-xs text-gray-600">
                      {column.data_type}
                    </span>
                    {column.is_nullable === 'YES' && (
                      <span className="text-xs text-blue-600">null</span>
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
