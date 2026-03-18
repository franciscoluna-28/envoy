import { useState } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'
import { 
  ChevronRight, Database, CheckCircle, XCircle, 
  AlertCircle, Server, Activity, Clock
} from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'

import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { Badge } from '@/components/ui/badge'
import { 
  Table, TableBody, TableCell, TableHead, 
  TableHeader, TableRow 
} from '@/components/ui/table'
import { 
  Dialog, DialogContent, DialogHeader, DialogTitle 
} from '@/components/ui/dialog'

import { ShadcnReactFlowSchema, useEnvironmentSchema, useGetProject } from '@/features/projects'
import { useGetEnvironment } from '@/features/projects/hooks/useEnvironments'
import { useGetEnvironmentMigrations } from '@/features/projects/hooks/useMigrations'
import { toast } from 'sonner'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'


export const Route = createFileRoute(
  '/app/projects/$projectId/environments/$envId/',
)({
  component: RouteComponent,
})

function StatusBadge({ status }: { status?: string }) {
  const configs = {
    completed: { icon: CheckCircle, color: 'text-green-600', bg: 'bg-green-50', label: 'Success' },
    failed: { icon: XCircle, color: 'text-red-600', bg: 'bg-red-50', label: 'Failed' },
    running: { icon: Activity, color: 'text-blue-600', bg: 'bg-blue-50', label: 'Running' },
  }

  const config = configs[status as keyof typeof configs] || { 
    icon: AlertCircle, color: 'text-stone-400', bg: 'bg-stone-50', label: 'Unknown' 
  }

  const Icon = config.icon

  return (
    <div className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full ${config.bg} border border-transparent transition-colors`}>
      <Icon className={`h-3 w-3 ${config.color}`} />
      <span className={`text-[11px] font-bold uppercase tracking-wide ${config.color}`}>
        {config.label}
      </span>
    </div>
  )
}

function RouteComponent() {
  const params = Route.useParams()
  const { projectId, envId } = params
  const [isModalOpen, setIsModalOpen] = useState(false)

  const { data: projectData, isPending: isLoadingProject } = useGetProject(projectId)
  const { data: environmentData, isPending: isLoadingEnvironment } = useGetEnvironment(projectId, envId)
  const { data: environmentSchema, isPending: isLoadingEnvironmentSchema } = useEnvironmentSchema(envId)
  const { data: migrations, isPending: isLoadingMigrations } = useGetEnvironmentMigrations(envId)

  if (isLoadingProject || isLoadingEnvironment) {
    return (
      <div className="flex h-[60vh] flex-col items-center justify-center gap-4">
        <Spinner className="h-8 w-8 text-blue-600" />
        <p className="text-xs text-stone-400 font-medium animate-pulse uppercase tracking-widest">
          Polling Node Infrastructure...
        </p>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full duration-500">
      <header className="flex flex-col gap-4 py-6 border-b bg-stone-50/30">
        <div className="flex items-end justify-between">
          <div className="space-y-1">
            <h1 className="text-xl font-bold tracking-tighter text-stone-900 flex items-center gap-3 uppercase">
              <div className="p-2 bg-white border border-stone-200 rounded-xl shadow-sm">
                <Database className="h-5 w-5 text-blue-600" />
              </div>
              {environmentData?.name}
            </h1>
            <div className="flex items-center gap-2 pl-12">
               <button 
                onClick={() => {
                  navigator.clipboard.writeText(envId)
                  toast.success('Node ID copied to clipboard')
                }}
                className="text-[10px] font-bold text-stone-400 hover:text-blue-600 transition-colors uppercase tracking-tight"
              >
                Node ID: #{envId.slice(0, 8)}
              </button>
            </div>
          </div>

          <div className="flex gap-3">
            <Button 
              variant="outline" 
              size="lg" 
              onClick={() => setIsModalOpen(true)} 
             
            >
              Preview Schema
            </Button>
            <Link to="/app/projects/$projectId/environments/$envId/migrations/new" params={{ projectId, envId }}>
              <Button size="lg">
                Create Migration
              </Button>
            </Link>
          </div>
        </div>
      </header>

      <main className="flex-1 py-8 space-y-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            { label: 'Network Status', value: 'Connected', icon: <div className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse" /> },
            { label: 'Node Type', value: environmentData?.type || 'Development', icon: <Server className="w-3 h-3 text-stone-400" /> },
            { label: 'Database Engine', value: 'PostgreSQL 16', icon: <Database className="w-3 h-3 text-stone-400" /> },
          ].map((stat) => (
            <div key={stat.label} className="border p-4 rounded-2xl bg-card flex flex-col gap-1.5">
              <span className="text-[10px] font-bold text-stone-400 uppercase">{stat.label}</span>
              <div className="flex items-center gap-2">
                {stat.icon}
                <span className="text-sm font-medium text-stone-900 capitalize tracking-tight">{stat.value}</span>
              </div>
            </div>
          ))}
        </div>
        <div className="space-y-6">
          <div className="flex items-center justify-between border-b border-stone-100 pb-4">
             <div className="flex items-center gap-2">
                <Clock className="h-4 w-4 text-stone-400" />
                <h3 className="text-[11px] font-bold text-stone-500 uppercase tracking-[0.2em]">
                  Migration History
                </h3>
             </div>
             <Badge variant="outline" className="text-[10px] border-stone-200 text-stone-500 font-bold uppercase">
                {migrations?.length || 0} Total
             </Badge>
          </div>
          
          <div className="rounded-2xl border border-stone-100 overflow-hidden shadow-sm bg-white">
            <Table>
              <TableHeader className="bg-stone-50/50">
                <TableRow className="hover:bg-transparent border-stone-100">
                  <TableHead className="text-[10px] font-bold uppercase tracking-widest text-stone-400 py-4 pl-6">Node Metadata</TableHead>
                  <TableHead className="text-[10px] font-bold uppercase tracking-widest text-stone-400 py-4">Status</TableHead>
                  <TableHead className="text-[10px] font-bold uppercase tracking-widest text-stone-400 py-4 text-right pr-8">Performance / Time</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoadingMigrations ? (
                  <TableRow><TableCell colSpan={3} className="py-12 text-center"><Spinner className="mx-auto h-4 w-4 text-stone-300" /></TableCell></TableRow>
                ) : migrations?.length === 0 ? (
                  <TableRow><TableCell colSpan={3} className="py-12 text-center text-[11px] text-stone-400 font-medium uppercase tracking-widest">Initial node state: No migrations found</TableCell></TableRow>
                ) : (
                  migrations?.map((m) => (
                    <>
                     <TableRow 
  key={m.id} 
  className="border-stone-50 transition-colors group hover:bg-stone-50/30"
>
  <TableCell className="py-5 pl-6">
    <div className="flex items-center gap-2">
      <div className={`font-bold text-sm tracking-tight ${m.status === 'failed' ? 'text-red-900' : 'text-stone-900'}`}>
        {m.name}
      </div>
      {m.status === 'failed' && (
        <Popover>
          <PopoverTrigger asChild>
            <Button 
              variant="ghost" 
              className="h-6 px-2 text-[10px] font-black uppercase tracking-widest text-red-500 bg-red-50/50 hover:bg-red-100 hover:text-red-700 transition-all rounded-md border border-red-100/50 ml-1"
            >
              View Logs
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-[450px] p-0 overflow-hidden border-stone-800 shadow-2xl bg-stone-950">
            <div className="bg-stone-900 px-3 py-2 border-b border-stone-800 flex items-center justify-between">
              <span className="text-[10px] font-bold text-red-400 uppercase tracking-widest">Engine Exception</span>
              <code className="text-[9px] text-stone-500 font-mono">ID: {m?.id?.slice(0,8)}</code>
            </div>
            <div className="p-4">
               <code className="text-[11px] font-mono text-red-400 leading-relaxed block whitespace-pre-wrap">
                {m.error_message}
              </code>
            </div>
          </PopoverContent>
        </Popover>
      )}
    </div>
    <div className="text-[10px] text-stone-400 uppercase font-bold tracking-widest mt-0.5">
      {m.description || 'System Migration'}
    </div>
  </TableCell>
  
  <TableCell>
    <StatusBadge status={m.status} />
  </TableCell>
  
  <TableCell className="text-right pr-8">
    <div className="text-xs font-bold text-stone-900 tabular-nums">
      {m.duration ? `${m.duration}ms` : '--'}
    </div>
    <div className="text-[10px] text-stone-400 font-bold uppercase tracking-tighter mt-0.5 whitespace-nowrap">
      {m.executed_at ? formatDistanceToNow(new Date(m.executed_at), { addSuffix: true }) : 'Scheduled'}
    </div>
  </TableCell>
</TableRow>
                    </>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </div>
      </main>
<Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
  <DialogContent className="max-w-[95vw] lg:max-w-6xl h-[85vh] flex flex-col p-0 overflow-hidden border-stone-200 rounded-3xl shadow-2xl">
    <DialogHeader className="p-6 border-b bg-white/80 backdrop-blur-md sticky top-0 z-20">
      <div className="flex items-center justify-between">
        <div className="">
          <DialogTitle className="text-xl font-bold tracking-tight text-stone-900 flex items-center gap-2">
            Infrastructure Schema Preview
          </DialogTitle>
          <p className="text-[10px] text-stone-400 font-bold uppercase">
            Live Database Topology • {environmentData?.name}
          </p>
        </div>
      </div>
    </DialogHeader>

    <div className="flex-1 overflow-y-auto p-0 relative group">
      <div className="absolute inset-0">
        <ShadcnReactFlowSchema 
          schema={environmentSchema || []} 
          isLoading={isLoadingEnvironmentSchema} 
        />
      </div>
    </div>


   <div className="p-4 border-t bg-white/80 backdrop-blur-md flex items-center justify-end gap-6 sticky bottom-0 z-20">
    
      <Button onClick={() => setIsModalOpen(false)}>
        Close
      </Button>
    </div>
  </DialogContent>
</Dialog>
    </div>
  )
}