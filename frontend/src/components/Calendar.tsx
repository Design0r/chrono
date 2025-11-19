import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate, useParams } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { ChronoClient } from "../api/chrono/client";
import type { User } from "../types/auth";
import type { EventUser, Month } from "../types/response";
import { hexToHSL, hsla } from "../utils/colors";
import { LoadingSpinner } from "./LoadingSpinner";
import { useToast } from "./Toast";

export function CalendarNavigation({
  year: currYear,
  month: currMonth,
  monthName,
}: {
  year: number;
  month: number;
  monthName: string;
}) {
  let year = currYear;
  let nextYear = currYear;
  let nextMonth = currMonth + 1;
  let prevMonth = currMonth - 1;
  if (prevMonth < 0) {
    prevMonth = 11;
    year--;
  }

  if (nextMonth > 11) {
    nextMonth = nextMonth % 12;
    nextYear++;
  }

  return (
    <div className="col-span-2 flex justify-start space-x-2 bg-base-100 p-3 rounded-xl gap-4">
      <div className="flex items-center gap-2 w-full">
        <div className="flex justify-center items-center">
          <Link
            to="/calendar/$year/$month"
            params={{ year: year.toString(), month: prevMonth.toString() }}
            className="btn btn-sm btn-soft btn-primary hover:text-neutral icon-outlined animate-color duration-500"
            search={(prev) => prev}
          >
            arrow_back
          </Link>
        </div>
        <div className="flex justify-center items-center">
          <Link
            to="/calendar/$year/$month"
            params={{ year: nextYear.toString(), month: nextMonth.toString() }}
            className="btn btn-sm btn-soft btn-primary hover:text-neutral icon-outlined animate-color duration-500"
            search={(prev) => prev}
          >
            arrow_forward
          </Link>
        </div>
        <div className="pl-4 text-lg">
          <p>
            {monthName} {year}
          </p>
        </div>
      </div>
    </div>
  );
}

export function UserFilter({
  users,
  userFilter,
  setUserFilter,
}: {
  users: User[];
  userFilter?: string;
  setUserFilter: (value?: string) => void;
}) {
  const navigate = useNavigate();
  const params = useParams({ from: "/_auth/calendar/$year/$month" });

  return (
    <select
      defaultValue={userFilter}
      onChange={(e) => {
        const filtered =
          e.target.value === "allUsers" ? undefined : e.target.value;
        setUserFilter(filtered);
        navigate({
          to: "/calendar/$year/$month",
          search: (prev) => ({ ...prev, user: filtered }),
          params: params,
        });
      }}
      className="w-full col-span-1 cursor-pointer bg-base-100 select hover:text-white border-0 hover:bg-[#6F78EA] text-center focus:outline-0 h-full text-lg rounded-xl animate-color"
    >
      <option value="allUsers">All Users</option>
      {users.map((u, i) => (
        <option key={i} value={u.username}>
          {u.username}
        </option>
      ))}
    </select>
  );
}

export function EventFilter({
  events,
  eventFilter,
  setEventFilter,
}: {
  events: string[];
  eventFilter?: string;
  setEventFilter: (value: string | undefined) => void;
}) {
  const navigate = useNavigate();
  const params = useParams({ from: "/_auth/calendar/$year/$month" });

  return (
    <select
      defaultValue={eventFilter}
      onChange={(e) => {
        const filtered =
          e.target.value === "allEvents" ? undefined : e.target.value;
        setEventFilter(filtered);
        navigate({
          to: "/calendar/$year/$month",
          search: (prev) => ({ ...prev, event: filtered }),
          params: params,
        });
      }}
      className="w-full col-span-1 cursor-pointer bg-base-100 select hover:text-white border-0 hover:bg-[#6F78EA] text-center focus:outline-0 h-full text-lg rounded-xl animate-color"
    >
      <option value={"allEvents"}>All Events</option>
      {events.map((e, i) => (
        <option key={i} value={e.toLowerCase()}>
          {e}
        </option>
      ))}
    </select>
  );
}

