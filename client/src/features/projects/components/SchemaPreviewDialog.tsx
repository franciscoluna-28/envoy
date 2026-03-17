import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { ShadcnReactFlowSchema } from "./ShadcnReactFlowSchema"
import { useEnvironmentSchema } from "../hooks/useEnvironmentSchema"

interface SchemaPreviewDialogProps {
  environmentId: string
  environmentName: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SchemaPreviewDialog({
  environmentId,
  environmentName,
  open,
  onOpenChange,
}: SchemaPreviewDialogProps) {
  const { data: schema, isLoading } = useEnvironmentSchema(environmentId)

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="w-full !max-w-[90vw] h-[85vh] max-h-[700px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl">
            Database Schema: {environmentName}
          </DialogTitle>
        </DialogHeader>
        <div className="flex-1 min-h-0">
          <ShadcnReactFlowSchema schema={schema} isLoading={isLoading} />
        </div>
       
      </DialogContent>
    </Dialog>
  )
}
