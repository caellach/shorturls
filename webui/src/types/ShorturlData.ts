export type ShorturlData = {
  id: string;
  shortUrl: string;
  url: string;
  useCount: number;
  lastUsed: Date;
};

export type UserMetadata = {
  userId: string;
  activeCount: number;
  createdCount: number;
  lastCreated: Date;
};
