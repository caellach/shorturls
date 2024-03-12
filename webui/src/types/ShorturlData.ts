export type ShorturlData = {
  id: string;
  short_url: string;
  url: string;
  use_count: number;
  last_used: Date;
};

export type UserMetadata = {
  user_id: string;
  active_count: number;
  created_count: number;
  last_created: Date;
};
