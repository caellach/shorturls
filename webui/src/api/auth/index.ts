import axios from "axios";

// Purpose: Contains the functions for authenticating the user.

axios.defaults.baseURL = "http://localhost:3000";

export const getOAuthUrl = async (provider: string) => {
  const response = await axios.get(`/api/auth/${provider}`);
  return response.data;
};

export const GetToken = async (provider: string, code: string) => {
  //const response = await axios.post("/api/auth", { provider, code });
  // mock the return for now; needs to be a valid JWT token
  code; // remove the unused variable warning
  const response = await generateAccessToken("123", "test", provider);
  console.log(response);
  return response.data;
};

const generateAccessToken = async (
  userId: string,
  username: string,
  provider: string,
) => {
  const header = {
    alg: "HS256",
    typ: "JWT",
  };
  const payload = {
    sub: userId,
    username: username,
    auth_provider: provider,
  };

  const jwt_parts = [
    btoa(JSON.stringify(header)),
    btoa(JSON.stringify(payload)),
    btoa("signature"),
  ];

  const accessToken = jwt_parts.join(".");

  const response = {
    data: {
      access_token: accessToken,
      refresh_token: "refresh_token_here",
    },
  };
  return response;
};
