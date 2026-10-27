import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { currentUserId, timestampQueries, userQueries } from "../api/queries";
import { TeamTimestamps } from "../components/Timestamps";
import { TimestampTableByWeek } from "../components/timestamps/TimestampTableByWeek";
import {
  durationFromTimestamps,
  formatCounter,
  isoToDateLocal,
  secondsToCounter,
  yearOfRange,
} from "../lib/timestamp-utils";

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
  loaderDeps: ({ search: { startDate, endDate } }) => ({ startDate, endDate }),
  loader: async ({ context: { auth, queryClient }, deps }) => {
    const userId = currentUserId(auth.userId);
    const [user] = await Promise.all([
      queryClient.ensureQueryData(userQueries.detail(userId)),
      queryClient.ensureQueryData(timestampQueries.forCurrentUser(deps)),
    ]);
    if (user.is_superuser) {
      void queryClient.prefetchQuery(timestampQueries.workedSecondsPerUser(deps));
      void queryClient.prefetchQuery(userQueries.list());
      void queryClient.prefetchQuery(
        timestampQueries.workHoursForAllUsers(yearOfRange(deps.startDate)),
      );
    }
    return { userId };
  },
});

function RouteComponent() {
  const { userId } = Route.useLoaderData();
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

  const { data: timestamps } = useSuspenseQuery(
    timestampQueries.forCurrentUser({ startDate, endDate }),
  );
  const { data: user } = useSuspenseQuery(userQueries.detail(userId));

  const counter = secondsToCounter(durationFromTimestamps(timestamps));

  return (
    <div className="flex flex-col container mx-auto justify-center align-middle gap-6">
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
          <span className="text-md text-left pl-0.5 text-base-content/90">Total Time</span>
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
      {user.is_superuser && <TeamTimestamps startDate={startDate} endDate={endDate} user={user} />}
    </div>
  );
}
