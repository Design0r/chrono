export type ProfileEditForm = {
  username: string;
  password: string;
  email: string;
  color: string;
  awork_id: string;
};

export type TeamEditForm = {
  vacation_days: number;
  enabled: boolean;
  role: string;
};
