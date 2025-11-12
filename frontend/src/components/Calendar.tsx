import { Link } from "@tanstack/react-router";
import type { User } from "../types/auth";
import { hexToHSL, hsla } from "../utils/colors";
import type { EventUser, Month, Event } from "../types/response";

export function CalendarNavigation({
  year: currYear,
  month: currMonth,
}: {
  year: number;
  month: number;
}) {
  let year = currYear;
  let nextYear = currYear;
  let nextMonth = currMonth + 1;
  let prevMonth = currMonth - 1;
  if (prevMonth <= 0) {
    prevMonth = 12;
    year--;
  }

  if (nextMonth > 12) {
    nextMonth = nextMonth % 12;
    nextYear++;
    if (nextMonth === 0) nextMonth = 1;
  }

  return (
    <div className="col-span-2 flex justify-start space-x-2 bg-base-100 p-3 rounded-xl gap-4">
      <div className="flex items-center gap-2 w-full">
        <div className="flex justify-center items-center">
          <Link
            to="/calendar/$year/$month"
            params={{ year: year.toString(), month: prevMonth.toString() }}
            className="btn btn-sm btn-soft btn-primary hover:text-neutral icon-outlined animate-color duration-500"
          >
            arrow_back
          </Link>
        </div>
        <div className="flex justify-center items-center">
          <Link
            to="/calendar/$year/$month"
            params={{ year: nextYear.toString(), month: nextMonth.toString() }}
            className="btn btn-sm btn-soft btn-primary hover:text-neutral icon-outlined animate-color duration-500"
          >
            arrow_forward
          </Link>
        </div>
        <div className="pl-4 text-lg">
          <p>month.Name strYear</p>
        </div>
      </div>
    </div>
  );
}

export function UserFilter({ users }: { users: User[] }) {
  return (
    <select className="w-full col-span-1 cursor-pointer bg-base-100 select hover:text-white border-0 hover:bg-[#6F78EA] text-center focus:outline-0 h-full text-lg rounded-xl animate-color">
      <option value="all">All Users</option>
      {users.map((u, i) => (
        <option key={i} onSelect={() => {}} value={u.username}>
          {u.username}
        </option>
      ))}
    </select>
  );
}

export function EventFilter({ events }: { events: string[] }) {
  return (
    <select className="w-full col-span-1 cursor-pointer bg-base-100 select hover:text-white border-0 hover:bg-[#6F78EA] text-center focus:outline-0 h-full text-lg rounded-xl animate-color">
      <option value="all">All Events</option>
      {events.map((e, i) => (
        <option key={i} onSelect={() => {}} value={e}>
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
        {used} <span className="!truncate text-info/80 px-2">used</span>
      </div>{" "}
      <span className="text-info/30 px-2">|</span>
      <div
        className="flex flex-wrap justify-center tooltip text-secondary cursor-help"
        data-tip="remaining"
      >
        {remaining}{" "}
        <span className="!truncate text-info/80 px-2">remaining</span>
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

export function Event({ event }: { event: EventUser }) {
  // {{
  // 	cfg := config.GetConfig()
  // 	h, s, l := domain.Color.HexToHSL(event.User.Color)
  // 	borderColor := fmt.Sprintf("hsla(%.0f, %.1f%%, %.1f%%, 0.6)", h, s*100, l*100)
  // 	bgColor := fmt.Sprintf("hsla(%.0f, %.1f%%, %.1f%%, 0.20)", h, s*100, l*100)
  // 	if !event.Event.IsVacation() && event.User.Username != cfg.BotName {
  // 		h, s, l := domain.Color.HexToHSL(event.User.Color)
  // 		borderColor = fmt.Sprintf("hsla(%.0f, %.1f%%, %.1f%%, 0.30)", h, s*100, l*100)
  // 		bgColor = fmt.Sprintf("hsla(%.0f, %.1f%%, %.1f%%, 0.10)", h, s*100, l*100)
  // 		// Auch hier Transparenz hinzufügen falls gewünscht
  // 	}
  // 	eventId := fmt.Sprintf("event-%v", event.Event.ID)
  // 	deleteUrl := fmt.Sprintf("#%v", eventId)
  // }}
  //
  //

  const hsl = hexToHSL(event.user.color);
  const bgColor = hsla(...hsl, 0.2);
  const borderColor = hsla(...hsl, 0.3);
  const deletable = true;
  return (
    <div className="indicator w-full">
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
      <div
        style={{ backgroundColor: bgColor, borderColor: borderColor }}
        className="group relative text-center border py-1 w-full rounded-lg"
      >
        {deletable ? (
          <>
            <span className="flex items-center justify-center text-transparent rounded-lg group-hover:mix-blend-revert group-hover:text-base-content absolute top-0 left-0 w-full h-full icon-outlined animate-all">
              <button className="h-9 w-12 text-center rounded-xl cursor-pointer hover:drop-shadow-lg icon-outlined group-hover:bg-neutral/30 group-hover:text-base-content hover:text-error-content hover:duration-1000 hover:bg-error/80 hover:w-full hover:h-full hover:rounded-lg animate-all">
                delete
              </button>
            </span>
            <div className="text-base-content group-hover:text-white/0 animate-all">
              {event.event.name.toUpperCase()}
            </div>
            <div className="pb-1 text-xs text-base-content/80 group-hover:text-white/0 animate-all">
              {event.user.username}
            </div>
          </>
        ) : (
          <>
            <div className="text-base-content animate-all">
              {event.event.name.toUpperCase()}
            </div>
            <div className="pb-1 text-xs text-base-content/80 animate-all">
              {event.user.username}
            </div>
          </>
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
}: {
  date: number;
  day: string;
  month: number;
  year: number;
  events: EventUser[];
}) {
  const now = new Date(Date.now());
  const isToday =
    now.getDay() === date &&
    now.getMonth() === month &&
    now.getFullYear() === year;
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
          {events?.map((e) => (
            <Event event={e} />
          ))}
        </div>
        <button className="my-2 btn btn-sm border border-dashed border-primary/30 hover:bg-primary/10 rounded-lg text-primary hover:text-base-content w-full hover:icon-filled">
          <span className="icon-outlined hover:icon-filled text-2xl leading-5">
            add
          </span>
        </button>
      </div>
    </div>
  );
}

export function Calendar({ month }: { month: Month }) {
  return (
    <div className="my-12 lg:my-8 lg:mt-0 mx-auto grid px-6 grid-cols-1 gap-y-6 lg:grid-cols-7 lg:px-4 gap-x-2 lg:gap-y-4 overflow-x-scroll">
      <WeekdayHeader label="Monday" />
      <WeekdayHeader label="Tuesday" />
      <WeekdayHeader label="Wednesday" />
      <WeekdayHeader label="Thursday" />
      <WeekdayHeader label="Friday" />
      <WeekdayHeader label="Saturday" />
      <WeekdayHeader label="Sunday" />
      {Array(month.offset).map((_, i) => (
        <div key={i} className="hidden lg:block"></div>
      ))}
      {month.days.map((d, i) => {
        const date = new Date(d.date);
        return (
          <Day
            key={i}
            date={d.number}
            day={d.name}
            month={date.getMonth()}
            year={date.getFullYear()}
            events={d.events}
          />
        );
      })}
    </div>
  );
}
