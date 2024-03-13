import axios from "axios";

// Purpose: Contains the functions for authenticating the user.

if (window.location.hostname === "localhost") {
  axios.defaults.baseURL = "http://localhost:3000";
}

export const getOAuthUrl = async (provider: string) => {
  const response = await axios.get(`/api/auth/${provider}`);
  return response.data;
};
