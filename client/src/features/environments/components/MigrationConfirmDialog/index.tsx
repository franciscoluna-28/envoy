import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void;
  environmentName: string;
  environmentType?: 'production' | 'staging' | 'development' | string;
  isLoading?: boolean;
}

export function MigrationConfirmDialog({
  open,
  onOpenChange,
  onConfirm,
  environmentName,
  environmentType = "development",
  isLoading = false,
}: Props) {
  
  const isProduction = environmentType === 'production';

  return (
<AlertDialog open={open} onOpenChange={onOpenChange}>
  <AlertDialogContent className="max-w-[400px] p-0 overflow-hidden border-stone-200 shadow-2xl bg-white gap-0">    
    <div className="p-6 space-y-6">
      <AlertDialogHeader className="space-y-2 text-left">
        <AlertDialogTitle className="text-xl font-semibold text-stone-900 leading-none">
          Run Migration?
        </AlertDialogTitle>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="text-xs font-medium border">
            {environmentName}
          </Badge>
          {isProduction && (
            <Badge variant={isProduction ? "destructive" : "secondary"} className="text-xs capitalize h-5 px-2 font-medium border-none">
              {environmentType}
            </Badge>
          )}
        </div>
      </AlertDialogHeader>

      <div className="space-y-3">
        {isProduction && (
          <div className="bg-red-50/50 border-l-4 border-red-600 p-3">
            <p className="text-[11px] font-bold text-red-900 uppercase tracking-widest mb-1">
              Critical Warning
            </p>
            <p className="text-[12px] text-red-800/90 leading-snug">
              Modifying persistent production schema. Ensure backups are verified.
            </p>
          </div>
        )}

       
      </div>

      <AlertDialogFooter className="flex gap-3 sm:justify-start">
        <Button 
          variant="outline" 
          onClick={() => onOpenChange(false)} 
          disabled={isLoading}
          className="flex-1 text-xs h-10"
        >
          Cancel
        </Button>
        <Button 
          variant={isProduction ? "destructive" : "default"} 
          onClick={onConfirm} 
          disabled={isLoading}
          className={`flex-1 text-xs h-10`}
        >
          {isLoading ? (
            <span className="flex items-center gap-2">
              <span className="h-3 w-3 animate-spin rounded-full border-2 border-current border-t-transparent" />
              Running...
            </span>
          ) : 'Confirm Migration'}
        </Button>
      </AlertDialogFooter>
    </div>
  </AlertDialogContent>
</AlertDialog>
  );
}