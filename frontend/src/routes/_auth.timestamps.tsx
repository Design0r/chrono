import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { ErrorPage } from "../components/ErrorPage";
import type { Timestamp } from "../types/response";
import { TimestampTable } from "../components/Timestamps";

type TimestampsSearchParams = {
  year?: string;
  month?: string;
};

export const Route = createFileRoute("/_auth/timestamps")({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>): TimestampsSearchParams => {
    return {
      year: search.user as string,
      month: search.event as string,
    };
  },
});

function RouteComponent() {
  const { chrono } = Route.useRouteContext();

  const params = Route.useSearch();
  const year = params.year ? Number(params.year) : undefined;
  const month = params.month ? Number(params.month) : undefined;

  const timestampQ = useQuery({
    queryKey: ["timestamps", "all"],
    queryFn: () => chrono.timestamps.getAllForUser(year, month),
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

  return (
    <div className="flex flex-col container mx-auto justify-center align-middle gap-6 p-4">
      <TimestampTable timestamps={timestamps} />
    </div>
  );
}
