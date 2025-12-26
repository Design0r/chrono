import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { ErrorPage } from "../components/ErrorPage";
import type { Timestamp } from "../types/response";
import {
  durationFromTimestamps,
  isoToDateLocal,
  secondsToCounter,
  TimestampTable,
} from "../components/Timestamps";
import { useEffect, useState } from "react";

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
  const { chrono } = Route.useRouteContext();
  const navigate = useNavigate();

  const params = Route.useSearch();

  const [startDate, setStartDate] = useState<string | undefined>();
  const [endDate, setEndDate] = useState<string | undefined>();

  useEffect(() => {
    if (params.startDate) setStartDate(params.startDate);
    else setStartDate(undefined);

    if (params.endDate) setEndDate(params.endDate);
    else setEndDate(undefined);
  }, [params]);

  useEffect(() => {
    navigate({
      to: "/timestamps",
      search: {
        startDate: startDate,
        endDate: endDate,
      },
    });
  }, [startDate, endDate]);

  const timestampQ = useQuery({
    queryKey: ["timestamps", startDate, endDate],
    queryFn: () => chrono.timestamps.getAllForUser(startDate, endDate),
    staleTime: 1000 * 60 * 1, // 1min
    gcTime: 1000 * 60 * 30, // 30min
    retry: false,
  });

  const queries = [timestampQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const timestamps = timestampQ.data! as Timestamp[];

  const counter = secondsToCounter(durationFromTimestamps(timestamps));

  return (
    <div className="flex flex-col container mx-auto justify-center align-middle gap-6 p-4">
      <div className="grid gap-4 grid-cols-1 md:grid-cols-2">
        <label className="space-x-2 flex flex-col lg:flex-row items-center">
          <span>Start Date</span>
          <input
            type="date"
            className="input"
            defaultValue={startDate && isoToDateLocal(startDate)}
            onChange={(e) => setStartDate(e.target.value)}
          />
        </label>

        <label className="space-x-2 flex flex-col lg:flex-row items-center">
          <span>End Date</span>
          <input
            type="date"
            className="input"
            defaultValue={endDate && isoToDateLocal(endDate)}
            onChange={(e) => setEndDate(e.target.value)}
          />
        </label>
      </div>

      <h2 className="text-lg text-center xl:text-left">
        Total Duration: {counter.hours}h {counter.minutes}m {counter.seconds}s
      </h2>

      <TimestampTable timestamps={timestamps} />
    </div>
  );
}
