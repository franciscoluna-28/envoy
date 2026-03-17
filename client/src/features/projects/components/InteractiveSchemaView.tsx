import React, { useCallback } from 'react'
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  type Connection,
  type Edge,
  type Node,
  Panel,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Database, Key, Hash, Type } from 'lucide-react'

interface ColumnNodeProps {
  data: {
    columnName: string
    dataType: string
    isNullable: string
    tableName: string
  }
}

const ColumnNode: React.FC<ColumnNodeProps> = ({ data }) => {
  const isPrimaryKey = data.columnName.toLowerCase().includes('id') || data.columnName.toLowerCase().endsWith('_id')
  const isNullable = data.isNullable === 'YES'

  return (
    <Card className="min-w-[200px] max-w-[250px] shadow-md">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-mono truncate">{data.columnName}</CardTitle>
          <div className="flex items-center gap-1">
            {isPrimaryKey && <Key className="h-3 w-3 text-yellow-500" />}
            {isNullable && <Hash className="h-3 w-3 text-blue-500" />}
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="flex items-center justify-between">
          <Badge variant="outline" className="text-xs">
            <Type className="h-3 w-3 mr-1" />
            {data.dataType}
          </Badge>
          <Badge 
            variant={isNullable ? 'secondary' : 'default'} 
            className="text-xs"
          >
            {isNullable ? 'NULL' : 'NOT NULL'}
          </Badge>
        </div>
      </CardContent>
    </Card>
  )
}

const nodeTypes = {
  column: ColumnNode,
}

interface InteractiveSchemaViewProps {
  schema: any
  isLoading?: boolean
}

export function InteractiveSchemaView({ schema, isLoading }: InteractiveSchemaViewProps) {
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

  // Create nodes for each column
  const initialNodes: Node[] = []
  const initialEdges: Edge[] = []

  let yOffset = 50
  let xOffset = 50

  Object.entries(groupedByTable).forEach(([tableName, columns]) => {
  const columnsArray = columns as any[]
  
  // Add table header node
  initialNodes.push({
    id: `table-${tableName}`,
    type: 'default',
    position: { x: xOffset, y: yOffset },
    data: { 
      label: (
        <div className="px-4 py-2 bg-primary text-primary-foreground rounded-md font-semibold">
          <Database className="h-4 w-4 inline mr-2" />
          {tableName}
        </div>
      )
    },
    style: {
      background: 'transparent',
      border: 'none',
    },
  })

  yOffset += 80

  // Add column nodes
  columnsArray.forEach((column: any, colIndex: number) => {
    const nodeId = `${tableName}-${column.column_name}`
    initialNodes.push({
      id: nodeId,
      type: 'column',
      position: { x: xOffset + (colIndex % 3) * 260, y: yOffset + Math.floor(colIndex / 3) * 120 },
      data: {
        columnName: column.column_name,
        dataType: column.data_type,
        isNullable: column.is_nullable,
        tableName: tableName,
      },
    })

    // Add edge from table to column
    initialEdges.push({
      id: `edge-${tableName}-${column.column_name}`,
      source: `table-${tableName}`,
      target: nodeId,
      type: 'smoothstep',
      style: { stroke: '#94a3b8', strokeWidth: 1 },
    })
  })

  yOffset += Math.ceil(columnsArray.length / 3) * 120 + 150
  if (xOffset > 800) {
    xOffset = 50
    yOffset = 50
  } else {
    xOffset += 800
  }
})

  const [nodes, , onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)

  const onConnect = useCallback(
    (params: Edge | Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  )

  return (
    <div className="h-[600px] w-full border rounded-lg">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        fitView
        attributionPosition="bottom-left"
      >
        <Background color="#e2e8f0" gap={16} />
        <Controls />
        <MiniMap 
          nodeColor={(node) => {
            if (node.id.startsWith('table-')) return '#3b82f6'
            return '#64748b'
          }}
          className="bg-background border"
        />
        <Panel position="top-right" className="bg-background border rounded-md p-2">
          <div className="text-xs text-muted-foreground space-y-1">
            <div className="flex items-center gap-2">
              <Key className="h-3 w-3 text-yellow-500" />
              <span>Primary Key</span>
            </div>
            <div className="flex items-center gap-2">
              <Hash className="h-3 w-3 text-blue-500" />
              <span>Nullable</span>
            </div>
          </div>
        </Panel>
      </ReactFlow>
    </div>
  )
}
