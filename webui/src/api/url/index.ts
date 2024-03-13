import { ShorturlData, UserMetadata } from "@/types/ShorturlData";
import axios from "axios";

export const getUrls = async (): Promise<ShorturlData[]> => {
  const response = await axios.get("/u/");
  return response.data as ShorturlData[];
};

export const getUserMetadata = async (): Promise<UserMetadata> => {
  const response = await axios.get("/u/metadata");
  return response.data;
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
  const response = await axios.put("/u/", { url });
  return response.data as ShorturlData;
};

export const getAuthenticatedWebsocket = (): WebSocket => {
  const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
  const wsBaseUrl =
    window.location.hostname === "localhost"
      ? axios.defaults.baseURL?.split("://")[1] ?? "localhost:3000"
      : window.location.host;
  const ws = new WebSocket(`${wsProtocol}://${wsBaseUrl}/u/ws`);
  // get the token set on the axios instance
  const authorization = axios.defaults.headers.common.Authorization;
  if (!authorization) {
    throw new Error("No token set on axios instance");
  }
  if (typeof authorization !== "string") {
    throw new Error("Token is not a string");
  }
  const token = authorization.split(" ")[1];
  if (!token) {
    throw new Error("No token found in Authorization header");
  }

  ws.onopen = () => {
    console.log("Websocket connected, authenticating...");
    ws.send(
      JSON.stringify({
        action: "auth",
        token,
      }),
    );
  };

  return ws;
};
