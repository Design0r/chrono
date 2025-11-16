import type { User } from "./auth";

export type ChronoResponse = {
  message: string;
  data: any | null;
};

export type State = "accepted" | "declined" | "pending";

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
  events: EventUser[] | null;
  date: Date;
};

export type Event = {
  id: number;
  scheduled_at: Date;
  created_at: Date;
  edited_at: Date;
  name: string;
  state: State;
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

export type BatchRequest = {
  start_date: string;
  end_date: string;
  event_count: number;
  request: RequestEventUser;
  conflicts: User[] | null;
};

export type RequestEventUser = {
  request_id: number;
  message: string | null;
  request_state: State;
  created_at: string;
  edited_at: string;
  edited_by: number | null;

  user_id: number;
  username: string;
  email: string;
  vacation_days: number;
  is_superuser: boolean;
  user_created_at: string;
  user_edited_at: string;
  color: string;
  role: string;
  enabled: boolean;

  event_id: number;
  scheduled_at: string;
  name: string;
  event_state: State;
  event_created_at: string;
  event_edited_at: string;
};

export type PatchRequestForm = {
  userId: number;
  state: string;
  reason: string;
  start_date: string;
  end_date: string;
};

export type Settings = {
  id: number;
  signup_enabled: boolean;
};
