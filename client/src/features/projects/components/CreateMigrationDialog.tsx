import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { useState } from "react"

interface CreateMigrationDialogProps {
  environmentId: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateMigrationDialog({
  environmentId,
  open,
  onOpenChange,
}: CreateMigrationDialogProps) {
  const [migrationName, setMigrationName] = useState("")
  const [migrationDescription, setMigrationDescription] = useState("")
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!migrationName.trim()) return

    setIsLoading(true)
    try {
      // TODO: Implement migration creation API call
      console.log("Creating migration:", {
        environmentId,
        name: migrationName,
        description: migrationDescription,
      })
      
      // Reset form
      setMigrationName("")
      setMigrationDescription("")
      onOpenChange(false)
    } catch (error) {
      console.error("Failed to create migration:", error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Create Migration</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <Label htmlFor="migration-name">Migration Name</Label>
            <Input
              id="migration-name"
              value={migrationName}
              onChange={(e) => setMigrationName(e.target.value)}
              placeholder="e.g., add_user_table"
              required
            />
          </div>
          <div>
            <Label htmlFor="migration-description">Description (optional)</Label>
            <Textarea
              id="migration-description"
              value={migrationDescription}
              onChange={(e) => setMigrationDescription(e.target.value)}
              placeholder="Describe what this migration does..."
              rows={3}
            />
          </div>
          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={!migrationName.trim() || isLoading}>
              {isLoading ? "Creating..." : "Create Migration"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
