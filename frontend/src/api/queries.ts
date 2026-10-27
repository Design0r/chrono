import { QueryClient, queryOptions } from "@tanstack/react-query";
import type { UserWithVacation } from "../types/auth";
import { chrono } from "./chrono/client";

const SECOND = 1000;
const MINUTE = 60 * SECOND;
const HOUR = 60 * MINUTE;

export function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        staleTime: MINUTE,
        gcTime: 30 * MINUTE,
      },
    },
  });
}

export const userQueries = {
  all: () => ["users"] as const,
  lists: () => [...userQueries.all(), "list"] as const,
  list: ({ includeInactive = false }: { includeInactive?: boolean } = {}) =>
    queryOptions({
      queryKey: [...userQueries.lists(), { includeInactive }] as const,
      queryFn: () => chrono.users.getUsers({ includeInactiveUsers: includeInactive }),
      staleTime: 30 * MINUTE,
    }),
  listWithVacation: (
    year: number,
    { includeInactive = false }: { includeInactive?: boolean } = {},
  ) =>
    queryOptions({
      queryKey: [...userQueries.lists(), { includeInactive, vacationYear: year }] as const,
      queryFn: async () =>
        (await chrono.users.getUsers({
          vacation: { year },
          includeInactiveUsers: includeInactive,
        })) as UserWithVacation[],
      staleTime: 30 * MINUTE,
    }),
  detail: (id: number) =>
    queryOptions({
      queryKey: [...userQueries.all(), "detail", id] as const,
      queryFn: () => chrono.users.getUserById(id),
      staleTime: 6 * HOUR,
    }),
  detailWithVacation: (id: number, year: number) =>
    queryOptions({
      queryKey: [...userQueries.detail(id).queryKey, { vacationYear: year }] as const,
      queryFn: async () => (await chrono.users.getUserById(id, { year })) as UserWithVacation,
      staleTime: 6 * HOUR,
    }),
};

export const eventQueries = {
  all: () => ["events"] as const,
  month: (year: number, month: number) =>
    queryOptions({
      queryKey: [...eventQueries.all(), "month", year, month] as const,
      queryFn: () => chrono.events.getEventsForMonth(year, month),
    }),
  vacationGraph: (year: number) =>
    queryOptions({
      queryKey: [...eventQueries.all(), "vacationGraph", year] as const,
      queryFn: () => chrono.events.getVacationGraph(year),
    }),
};

export const timestampQueries = {
  all: () => ["timestamps"] as const,
  today: () =>
    queryOptions({
      queryKey: [...timestampQueries.all(), "today"] as const,
      queryFn: () => chrono.timestamps.getForToday(),
      staleTime: 10 * MINUTE,
    }),
  latest: () =>
    queryOptions({
      queryKey: [...timestampQueries.all(), "latest"] as const,
      queryFn: () => chrono.timestamps.getLatest(),
      staleTime: 30 * SECOND,
    }),
  forCurrentUser: (range: { startDate?: string; endDate?: string }) =>
    queryOptions({
      queryKey: [...timestampQueries.all(), "currentUser", range] as const,
      queryFn: () => chrono.timestamps.getAllForUser(range.startDate, range.endDate),
    }),
  forAllUsers: (range: { startDate?: string; endDate?: string }) =>
    queryOptions({
      queryKey: [...timestampQueries.all(), "allUsers", range] as const,
      queryFn: () => chrono.timestamps.getAll(range.startDate, range.endDate),
    }),
  workHours: () => [...timestampQueries.all(), "workHours"] as const,
  workHoursForCurrentUser: (year: number) =>
    queryOptions({
      queryKey: [...timestampQueries.workHours(), "currentUser", year] as const,
      queryFn: () => chrono.timestamps.getWorkHours(year),
      staleTime: HOUR,
    }),
  workHoursForAllUsers: (year: number) =>
    queryOptions({
      queryKey: [...timestampQueries.workHours(), "allUsers", year] as const,
      queryFn: () => chrono.timestamps.getWorkHoursForAllUsers(year),
    }),
};

export const aworkQueries = {
  all: () => ["awork"] as const,
  workTimes: () => [...aworkQueries.all(), "workTimes"] as const,
  workTimesForYear: (year: number) =>
    queryOptions({
      queryKey: [...aworkQueries.workTimes(), year] as const,
      queryFn: () => chrono.awork.getWorkTimesforYear(year),
      staleTime: HOUR,
    }),
  users: () =>
    queryOptions({
      queryKey: [...aworkQueries.all(), "users"] as const,
      queryFn: () => chrono.awork.getUsers(),
      staleTime: 6 * HOUR,
    }),
};

export const notificationQueries = {
  list: () =>
    queryOptions({
      queryKey: ["notifications"] as const,
      queryFn: () => chrono.notifications.get(),
      staleTime: 10 * MINUTE,
    }),
};

export const requestQueries = {
  list: () =>
    queryOptions({
      queryKey: ["requests"] as const,
      queryFn: () => chrono.requests.getRequests(),
      staleTime: 0,
    }),
};

export const settingsQueries = {
  detail: () =>
    queryOptions({
      queryKey: ["settings"] as const,
      queryFn: () => chrono.settings.getSettings(),
      staleTime: 30 * SECOND,
    }),
};

export function currentUserId(userId: number | null): number {
  if (userId === null) throw new Error("Not logged in");
  return userId;
}
