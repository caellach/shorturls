import { createSlice } from "@reduxjs/toolkit";

export type AuthClaims = {
  sub: string;
  username: string;
  avatar?: string;
  provider: string;
  provider_sub: string;
  exp: number;
};

export type AuthState = {
  isLoggedIn: boolean;
  claims: AuthClaims;
  accessToken: string;
  refreshToken: string;
};

const initialState: AuthState = {
  isLoggedIn: false,
  claims: {
    sub: "",
    username: "",
    provider: "",
    provider_sub: "",
    exp: 0,
  },
  accessToken: "",
  refreshToken: "",
};

export const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    login: (state, action) => {
      if (action.payload.authTokens) {
        const decodedToken = JSON.parse(
          atob(action.payload.authTokens.access_token.split(".")[1]),
        );

        state.claims.sub = decodedToken.sub ?? "";
        state.claims.username = decodedToken.username ?? "";
        state.claims.provider = decodedToken.provider ?? "";
        state.claims.provider_sub = decodedToken.provider_sub ?? "";
        state.claims.exp = decodedToken.exp ?? 0;
        state.claims.avatar = decodedToken.avatar ?? "";
        state.accessToken = action.payload.authTokens.access_token ?? "";
        state.refreshToken = action.payload.authTokens.refresh_token ?? "";

        state.isLoggedIn = true;
      }
    },
    logout: (state) => {
      //localStorage.removeItem("token");
      state.isLoggedIn = false;
      state.claims = {
        sub: "",
        username: "",
        provider: "",
        provider_sub: "",
        exp: 0,
        avatar: "",
      };
      state.accessToken = "";
      state.refreshToken = "";
    },
  },
});

export const { login, logout } = authSlice.actions;

export default authSlice.reducer;
