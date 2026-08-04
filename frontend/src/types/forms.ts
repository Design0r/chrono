export type ProfileEditForm = {
  username: string;
  password: string;
  email: string;
  color: string;
  awork_id: string;
  workday_hours: number;
  workdays_week: number;
  start_date: string;
};

export type TeamEditForm = {
  vacation_days: number;
  enabled: boolean;
  role: string;
};
