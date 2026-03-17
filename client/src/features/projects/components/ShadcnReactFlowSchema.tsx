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
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
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
    <Card className="min-w-[300px] border-2 shadow-lg overflow-hidden border-slate-200/50">
      <CardHeader className="bg-slate-50/50 pb-3 border-b">
        <CardTitle className="text-sm font-bold flex items-center justify-between">
          <div className="flex items-center gap-2 text-slate-900">
            <Database className="h-4 w-4 text-blue-500" />
            {data.tableName}
          </div>

        </CardTitle>
      </CardHeader>
      <CardContent className="p-0">
        <div className="divide-y divide-slate-100">
          {data.columns.map((column: any, index: number) => {
            const isPK = column.column_name === 'id'; 
            
            return (
              <div key={index} className="group flex justify-between items-center px-4 py-2 hover:bg-blue-50/30 transition-colors">
                <div className="flex items-center gap-2 overflow-hidden">
                  {isPK ? (
                    <span className="text-[10px] font-bold text-black">PK</span>
                  ) : (
                    <div className="w-1 h-1 rounded-full bg-slate-300 group-hover:bg-blue-400" />
                  )}
                  <span className={`font-mono text-xs truncate ${isPK ? 'font-bold text-slate-900' : 'text-slate-700'}`}>
                    {column.column_name}
                  </span>
                </div>
                
                <div className="flex items-center gap-2 shrink-0">
                  <span className="text-[10px] font-medium text-slate-400 uppercase tracking-tighter">
                    {column.data_type.replace('without time zone', '')}
                  </span>
                  {column.is_nullable === 'NO' && (
                    <span className="text-red-500 text-[10px] font-bold" title="Required">*</span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  )
})

TableNode.displayName = 'TableNode'

const nodeTypes = {
  table: TableNode,
}

interface ShadcnReactFlowSchemaProps {
  schema: any
  isLoading?: boolean
}

export function ShadcnReactFlowSchema({ schema, isLoading }: ShadcnReactFlowSchemaProps) {
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
  const groupedByTable: Record<string, any[]> = (schema || []).reduce((acc: Record<string, any[]>, item: any) => {
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
  const maxCols = 4
  const cardWidth = 340
  const gap = 40

Object.entries(groupedByTable).forEach(([tableName, columns]) => {
  const nodeId = `table-${tableName}`;
  
  const estimatedHeight = 60 + (columns.length * 32);

  initialNodes.push({
    id: nodeId,
    type: 'table',
    position: { x: xOffset, y: yOffset },
    data: { tableName, columns },
  });

  col += 1;
  if (col >= maxCols) {
    col = 0;
    xOffset = 50;
    yOffset += estimatedHeight + 80; 
  } else {
    xOffset += cardWidth + gap;
  }
});

  const [nodes, , onNodesChange] = useNodesState(initialNodes)
  const [edges, , onEdgesChange] = useEdgesState(initialEdges)

  return (
    <div className="h-[600px] w-full border rounded-lg">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        fitView
        attributionPosition="bottom-left"
        defaultViewport={{ x: 0, y: 0, zoom: 0.9 }}
      >
        <Background />
        <Controls 
          showZoom={true}
          showFitView={true}
          showInteractive={false}
        />
      </ReactFlow>
    </div>
  )
}
