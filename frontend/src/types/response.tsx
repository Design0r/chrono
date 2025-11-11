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
