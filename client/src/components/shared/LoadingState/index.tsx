import { Spinner } from "@/components/ui/spinner";

/** Placeholder used in MVP to show loading state. We're not creating loading skeletons yet. */
export function LoadingState() {
  return (
    <div className="flex flex-col items-center justify-center h-[60vh] gap-4">
      <Spinner className="w-8 h-8 text-blue-600" />
      <p className="text-xs text-stone-400 animate-pulse">
        Loading infrastructure resources...
      </p>
    </div>
  );
}