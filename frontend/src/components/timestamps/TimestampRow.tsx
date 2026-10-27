import {
  formatCounter,
  formatTime,
  formatWeekdayShort,
  secondsToCounter,
} from "../../lib/timestamp-utils";
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
            {formatWeekdayShort(start)}
            .,
          </span>
          <span className="block">{formatTime(start)}</span>
        </div>
      </td>
      <td>
        {end ? (
          <div className="flex gap-1.5 items-center text-info/60">
            <span className="w-6">
              {formatWeekdayShort(end)}
              .,
            </span>
            <span className="block">{formatTime(end)}</span>
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
