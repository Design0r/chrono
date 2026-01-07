export type LoginRequest = {
  email: string;
  password: string;
};

export type SignupRequest = {
  email: string;
  username: string;
  password: string;
};

export type User = {
  id: number;
  username: string;
  email: string;
  vacation_days: number;
  is_superuser: boolean;
  created_at: string;
  edited_at: string;
  color: string;
  role: string;
  enabled: boolean;
  awork_id: string | null;
  workday_hours: number;
  workdays_week: number;
};

export type UserWithVacation = User & {
  vacation_remaining: number;
  vacation_used: number;
  pending_events: number;
};
