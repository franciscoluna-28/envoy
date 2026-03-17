import { memo } from 'react'
import {
  ReactFlow,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  type Node,
  type Edge,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { Badge } from '@/components/ui/badge'
import { Database } from 'lucide-react'

interface TableNodeProps {
  data: {
    tableName: string
    columns: Array<{
      column_name: string
      data_type: string
      is_nullable: string
    }>
  }
}

const TableNode = memo<TableNodeProps>(({ data }) => {
  return (
    <div className="border rounded-lg bg-white shadow-sm min-w-[320px] max-w-[380px]">
      <div className="bg-linear-to-r from-blue-50 to-indigo-50 px-4 py-3 border-b">
        <h3 className="font-semibold text-gray-900 flex items-center gap-2">
          <Database className="h-4 w-4 text-blue-600" />
          {data.tableName}
        </h3>
        <p className="text-xs text-gray-500 mt-1">
          {data.columns.length} columns
        </p>
      </div>
      
      <div className="p-4 space-y-3">
        {data.columns.map((column: any, index: number) => (
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
  )
})

TableNode.displayName = 'TableNode'

const nodeTypes = {
  table: TableNode,
}

interface ReactFlowGridSchemaProps {
  schema: any
  isLoading?: boolean
}

export function ReactFlowGridSchema({ schema, isLoading }: ReactFlowGridSchemaProps) {
  if (isLoading) {
    return (
      <div className="h-[600px] flex items-center justify-center">
        <div className="text-center">
          <Database className="h-12 w-12 text-muted-foreground animate-pulse mx-auto mb-4" />
          <p className="text-muted-foreground">Loading database schema...</p>
        </div>
      </div>
    )
  }

  if (!schema || schema.length === 0) {
    return (
      <div className="h-[600px] flex items-center justify-center">
        <div className="text-center">
          <Database className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <p className="text-muted-foreground">No schema data available</p>
        </div>
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

  // Create nodes in a grid layout
  const initialNodes: Node[] = []
  const initialEdges: Edge[] = []

  let yOffset = 50
  let xOffset = 50
  let col = 0
  const maxCols = 3 // Max columns per row
  const cardWidth = 400
  const cardHeight = 300
  const gap = 50

  Object.entries(groupedByTable).forEach(([tableName, columns]) => {
    const nodeId = `table-${tableName}`
    
    initialNodes.push({
      id: nodeId,
      type: 'table',
      position: { x: xOffset, y: yOffset },
      data: {
        tableName: tableName,
        columns: columns,
      },
    })

    // Grid positioning
    col += 1
    if (col >= maxCols) {
      col = 0
      xOffset = 50
      yOffset += cardHeight + gap
    } else {
      xOffset += cardWidth + gap
    }
  })

  const [nodes, , onNodesChange] = useNodesState(initialNodes)
  const [edges, , onEdgesChange] = useEdgesState(initialEdges)

  return (
    <div className="h-[600px] w-full border rounded-lg bg-gray-50">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        fitView
        attributionPosition="bottom-left"
        defaultViewport={{ x: 0, y: 0, zoom: 0.8 }}
      >
        <Background color="#e2e8f0" gap={16} />
        <Controls 
          showZoom={true}
          showFitView={true}
          showInteractive={false}
          position="top-right"
        />
      </ReactFlow>
    </div>
  )
}
