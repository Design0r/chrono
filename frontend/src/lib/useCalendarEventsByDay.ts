import { useQueries } from "@tanstack/react-query";
import { useMemo } from "react";
import { eventQueries } from "../api/queries";
import type { Month, Timestamp } from "../types/response";
import { dateKeyFromDate } from "./timestamp-utils";

export type EventsByDay = Map<string, string[]>;

export function useCalendarEventsByDay(timestamps: Timestamp[], userId: number): EventsByDay {
  const months = useMemo(() => {
    const unique = new Map<string, [number, number]>();
    for (const t of timestamps) {
      const start = new Date(t.start_time);
      const year = start.getFullYear();
      const month = start.getMonth() + 1;
      unique.set(`${year}-${month}`, [year, month]);
    }
    return [...unique.values()];
  }, [timestamps]);

  return useQueries({
    queries: months.map(([year, month]) => eventQueries.month(year, month)),
    combine: (results) =>
      eventsByDayForUser(
        results.map((r) => r.data),
        userId,
      ),
  });
}

function eventsByDayForUser(months: (Month | undefined)[], userId: number): EventsByDay {
  const map: EventsByDay = new Map();
  for (const month of months) {
    for (const day of month?.days ?? []) {
      const names = new Set(
        (day.events ?? []).filter((e) => e.event.user_id === userId).map((e) => e.event.name),
      );
      if (names.size) map.set(dateKeyFromDate(new Date(day.date)), [...names]);
    }
  }
  return map;
}
