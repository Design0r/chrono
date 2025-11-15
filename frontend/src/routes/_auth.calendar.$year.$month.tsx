import { createFileRoute } from "@tanstack/react-router";
import {
  Calendar,
  CalendarNavigation,
  EventFilter,
  UserFilter,
  VacationCounter,
} from "../components/Calendar";
import type { UserWithVacation } from "../types/auth";
import { useQuery } from "@tanstack/react-query";
import { LoadingSpinnerPage } from "../components/LoadingSpinner";
import { ErrorPage } from "../components/ErrorPage";
import { useState } from "react";

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
  const year = Number.parseInt(params.year);
  const month = Number.parseInt(params.month);

  const [userFilter, setUserFilter] = useState<string | undefined>(search.user);
  const [eventFilter, setEventFilter] = useState<string | undefined>(
    search.event,
  );
  const [selectedEvent, setSelectedEvent] = useState<string>("Urlaub");

  const usersQ = useQuery({
    queryKey: ["users", "vacation", year],
    queryFn: () => chrono.users.getUsers({ year: year }),
    staleTime: 1000 * 60 * 30, // 30min
    gcTime: 1000 * 60 * 60 * 1, // 1h
  });

  const currUserQ = useQuery({
    queryKey: ["user", auth.userId, "vacation", year],
    queryFn: () => chrono.users.getUserById(auth.userId!, { year: year }),
    staleTime: 1000 * 60 * 60 * 6, // 6h
    gcTime: 1000 * 60 * 60 * 7, // 7h
  });

  const monthQ = useQuery({
    queryKey: ["month", params.year, params.month],
    queryFn: () => chrono.events.getEventsForMonth(year, month),
    staleTime: 1000 * 60 * 1, // 1min
    gcTime: 1000 * 60 * 30, // 30min
  });

  const queries = [usersQ, currUserQ, monthQ];
  const anyPending = queries.some((q) => q.isPending);
  const firstError = queries.find((q) => q.isError)?.error;

  if (anyPending) return <LoadingSpinnerPage />;
  if (firstError) return <ErrorPage error={firstError} />;

  const users = usersQ.data!;
  const currUser = currUserQ.data! as UserWithVacation;
  const monthData = monthQ.data!;

  return (
    <div>
      <div className="grid grid-cols-7">
        <div className="px-6 col-span-7 grid grid-cols-1 grid-rows-4 lg:grid-rows-1 items-center gap-y-2 lg:gap-y-0 lg:gap-x-2 mt-2 lg:mb-16 lg:grid-cols-7 lg:px-4 ">
          <select
            onChange={(e) => setSelectedEvent(e.target.value)}
            className=" col-span-1 lg:col-span-1 cursor-pointer bg-base-100 select hover:text-white border-0 hover:bg-[#6F78EA] text-center focus:outline-0 h-full w-full text-lg rounded-xl animate-color"
            name="eventName"
            id="eventName"
          >
            <option value="urlaub">Urlaub</option>
            <option value="urlaub halbtags">Urlaub Halbtags</option>
            <option value="workation">Workation</option>
            <option value="krank">Krank</option>
            <option value="home office">Home Office</option>
          </select>
          <CalendarNavigation
            year={year}
            month={month}
            monthName={monthData.name}
          />
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
                    events={[
                      "Krank",
                      "Home Office",
                      "Urlaub",
                      "Urlaub Halbtags",
                      "Workation",
                    ]}
                  />
                </div>
                <div className="flex col-span-2 w-full items-center align-middle h-full ">
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
      <Calendar
        selectedEvent={selectedEvent}
        eventFilter={eventFilter}
        userFilter={userFilter}
        month={monthData}
        currUser={currUser}
      />
    </div>
  );
}
