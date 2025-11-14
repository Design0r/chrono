import type { User } from "./auth";

export type ChronoResponse = {
  message: string;
  data: any | null;
};

export type Notification = {
  id: string;
  message: string;
  created_at: Date;
  viewed_at: Date;
};

export type Month = {
  name: string;
  number: number;
  year: number;
  days: Day[];
  offset: number;
};

export type Day = {
  name: string;
  number: number;
  events: EventUser[];
  date: Date;
};

export type Event = {
  id: number;
  scheduled_at: Date;
  created_at: Date;
  edited_at: Date;
  name: string;
  state: "accepted" | "declined" | "pending";
  user_id: number;
};

export type EventUser = { event: Event; user: User };

export type VacationGraphMonth = {
  is_holiday: boolean;
  count: number;
  last_day_of_month: number;
  is_current_week: boolean;
  usernames: string[] | null;
  date: string;
};

export type VacationGraph = {
  month_gaps: number[];
  year_offset: number;
  vacation_data: VacationGraphMonth[];
};
