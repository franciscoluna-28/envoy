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
import {
  DatabaseSchemaNode,
  DatabaseSchemaNodeHeader,
  DatabaseSchemaNodeBody,
  DatabaseSchemaTableRow,
  DatabaseSchemaTableCell,
} from '@/components/database-schema-node'
import { Database } from 'lucide-react'

interface DatabaseSchemaNodeData {
  data: {
    label: string
    schema: { title: string; type: string; nullable: string }[]
  }
}

const TableNode = memo<DatabaseSchemaNodeData>(({ data }) => {
  return (
    <DatabaseSchemaNode className="min-w-[400px] max-w-[500px]">
      <DatabaseSchemaNodeHeader>{data.label}</DatabaseSchemaNodeHeader>
      <DatabaseSchemaNodeBody>
        {data.schema.map((entry) => (
          <DatabaseSchemaTableRow key={entry.title}>
            <DatabaseSchemaTableCell className="font-mono text-sm py-3">
              {entry.title}
            </DatabaseSchemaTableCell>
            <DatabaseSchemaTableCell className="text-sm py-3 text-right">
              <div className="flex items-center justify-end gap-2">
                <span className="text-gray-600">{entry.type}</span>
                {entry.nullable === 'YES' && (
                  <span className="text-xs text-blue-600 bg-blue-50 px-2 py-1 rounded">
                    nullable
                  </span>
                )}
              </div>
            </DatabaseSchemaTableCell>
          </DatabaseSchemaTableRow>
        ))}
      </DatabaseSchemaNodeBody>
    </DatabaseSchemaNode>
  )
})

TableNode.displayName = 'TableNode'

const nodeTypes = {
  table: TableNode,
}

interface SimpleReactFlowSchemaProps {
  schema: any
  isLoading?: boolean
}

export function SimpleReactFlowSchema({ schema, isLoading }: SimpleReactFlowSchemaProps) {
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

  // Create nodes for each table
  const initialNodes: Node[] = []
  const initialEdges: Edge[] = []

  let yOffset = 50
  let xOffset = 50

  Object.entries(groupedByTable).forEach(([tableName, columns]) => {
    const nodeId = `table-${tableName}`
    const schemaData = (columns as any[]).map(col => ({
      title: col.column_name,
      type: col.data_type,
      nullable: col.is_nullable
    }))
    
    initialNodes.push({
      id: nodeId,
      type: 'table',
      position: { x: xOffset, y: yOffset },
      data: {
        label: tableName,
        schema: schemaData,
      },
    })

    yOffset += 300 // More vertical space between tables
    if (yOffset > 600) {
      yOffset = 50
      xOffset += 520 // More horizontal space between tables
    }
  })

  const [nodes, , onNodesChange] = useNodesState(initialNodes)
  const [edges, , onEdgesChange] = useEdgesState(initialEdges)

  return (
    <div className="h-[600px] w-full border rounded-lg bg-background">
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
        />
      </ReactFlow>
    </div>
  )
}
