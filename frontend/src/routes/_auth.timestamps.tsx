import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { ErrorPage } from "../components/ErrorPage";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { TeamTimestamps } from "../components/Timestamps";
import { TimestampTableByWeek } from "../components/timestamps/TimestampTableByWeek";
import {
  durationFromTimestamps,
  formatCounter,
  isoToDateLocal,
  secondsToCounter,
} from "../lib/timestamp-utils";
import type { User } from "../types/auth";
import type { Timestamp } from "../types/response";

type TimestampsSearchParams = {
  startDate?: string;
  endDate?: string;
};

export const Route = createFileRoute("/_auth/timestamps")({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>): TimestampsSearchParams => {
    return {
      startDate: search.startDate as string,
      endDate: search.endDate as string,
    };
  },
});

function RouteComponent() {
  const { auth, chrono } = Route.useRouteContext();
  const navigate = useNavigate();

  const params = Route.useSearch();

  // Die Search-Params sind die Wahrheitsquelle. Vorher spiegelte lokaler State
  // sie und navigierte per Effect zurück in die URL — eine Rückkopplung, die
  // bei jedem Render-Durchlauf eine Navigation auslösen konnte.
  const startDate = params.startDate;
  const endDate = params.endDate;

  const setRange = (range: { startDate?: string; endDate?: string }) =>
    navigate({
      to: "/timestamps",
      search: {
        startDate: range.startDate || undefined,
        endDate: range.endDate || undefined,
      },
    });

  const timestampQ = useQuery({
    queryKey: ["timestamps", startDate, endDate],
    queryFn: () => chrono.timestamps.getAllForUser(startDate, endDate),
    staleTime: 1000 * 60 * 1, // 1min
    gcTime: 1000 * 60 * 30, // 30min
    retry: false,
  });

  const userQ = useQuery({
    queryKey: ["user", auth.userId],
    queryFn: () => chrono.users.getUserById(auth.userId!),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
    retry: false,
  });

  const queries = [timestampQ, userQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const timestamps = timestampQ.data! as Timestamp[];
  const user = userQ.data! as User;

  const counter = secondsToCounter(durationFromTimestamps(timestamps));

  return (
    <div className="flex flex-col container mx-auto justify-center align-middle gap-6 p-4">
      <div className="grid gap-4 grid-cols-1 lg:flex justify-between mb-6 items-center">
        <div className="flex gap-4">
          <label className="gap-1 flex flex-col ">
            <span className="pl-0.5">Filter Start</span>
            <input
              type="date"
              className="input border-info/15"
              defaultValue={startDate && isoToDateLocal(startDate)}
              onChange={(e) => setRange({ startDate: e.target.value, endDate })}
            />
          </label>
          <label className="gap-1 flex flex-col ">
            <span className="pl-0.5">Filter End</span>
            <input
              type="date"
              className="input border-info/15"
              defaultValue={endDate && isoToDateLocal(endDate)}
              onChange={(e) => setRange({ startDate, endDate: e.target.value })}
            />
          </label>
        </div>

        <div className="flex flex-col gap-1">
          <span className="text-md text-left pl-0.5 text-base-content/90">
            Total Time
          </span>
          <h2 className="w-fit text-lg text-primary bg-base-200 px-4 py-1.5 rounded-lg border ">
            {(() => {
              const f = formatCounter(counter);
              return (
                <>
                  <span>{f.hours}</span>
                  <span>:</span>
                  <span>{f.minutes}</span> h
                </>
              );
            })()}
          </h2>
        </div>
      </div>

      <TimestampTableByWeek timestamps={timestamps} user={user} />
      <TeamTimestamps startDate={startDate} endDate={endDate} user={user} />
    </div>
  );
}
