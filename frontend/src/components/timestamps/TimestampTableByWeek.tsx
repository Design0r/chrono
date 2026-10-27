import { useMemo } from "react";
import type { User } from "../../types/auth";
import type { Timestamp } from "../../types/response";
import {
  formatWeekRange,
  getWeekNumber,
  groupTimestampsByWeek,
  durationFromTimestamps,
} from "../../lib/timestamp-utils";
import { useCalendarEventsByDay } from "../../lib/useCalendarEventsByDay";
import { TimestampTable } from "./TimestampTable";

export function TimestampTableByWeek({
  timestamps,
  user,
}: {
  timestamps: Timestamp[];
  user: User;
}) {
  const groups = useMemo(
    () => groupTimestampsByWeek(timestamps, durationFromTimestamps),
    [timestamps],
  );
  const eventsByDay = useCalendarEventsByDay(timestamps, user.id);

  return (
    <div className="space-y-6">
      {groups.map((g) => (
        <section
          key={g.weekMonday.getTime()}
          className="[content-visibility:auto] [contain-intrinsic-size:auto_24rem]"
        >
          <h3 className="text-base font-semibold mb-3.5 px-4">
            {getWeekNumber(g.weekMonday)}. KW{" "}
            <span className="text-accent/50 font-normal pl-1 pr-1.5">|</span>
            {formatWeekRange(g.weekMonday, g.weekSunday)}
          </h3>
          <TimestampTable
            timestamps={g.timestamps}
            user={user}
            eventsByDay={eventsByDay}
            footerTotalSeconds={g.totalSeconds}
          />
        </section>
      ))}
    </div>
  );
}
