import CodeMirror from "@uiw/react-codemirror";
import { sql as sqlLanguage } from "@codemirror/lang-sql";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Terminal } from "lucide-react";

type PreviewSQLModalProps = {
  sql: string;
  name: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function PreviewSQLModal({ sql, name, open, onOpenChange }: PreviewSQLModalProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl! p-0! overflow-hidden gap-0">
        <DialogHeader className="p-4 border-b">
          <div className="flex items-center">
            <DialogTitle className="text-sm">
              {name || "migration_preview.sql"}
            </DialogTitle>
          </div>
        </DialogHeader>
        
        <ScrollArea className="max-h-[80vh] overflow-x-auto">
          <CodeMirror
            value={sql}
            height="auto"
            extensions={[sqlLanguage()]}
            editable={false}
            basicSetup={{
              lineNumbers: true,
              foldGutter: true,
              highlightActiveLine: false,
            }}
            className="text-[11px]"
          />
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
}