export function VacationCounter({
  pending,
  used,
  remaining,
}: {
  pending: number;
  used: number;
  remaining: number;
}) {
  return (
    <div className="flex px-3 w-full justify-center bg-base-100 items-center rounded-xl align-middle h-full text-center">
      <div
        className="flex flex-wrap justify-center tooltip text-primary cursor-help"
        data-tip="pending"
      >
        {pending}
        <span className="text-info/80 px-2">pending</span>
      </div>
      <span className="text-info/30 px-2">|</span>
      <div
        className="flex flex-wrap justify-center tooltip text-warning cursor-help"
        data-tip="used"
      >
        {used} <span className="truncate! text-info/80 px-2">used</span>
      </div>{" "}
      <span className="text-info/30 px-2">|</span>
      <div
        className="flex flex-wrap justify-center tooltip text-secondary cursor-help"
        data-tip="remaining"
      >
        {remaining}{" "}
        <span className="truncate! text-info/80 px-2">remaining</span>
      </div>
    </div>
  );
}

export function WeekdayHeader({
  label,
  highlighted = false,
}: {
  label: string;
  highlighted?: boolean;
}) {
  return (
    <div
      className={`lg:block hidden truncate text-sm ${highlighted ? "text-primary" : ""} rounded-xl p-1 text-center lg:text-lg`}
    >
      {label}
    </div>
  );
}

export function Event({
  event,
  currUser,
}: {
  event: EventUser;
  currUser: User;
}) {
  const chrono = new ChronoClient();
  const queryClient = useQueryClient();
  const { addToast, addErrorToast } = useToast();

  const hsl = hexToHSL(event.user.color);
  const bgColor = hsla(...hsl, 0.2);
  const borderColor = hsla(...hsl, 0.3);

  const [visible, setVisible] = useState(true);

  const currDate = new Date();
  const eventDate = new Date(event.event.scheduled_at);
  const isInFuture = eventDate >= currDate;
  const isFromCurrUser = event.event.user_id === currUser.id;
  const isAdmin = currUser.is_superuser;
  const isHoliday = event.user.id === 1;

  const shortName = isHoliday
    ? event.event.name
    : chrono.events.getShortEventName(event.event.name);

  const isDeletable =
    isAdmin ||
    (isFromCurrUser && (isInFuture || event.event.state !== "accepted"));

  const mutation = useMutation({
    mutationKey: ["deleteEvent", event.event.id],
    mutationFn: () => chrono.events.deleteEvent(event.event.id),
    onSuccess: () => {
      addToast(`Successfully deleted event ${event.event.id}`, "success");
      setVisible(false);
      return queryClient.invalidateQueries({
        queryKey: ["month"],
      });
    },
    onError: (error) => addErrorToast(error),
  });

  if (!visible) return <></>;
  return (
    <div className="indicator w-full">
      {!isHoliday && (
        <div className="absolute top-1.5 right-1.5 z-10 w-3 bg-neutral/50 aspect-square rounded-full flex items-center justify-center">
          <span
            className={
              event.event.state === "pending"
                ? "bg-accent status status-sm status-accent animate-ping"
                : event.event.state === "declined"
                  ? "status status-md status-error"
                  : "status status-md status-success"
            }
          ></span>
        </div>
      )}
      <div
        style={{ backgroundColor: bgColor, borderColor: borderColor }}
        className={`group gap-2 relative ${!isHoliday && "flex"} text-center border py-1 w-full rounded-lg`}
      >
        {isDeletable && (
          <>
            <span className="flex items-center justify-center text-transparent rounded-lg group-hover:mix-blend-revert group-hover:text-base-content absolute top-0 left-0 w-full h-full icon-outlined animate-all">
              <button
                onClick={() => mutation.mutate()}
                className="h-9 w-12 text-center rounded-xl cursor-pointer hover:drop-shadow-lg icon-outlined group-hover:bg-neutral/30 group-hover:text-base-content hover:text-error-content hover:duration-1000 hover:bg-error/80 hover:w-full hover:h-full hover:rounded-lg animate-all"
              >
                delete
              </button>
            </span>
          </>
        )}
        <div
          className={`${!isHoliday && "bg-black/20 rounded px-1"} text-base-content ${isDeletable && "group-hover:text-white/0"}  animate-all`}
        >
          {shortName}
        </div>
        {!isHoliday && (
          <div
            className={`content-center text-base-content/80 ${isDeletable && "group-hover:text-white/0"} animate-all`}
          >
            {event.user.username}
          </div>
        )}
      </div>
    </div>
  );
}

