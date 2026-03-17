interface MinimalSchemaProps {
  schema: any
  isLoading?: boolean
}

export function MinimalSchema({ schema, isLoading }: MinimalSchemaProps) {
  if (isLoading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500">Loading...</div>
      </div>
    )
  }

  if (!schema || schema.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500">No schema data</div>
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
        <div key={tableName} className="border border-black">
          <div className="border-b border-black px-4 py-2 bg-black text-white">
            <div className="font-bold">{tableName}</div>
            <div className="text-xs opacity-75">{(columns as any[]).length} columns</div>
          </div>
          
          <div className="divide-y divide-black">
            {(columns as any[]).map((column: any, index: number) => (
              <div key={index} className="flex justify-between px-4 py-2">
                <div className="font-mono text-sm">{column.column_name}</div>
                <div className="text-sm text-gray-600">
                  {column.data_type}
                  {column.is_nullable === 'YES' && (
                    <span className="ml-2 text-xs">(null)</span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
