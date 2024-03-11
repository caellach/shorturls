import { ShorturlData } from "@/types/ShorturlData";
import axios from "axios";

// Purpose: Contains the functions for authenticating the user.

axios.defaults.baseURL = "http://localhost:3000";

export const getUrls = async (): Promise<ShorturlData[]> => {
  const response = await axios.get("/u/");
  return response.data as ShorturlData[];
};

export const deleteUrl = async (id: string): Promise<boolean> => {
  const response = await axios.delete(`/u/${id}`);
  if (response.status !== 200) {
    console.error(`Failed to delete shorturl: ${response.status} ${id}`);
    return false;
  }
  return true;
};

export const createUrl = async (url: string): Promise<ShorturlData> => {
  const response = await axios.post("/u/", { url });
  return response.data as ShorturlData;
};
