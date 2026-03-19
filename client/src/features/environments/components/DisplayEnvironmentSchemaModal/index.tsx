import { CurrentDatabaseSchema } from "@/features/environments/components/CurrentEnvironmentSchema";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import type { DatabaseSchemaItem, Environment } from "@/features/types";

type Props = {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    environmentData: Environment;
    environmentSchema: DatabaseSchemaItem[];
    isLoadingEnvironmentSchema: boolean;
    onClose: () => void;
}

export function DisplayEnvironmentSchemaModal({ open, onOpenChange, environmentData, environmentSchema, isLoadingEnvironmentSchema, onClose }: Props) {
    return (
              <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="max-w-[95vw] lg:max-w-6xl h-[85vh] flex flex-col p-0 overflow-hidden border-stone-200 rounded-3xl shadow-2xl">
          <DialogHeader className="p-6 border-b bg-background backdrop-blur-md sticky top-0 z-20">
            <div className="flex items-center justify-between">
              <div>
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
              <CurrentDatabaseSchema
                schema={environmentSchema || []}
                isLoading={isLoadingEnvironmentSchema}
              />
            </div>
          </div>

          <div className="p-4 border-t bg-white/80 backdrop-blur-md flex items-center justify-end gap-6 sticky bottom-0 z-20">
            <Button onClick={onClose}>Close</Button>
          </div>
        </DialogContent>
      </Dialog>
    )
}