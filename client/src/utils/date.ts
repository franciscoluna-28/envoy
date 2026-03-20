import { formatDistanceToNow } from "date-fns";

export function formatDate(dateString: string) {
  return formatDistanceToNow(new Date(dateString), {
    addSuffix: true,
  });
}
