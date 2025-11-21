import dayjs from "dayjs";
import isoWeek from "dayjs/plugin/isoWeek";
import { useMemo } from "react";
import type { VacationGraphMonth } from "../types/response";
import { clamp } from "../utils/math";

dayjs.extend(isoWeek);

export function OverviewDay({ day }: { day: VacationGraphMonth }) {
  const greens = [
    "#313745",
    "#a7f3d0",
    "#6ee7b7",
    "#34d399",
    "#10b981",
    "#059669",
    "#047857",
    "#065f46",
    "#064e3b",
  ];
  const holidayColor = "#7C85FF";
  const color = day.is_holiday ? holidayColor : greens[clamp(day.count, 0, 8)];

  return (
    <div className="tooltip">
      <div className="tooltip-content">
        <div>{day.date}</div>
        <div>Count: {day.count}</div>
        <>
          {day.usernames?.map((u, i) => (
            <p key={i}>{u}</p>
          ))}
        </>
      </div>
      <div
        className={`w-4 h-4 rounded-full`}
        style={
          day.last_day_of_month
            ? { boxShadow: "0 1.5rem 0 -0.25rem white", backgroundColor: color }
            : { backgroundColor: color }
        }
      ></div>
    </div>
  );
}

export function VacationGraph({
  gaps,
  yearOffset,
  data,
}: {
  gaps: number[];
  yearOffset: number;
  data: VacationGraphMonth[];
}) {
  const currWeek = dayjs().isoWeek();

  const cells = useMemo(() => {
    let week = 1;
    const out: Array<number> = [];

    for (const g of gaps) {
      out.push(week);
      for (let i = 0; i < g; i++) {
        week++;
        out.push(week);
      }
      week++;
    }

    return out;
  }, [gaps]);

  return (
    <div className="grid grid-cols-12 p-5 bg-base-100 rounded-2xl  xl:overflow-x-hidden overflow-x-auto mb-12">
      <div className="col-span-1" />
      <div className="col-span-11 grid grid-rows-1 grid-flow-col h-10 gap-1 text-secondary/70">
        {gaps.map((g, i) => (
          <div key={`top-${i}`} className="contents">
            <p className="h-4 w-4 text-center">{i + 1}.</p>
            {Array.from({ length: g }).map((j) => (
              <p key={`top-gap-${i}-${j}`} className="w-4 h-4 -z-10" />
            ))}
          </div>
        ))}
      </div>

      <div className="col-span-1" />
      <div className="col-span-11 grid grid-rows-1 grid-flow-col h-10 gap-1 text-sm text-base-content/40">
        {cells.map((week) => (
          <p
            key={week}
            className={`w-4 h-4 text-center ${currWeek === week ? "text-primary font-bold" : "font-light text-secondary/60"}`}
          >
            {week}
          </p>
        ))}
      </div>

      <div className="col-span-1 grid grid-rows-7 text-secondary/60 text-sm">
        <p className="truncate">Montag</p>
        <p className="truncate">Dienstag</p>
        <p className="truncate">Mittwoch</p>
        <p className="truncate">Donnerstag</p>
        <p className="truncate">Freitag</p>
        <p className="opacity-50 truncate">Samstag</p>
        <p className="opacity-50 truncate">Sonntag</p>
      </div>

      <div className="col-span-11 grid grid-rows-7 h-80 grid-flow-col gap-1">
        <div className="contents">
          {Array.from({ length: yearOffset }).map((_, i) => (
            <p key={`yo-${i}`} />
          ))}
          {data.map((d, i) => (
            <OverviewDay key={i} day={d} />
          ))}
        </div>
      </div>
    </div>
  );
}
