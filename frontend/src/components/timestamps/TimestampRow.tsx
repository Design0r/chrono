import { formatCounter, secondsToCounter } from "../../lib/timestamp-utils";
import type { User } from "../../types/auth";
import type { Timestamp } from "../../types/response";

export function TimestampRow({
  t,
  onRowClick,
}: {
  t: Timestamp;
  user: User;
  onRowClick: () => void;
}) {
  const start = new Date(t.start_time);
  const end = t.end_time && new Date(t.end_time);
  const duration = end
    ? secondsToCounter((end.getTime() - start.getTime()) / 1000)
    : { hours: 0, minutes: 0, seconds: 0 };
  const f = formatCounter(duration);

  return (
    <tr
      onClick={onRowClick}
      key={t.id}
      className="hover:bg-base-300 *:font-extralight *:text-info/70 bg-base-200/40"
    >
      <td>
        <div className="flex gap-1.5 items-center pl-6 text-info/60">
          <span className="w-6">
            {start
              .toLocaleDateString("de-DE", { weekday: "short" })
              .slice(0, 2)}
            .,
          </span>
          <span className="block">
            {start.toLocaleTimeString("de-DE", {
              hour: "2-digit",
              minute: "2-digit",
            })}
          </span>
        </div>
      </td>
      <td>
        {end ? (
          <div className="flex gap-1.5 items-center text-info/60">
            <span className="w-6">
              {end
                .toLocaleDateString("de-DE", { weekday: "short" })
                .slice(0, 2)}
              .,
            </span>
            <span className="block">
              {end.toLocaleTimeString("de-DE", {
                hour: "2-digit",
                minute: "2-digit",
              })}
            </span>
          </div>
        ) : (
          "–"
        )}
      </td>
      <td className="text-info/60">
        <div className="text-info/60">
          <span>{f.hours}</span>
          <span>:</span>
          <span>{f.minutes}</span> h
        </div>
      </td>
    </tr>
  );
}
