import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { Terminal, Copy } from "lucide-react";
import type { TablePermission } from "@/features/types";
import { toast } from "sonner";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  permissions: TablePermission[];
  title: string;
  databaseUser: string;
};

export function PermissionResultsModal({
  isOpen,
  onClose,
  permissions,
  title,
  databaseUser,
}: Props) {
  const hasIssues = permissions.some((p) => p.is_missing);

  const generateGrant = (p: TablePermission) => {
    if (!p.privileges || p.privileges.length === 0)
      return `-- No specific privileges detected for ${p.table_name}`;

    const privs = p.privileges.join(", ");

    // Dare to replace this with a "GRANT ALL" and I won't accept your PR.
    return `GRANT ${privs} ON TABLE ${p.table_name} TO ${databaseUser};`;
  };

  const copyAllGrants = () => {
    const missingPermissions = permissions.filter((p) => p.is_missing);
    const allGrants = missingPermissions.map(generateGrant).join("\n");

    navigator.clipboard.writeText(allGrants);
    toast.success(`Copied ${missingPermissions.length} granular grants`);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success("Copied to clipboard");
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl! h-[85vh] p-0 flex flex-col overflow-y-auto">
        <DialogHeader className="p-6 border-b shrink-0">
          <div className="flex items-center justify-between">
            <DialogTitle className="flex items-center gap-2 text-lg font-semibold">
              {title}
            </DialogTitle>
          </div>
          {hasIssues && (
            <Button
              variant="outline"
              size="sm"
              onClick={copyAllGrants}
              className="gap-2 border-destructive text-destructive hover:bg-destructive hover:text-white"
            >
              <Copy className="h-4 w-4" />
              Copy All CRUD Required Grants
            </Button>
          )}
        </DialogHeader>

        <ScrollArea className="flex-1">
          <div className="p-6 space-y-4">
            {permissions.map((permission, index) => (
              <div
                key={index}
                className="border rounded-lg p-4 transition-colors hover:bg-accent/50"
              >
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <Terminal className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm font-semibold">
                      {permission.table_name}
                    </span>
                  </div>
                  {permission.is_missing && (
                    <span className="text-[10px] font-bold text-destructive uppercase">
                      Missing Access
                    </span>
                  )}
                </div>

                <div className="flex flex-wrap gap-2">
                  {permission.privileges.map((priv) => (
                    <Badge
                      key={priv}
                      variant={permission.is_missing ? "outline" : "secondary"}
                      className="text-[10px]"
                    >
                      {priv}
                    </Badge>
                  ))}
                </div>

                {permission.is_missing && (
                  <div className="mt-4 pt-4 border-t border-dashed">
                    <div className="flex items-center justify-between bg-muted p-2 rounded-md">
                      <code className="text-xs text-muted-foreground truncate mr-2">
                        {generateGrant(permission)}
                      </code>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 shrink-0"
                        onClick={() =>
                          copyToClipboard(generateGrant(permission))
                        }
                      >
                        <Copy className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </ScrollArea>

        <div className="p-4 border-t bg-muted/50 flex justify-between items-center shrink-0">
          <p className="text-xs text-muted-foreground">
            Envoy Security Audit v1.0
          </p>
          <Button onClick={onClose} size="sm">
            Close Inspector
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