export function Day({
  date,
  day,
  month,
  year,
  events,
  selectedEvent,
  currUser,
}: {
  date: number;
  day: string;
  month: number;
  year: number;
  events: EventUser[];
  selectedEvent: string;
  currUser: User;
}) {
  const now = new Date();
  const isToday =
    now.getDate() === date &&
    now.getMonth() + 1 === month &&
    now.getFullYear() === year;

  const [evts, setEvts] = useState<EventUser[]>(events);

  useEffect(() => setEvts(events), [events]);
  const chrono = new ChronoClient();
  const queryClient = useQueryClient();
  const { addToast, addErrorToast } = useToast();

  const mutation = useMutation({
    mutationKey: ["createEvent", year, month, date, selectedEvent],
    mutationFn: () =>
      chrono.events.createEvent({
        year: year,
        month: month,
        day: date,
        event: selectedEvent,
      }),
    onSuccess: (e) => {
      setEvts((ev) => [...ev, e]);
      addToast(`Successfully created event ${selectedEvent}`, "success");
      return queryClient.invalidateQueries({
        queryKey: ["month"],
      });
    },
    onError: (error) => addErrorToast(error),
  });

  return (
    <div
      className={
        isToday
          ? "border bg-primary/90 border-primary text-neutral rounded-xl flex flex-col overflow-hidden"
          : "bg-base-100 rounded-xl border border-base-300 flex flex-col overflow-hidden"
      }
    >
      <div className="pt-2 pb-2 px-4 text-lg lg:text-center ">
        <div className="lg:hidden text-base font-medium">
          {date} <span className="pl-1.5 opacity-40">{day}</span>
        </div>
        <div className="hidden lg:block">{date}</div>
      </div>
      <div className="flex flex-col px-2 h-full lg:bg-base-200/65 rounded-t-none rounded-b-[0.65rem]">
        <div className="flex flex-col gap-2 h-fit rounded-[0.7rem] *:first:mt-2">
          {evts.map((e, i) => (
            <Event key={i} event={e} currUser={currUser} />
          ))}
        </div>
        <button
          onClick={() => mutation.mutate()}
          className="my-2 btn btn-sm border border-dashed border-primary/30 hover:bg-primary/10 rounded-lg text-primary hover:text-base-content w-full hover:icon-filled"
        >
          {mutation.isPending ? (
            <div className="w-5 flex">
              <LoadingSpinner />
            </div>
          ) : (
            <span className="icon-outlined hover:icon-filled text-2xl leading-5">
              add
            </span>
          )}
        </button>
      </div>
    </div>
  );
}

export function Calendar({
  month,
  eventFilter,
  userFilter,
  selectedEvent,
  currUser,
}: {
  month: Month;
  eventFilter: string | undefined;
  userFilter: string | undefined;
  selectedEvent: string;
  currUser: User;
}) {
  return (
    <div className="my-12 lg:my-8 lg:mt-0 mx-auto grid px-6 grid-cols-1 gap-y-6 lg:grid-cols-7 lg:px-4 gap-x-2 lg:gap-y-4 overflow-x-scroll">
      <WeekdayHeader label="Monday" />
      <WeekdayHeader label="Tuesday" />
      <WeekdayHeader label="Wednesday" />
      <WeekdayHeader label="Thursday" />
      <WeekdayHeader label="Friday" />
      <WeekdayHeader label="Saturday" />
      <WeekdayHeader label="Sunday" />
      {Array.from({ length: month.offset }).map((_, i) => (
        <div key={i} className="hidden lg:block"></div>
      ))}
      {month.days.map((d, i) => {
        const date = new Date(d.date);
        console.log(date, d.date, d.number);

        return (
          <Day
            selectedEvent={selectedEvent}
            key={i}
            date={d.number}
            day={d.name}
            month={date.getMonth() + 1}
            year={date.getFullYear()}
            currUser={currUser}
            events={
              d.events?.filter((e) => {
                let event = true;
                let user = true;
                if (eventFilter && e.user.id !== 1) {
                  event = e.event.name === eventFilter;
                }
                if (userFilter && e.user.id !== 1) {
                  user = e.user.username === userFilter;
                }

                return event && user;
              }) || []
            }
          />
        );
      })}
    </div>
  );
}
