import type { QueryClient } from "@tanstack/react-query";
import { aworkQueries, eventQueries, timestampQueries, userQueries } from "../api/queries";

export function invalidateEventData(queryClient: QueryClient) {
  return Promise.all([
    queryClient.invalidateQueries({ queryKey: eventQueries.all() }),
    queryClient.invalidateQueries({ queryKey: userQueries.all() }),
    queryClient.invalidateQueries({ queryKey: timestampQueries.workHours() }),
    queryClient.invalidateQueries({ queryKey: aworkQueries.workTimes() }),
  ]);
}
