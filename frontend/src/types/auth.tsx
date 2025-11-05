export type LoginRequest = {
  username: string;
  password: string;
};

export type User = {
  id: number;
  username: string;
  email: string;
  vacation_days: number;
  is_superuser: boolean;
  created_at: Date;
  edited_at: Date;
  color: string;
  role: string;
  enabled: boolean;
};
