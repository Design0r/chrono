import type { Timestamp } from "../types/response";

export type TimeCounter = {
  hours: number;
  minutes: number;
  seconds: number;
};

export type WeekGroup = {
  weekMonday: Date;
  weekSunday: Date;
  timestamps: Timestamp[];
  totalSeconds: number;
};

export type DayGroup = {
  dayDate: Date;
  dateKey: string;
  timestamps: Timestamp[];
  totalSeconds: number;
};

const weekdayFormat = new Intl.DateTimeFormat("de-DE", { weekday: "short" });
const timeFormat = new Intl.DateTimeFormat("de-DE", { hour: "2-digit", minute: "2-digit" });
const dateFormat = new Intl.DateTimeFormat("de-DE", {
  day: "2-digit",
  month: "2-digit",
  year: "numeric",
});
const dayLabelFormat = new Intl.DateTimeFormat("de-DE", {
  weekday: "short",
  day: "2-digit",
  month: "2-digit",
  year: "numeric",
});

export function formatWeekdayShort(date: Date): string {
  return weekdayFormat.format(date).slice(0, 2);
}

export function formatTime(date: Date): string {
  return timeFormat.format(date);
}

export function formatDayLabel(date: Date): string {
  return dayLabelFormat.format(date);
}

export function durationFromTimestamps(timestamps: Timestamp[]): number {
  return timestamps
    .map((t) => {
      const start = new Date(t.start_time);
      const end = t.end_time && new Date(t.end_time);
      return end ? (end.getTime() - start.getTime()) / 1000 : 0;
    })
    .reduce((acc, curr) => acc + curr, 0);
}

/** ISO 8601 week number (1–53). Week = Mon–Sun, week 1 = week with first Thursday. */
export function getWeekNumber(date: Date): number {
  const d = new Date(date);
  d.setHours(0, 0, 0, 0);
  const day = d.getDay() || 7; // 1 = Mon, 7 = Sun
  d.setDate(d.getDate() + 4 - day); // Thursday of this week
  const jan1 = new Date(d.getFullYear(), 0, 1);
  return 1 + Math.floor((d.getTime() - jan1.getTime()) / 86400000 / 7);
}

/** Monday 00:00:00 (local) of the week containing the given date. Week = Mon–Sun. */
export function getWeekMonday(date: Date): Date {
  const d = new Date(date);
  d.setHours(0, 0, 0, 0);
  const day = d.getDay();
  const diff = day === 0 ? -6 : 1 - day;
  d.setDate(d.getDate() + diff);
  return d;
}

export function formatWeekRange(monday: Date, sunday: Date): string {
  return `${dateFormat.format(monday)} – ${dateFormat.format(sunday)}`;
}

export function isWeekend(date: Date): boolean {
  const day = date.getDay();
  return day === 0 || day === 6;
}

export function secondsToCounter(totalSeconds: number): TimeCounter {
  const seconds = Math.max(0, Math.floor(totalSeconds));
  const hours = Math.floor(seconds / 60 / 60);
  const minutes = Math.floor(seconds / 60) % 60;
  const s = seconds % 60;
  return { hours, minutes, seconds: s };
}

export function formatCounter(c: TimeCounter): { hours: string; minutes: string; seconds: string } {
  return {
    hours: String(c.hours).padStart(2, "0"),
    minutes: String(c.minutes).padStart(2, "0"),
    seconds: String(c.seconds).padStart(2, "0"),
  };
}

/** Groups timestamps by calendar day (local date), newest first. */
export function groupTimestampsByDay(
  timestamps: Timestamp[],
  durationFn: (ts: Timestamp[]) => number,
): DayGroup[] {
  const byDay = new Map<string, Timestamp[]>();
  for (const t of sortByStart(timestamps)) {
    const key = dateKeyFromDate(new Date(t.start_time));
    const day = byDay.get(key);
    if (day) day.push(t);
    else byDay.set(key, [t]);
  }
  const groups: DayGroup[] = [];
  byDay.forEach((ts, key) => {
    const [y, m, d] = key.split("-").map(Number);
    groups.push({
      dayDate: new Date(y, m - 1, d),
      dateKey: key,
      timestamps: ts,
      totalSeconds: durationFn(ts),
    });
  });
  groups.sort((a, b) => b.dayDate.getTime() - a.dayDate.getTime());
  return groups;
}

/** Groups timestamps by calendar week (Mon–Sun), newest week first. */
export function groupTimestampsByWeek(
  timestamps: Timestamp[],
  durationFn: (ts: Timestamp[]) => number,
): WeekGroup[] {
  const byWeek = new Map<number, Timestamp[]>();
  for (const t of sortByStart(timestamps)) {
    const key = getWeekMonday(new Date(t.start_time)).getTime();
    const week = byWeek.get(key);
    if (week) week.push(t);
    else byWeek.set(key, [t]);
  }
  const weekSunday = (m: Date) => {
    const s = new Date(m);
    s.setDate(s.getDate() + 6);
    return s;
  };
  const groups: WeekGroup[] = [];
  byWeek.forEach((ts, key) => {
    const weekMonday = new Date(key);
    groups.push({
      weekMonday,
      weekSunday: weekSunday(weekMonday),
      timestamps: ts,
      totalSeconds: durationFn(ts),
    });
  });
  groups.sort((a, b) => b.weekMonday.getTime() - a.weekMonday.getTime());
  return groups;
}

function sortByStart(timestamps: Timestamp[]): Timestamp[] {
  return timestamps
    .map((t) => ({ t, start: Date.parse(t.start_time) }))
    .sort((a, b) => a.start - b.start)
    .map(({ t }) => t);
}

export function dateKeyFromDate(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

/** ISO ("2025-12-23T20:44:00Z") → datetime-local ("2025-12-23T20:44") */
export function isoToDatetimeLocal(iso: string): string {
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, "0");
  if (d.getSeconds() === 0) {
    d.setSeconds(new Date().getSeconds());
  }
  const yyyy = d.getFullYear();
  const mm = pad(d.getMonth() + 1);
  const dd = pad(d.getDate());
  const hh = pad(d.getHours());
  const min = pad(d.getMinutes());
  return `${yyyy}-${mm}-${dd}T${hh}:${min}`;
}

export function isoToDateLocal(iso: string): string {
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, "0");
  const yyyy = d.getFullYear();
  const mm = pad(d.getMonth() + 1);
  const dd = pad(d.getDate());
  return `${yyyy}-${mm}-${dd}`;
}

/** date ("2025-12-23") → ISO UTC ("2025-12-23T00:00:00Z") */
export function dateLocalToIso(value: string): string {
  return `${value}T00:00:00Z`;
}

/** datetime-local ("2025-12-23T21:44") → ISO UTC ("2025-12-23T20:44:00Z") */
export function datetimeLocalToIso(value: string): string {
  const d = new Date(value);
  if (d.getSeconds() === 0) {
    d.setSeconds(new Date().getSeconds());
  }
  const iso = d.toISOString();
  return `${iso.split(".")[0]}Z`;
}

export function yearOfRange(startDate?: string): number {
  return startDate ? new Date(startDate).getFullYear() : new Date().getFullYear();
}
