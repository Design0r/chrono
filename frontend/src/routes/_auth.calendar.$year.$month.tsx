import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import {
  Calendar,
  CalendarNavigation,
  EventFilter,
  UserFilter,
  VacationCounter,
} from "../components/Calendar";
import { ErrorPage } from "../components/ErrorPage";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { getEventLabel, getEventMap } from "../lib/events";
import type { UserWithVacation } from "../types/auth";

type TeamSearchParams = {
  user?: string;
  event?: string;
};

export const Route = createFileRoute("/_auth/calendar/$year/$month")({
  component: CalendarComponent,
  validateSearch: (search: Record<string, unknown>): TeamSearchParams => {
    return {
      user: search.user as string,
      event: search.event as string,
    };
  },
});

function CalendarComponent() {
  const { chrono, auth } = Route.useRouteContext();
  const params = Route.useParams();
  const search = Route.useSearch();
  const year = Number(params.year);
  const month = Number(params.month);

  const [userFilter, setUserFilter] = useState<string | undefined>(search.user);
  const [eventFilter, setEventFilter] = useState<string | undefined>(
    search.event,
  );
  const [selectedEvent, setSelectedEvent] = useState<string>(
    getEventLabel("urlaub"),
  );

  const usersQ = useQuery({
    queryKey: ["users", "vacation", year],
    queryFn: () => chrono.users.getUsers({ year: year }),
    staleTime: 1000 * 60 * 30, // 30min
    gcTime: 1000 * 60 * 60 * 1, // 1h
    retry: false,
  });

  const currUserQ = useQuery({
    queryKey: ["user", auth.userId, "vacation", year],
    queryFn: () => chrono.users.getUserById(auth.userId!, { year: year }),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
    retry: false,
  });

  const monthQ = useQuery({
    queryKey: ["month", params.year, params.month],
    queryFn: () => chrono.events.getEventsForMonth(year, month),
    staleTime: 1000 * 60 * 1, // 1min
    gcTime: 1000 * 60 * 30, // 30min
    retry: false,
  });

  const queries = [usersQ, currUserQ, monthQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const users = usersQ.data!;
  const currUser = currUserQ.data! as UserWithVacation;
  const monthData = monthQ.data!;
  const eventMap = getEventMap();

  return (
    <div className="pb-24 lg:pb-0">
      <div className="grid grid-cols-7">
        <div className="px-6 col-span-7 grid grid-cols-1 grid-rows-3  lg:grid-rows-1 items-center gap-y-2 lg:gap-y-0 lg:gap-x-2 mt-2 lg:mb-16 lg:grid-cols-7 lg:px-4 ">
          <select
            onChange={(e) => setSelectedEvent(e.target.value)}
            className="hidden lg:min-h-12 lg:flex w-full col-span-1 cursor-pointer bg-base-300 select hover:text-white focus-within:text-white text-center focus-within:outline-0 h-full pl-4 text-base border-0 rounded-lg animate-all"
            name="eventName"
            id="eventName"
          >
            <>
              {Object.entries(eventMap).map(([event, label]) => (
                <option key={event} value={event}>
                  {label}
                </option>
              ))}
            </>
          </select>
          <div className="hidden lg:block col-span-2">
            <CalendarNavigation
              year={year}
              month={month}
              monthName={monthData.name}
            />
          </div>
          <div className="row-span-2 col-span-2 lg:row-span-1 lg:col-span-4 h-full text-lg">
            <div className="h-full items-center rounded-xl bg-base-200">
              <div className="grid grid-cols-2 lg:grid-cols-4 w-full h-full gap-x-2 gap-y-2 lg:gap-y-0">
                <div className="col-span-1 w-full justify-center ">
                  <UserFilter
                    users={users}
                    userFilter={userFilter}
                    setUserFilter={setUserFilter}
                  />
                </div>
                <div className="col-span-1 w-full justify-center ">
                  <EventFilter
                    setEventFilter={setEventFilter}
                    eventFilter={eventFilter}
                    events={Object.values(eventMap)}
                  />
                </div>
                <div className="flex col-span-2 w-full items-center align-middle h-full">
                  <VacationCounter
                    pending={currUser.pending_events}
                    used={currUser.vacation_used}
                    remaining={currUser.vacation_remaining}
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div className="px-6 lg:px-4">
        <Calendar
          selectedEvent={selectedEvent}
          eventFilter={eventFilter}
          userFilter={userFilter}
          month={monthData}
          currUser={currUser}
        />
      </div>

      {/* Fixed bottom bar on mobile: event select + calendar navigation */}
      <div className="fixed bottom-19 shadow-xl drop-shadow-xs border border-info/10 w-fit mx-auto bg-base-200/50 backdrop-blur-xl rounded-full left-0 right-0 flex justify-between items-center gap-2 px-2 py-1 lg:hidden">
        <select
          onChange={(e) => setSelectedEvent(e.target.value)}
          value={selectedEvent}
          className="min-w-22 max-w-40 cursor-pointer border-0! text-sm bg-base-300 select hover:text-white focus-within:text-white focus-within:outline-0 h-12 pl-4 rounded-full animate-all"
          name="eventNameMobile"
          id="eventNameMobile"
        >
          <>
            {Object.entries(eventMap).map(([event, label]) => (
              <option key={event} value={event}>
                {label}
              </option>
            ))}
          </>
        </select>
        <CalendarNavigation
          year={year}
          month={month}
          monthName={monthData.name}
          compact
        />
      </div>
    </div>
  );
}
