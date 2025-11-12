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
