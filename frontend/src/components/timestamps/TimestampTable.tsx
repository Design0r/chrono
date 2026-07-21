import { useQueries } from "@tanstack/react-query";
import { Fragment, useMemo, useState } from "react";
import { ChronoClient } from "../../api/chrono/client";
import {
  dateKeyFromDate,
  durationFromTimestamps,
  formatCounter,
  groupTimestampsByDay,
  secondsToCounter,
} from "../../lib/timestamp-utils";
import type { User } from "../../types/auth";
import type { Timestamp } from "../../types/response";
import { DayBar } from "./DayBar";
import { EditTimestampModal } from "./EditTimestampModal";
import { TimestampRow } from "./TimestampRow";

export function TimestampTable({
  timestamps,
  user,
  footerTotalSeconds,
}: {
  timestamps: Timestamp[];
  user: User;
  footerTotalSeconds?: number;
}) {
  const chrono = useMemo(() => new ChronoClient(), []);
  const [modal, setModal] = useState<Timestamp | null>(null);
  const [expandedDays, setExpandedDays] = useState<Set<string>>(new Set());

  const dayGroups = useMemo(
    () => groupTimestampsByDay(timestamps, durationFromTimestamps),
    [timestamps],
  );

  const monthKeys = useMemo(() => {
    const set = new Set<string>();
    for (const g of dayGroups) {
      const y = g.dayDate.getFullYear();
      const m = g.dayDate.getMonth() + 1;
      set.add(`${y}-${m}`);
    }
    return Array.from(set);
  }, [dayGroups]);

  const monthQueries = useQueries({
    queries: monthKeys.map((key) => {
      const [y, m] = key.split("-").map(Number);
      return {
        queryKey: ["events", "month", y, m],
        queryFn: () => chrono.events.getEventsForMonth(y, m),
        staleTime: 1000 * 60 * 5,
      };
    }),
  });

  const calendarEventsByDay = useMemo(() => {
    const map = new Map<string, string[]>();
    for (const result of monthQueries) {
      const month = result.data;
      if (!month?.days) continue;
      for (const day of month.days) {
        if (!day.events?.length) continue;
        const myEvents = day.events.filter((e) => e.event.user_id === user.id);
        if (!myEvents.length) continue;
        const d = new Date(day.date);
        const key = dateKeyFromDate(d);
        const names = [...new Set(myEvents.map((e) => e.event.name))];
        map.set(key, names);
      }
    }
    return map;
  }, [monthQueries, user.id]);

  const toggleDay = (dateKey: string) => {
    setExpandedDays((prev) => {
      const next = new Set(prev);
      if (next.has(dateKey)) next.delete(dateKey);
      else next.add(dateKey);
      return next;
    });
  };

  const footerCounter =
    footerTotalSeconds !== undefined
      ? secondsToCounter(footerTotalSeconds)
      : null;

  return (
    <>
      <table className="table bg-base-300/50 lg:rounded-none">
        <thead>
          <tr className="text-accent/80 *:w-1/3 *:font-semibold">
            <th>Start</th>
            <th>End</th>
            <th>Duration</th>
          </tr>
        </thead>
        <tbody>
          {dayGroups.map((group) => {
            const isExpanded = expandedDays.has(group.dateKey);
            const dayCounter = secondsToCounter(group.totalSeconds);
            const dayFmt = formatCounter(dayCounter);
            return (
              <Fragment key={group.dateKey}>
                <tr
                  onClick={() => toggleDay(group.dateKey)}
                  className="cursor-pointer hover:bg-base-300/80 bg-base-200/60 select-none border-b border-base-300/60"
                  role="button"
                  tabIndex={0}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      e.preventDefault();
                      toggleDay(group.dateKey);
                    }
                  }}
                >
                  <td className="text-info/90">
                    <span className="inline-flex items-center gap-2">
                      <span
                        className="icon-outlined -ml-1 text-lg transition-transform"
                        style={{
                          transform: isExpanded
                            ? "rotate(90deg)"
                            : "rotate(0deg)",
                        }}
                      >
                        chevron_right
                      </span>
                      {group.dayDate.toLocaleDateString("de-DE", {
                        weekday: "short",
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                      })}
                      {(() => {
                        const eventNames = calendarEventsByDay.get(
                          group.dateKey,
                        );
                        if (!eventNames?.length) return null;
                        const label = eventNames.join(", ");
                        return (
                          <div className="tooltip">
                            <div className="tooltip-content">
                              <div className="text-primary font-bold first-letter:uppercase">
                                {label}
                              </div>
                            </div>
                            <button
                              className="icon-outlined text-base text-primary"
                              data-tip={label}
                              aria-label={`Kalender: ${label}`}
                            >
                              event
                            </button>
                          </div>
                        );
                      })()}
                    </span>
                  </td>
                  <td className="py-1.5">
                    <DayBar
                      totalSeconds={group.totalSeconds}
                      dayDate={group.dayDate}
                    />
                  </td>
                  <td className="text-info/90">
                    <span>{dayFmt.hours}</span>
                    <span>:</span>
                    <span>{dayFmt.minutes}</span>{" "}
                    <span className="text-info/70">h</span>
                  </td>
                </tr>
                {isExpanded &&
                  group.timestamps.map((t) => (
                    <TimestampRow
                      key={t.id}
                      t={t}
                      user={user}
                      onRowClick={() => setModal(t)}
                    />
                  ))}
              </Fragment>
            );
          })}
        </tbody>
        {footerCounter !== null && (
          <tfoot>
            <tr className="bg-base-200/70 border-b-lg font-semibold">
              <td></td>
              <td></td>
              <td className="text-primary">
                {(() => {
                  const f = formatCounter(footerCounter);
                  return (
                    <>
                      <span>{f.hours}</span>
                      <span>:</span>
                      <span>{f.minutes}</span>{" "}
                      <span className="text-primary/80">h</span>
                    </>
                  );
                })()}
              </td>
            </tr>
          </tfoot>
        )}
      </table>
      {modal && (
        <EditTimestampModal
          key={modal.id}
          timestamp={modal}
          onClose={() => setModal(null)}
        />
      )}
    </>
  );
}
