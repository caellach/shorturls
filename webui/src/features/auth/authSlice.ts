import { createSlice } from "@reduxjs/toolkit";

export type AuthClaims = {
  sub: string;
  username: string;
  displayName?: string;
  avatar?: string;
  locale?: string;
  provider: string;
  providerSub: string;
  providerMfa?: boolean;
  providerVerified?: boolean;
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
    providerSub: "",
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
          atob(action.payload.authTokens.accessToken.split(".")[1]),
        );

        state.claims.sub = decodedToken.sub ?? "";
        state.claims.username = decodedToken.username ?? "";
        state.claims.displayName = decodedToken.displayName ?? "";
        state.claims.avatar = decodedToken.avatar ?? "";
        state.claims.locale = decodedToken.locale ?? "";
        state.claims.provider = decodedToken.provider ?? "";
        state.claims.providerSub = decodedToken.providerSub ?? "";
        state.claims.providerMfa = decodedToken.providerMfa ?? false;
        state.claims.providerVerified = decodedToken.providerVerified ?? false;
        state.claims.exp = decodedToken.exp ?? 0;
        state.accessToken = action.payload.authTokens.accessToken ?? "";
        state.refreshToken = action.payload.authTokens.refreshToken ?? "";

        state.isLoggedIn = true;
      }
    },
    logout: (state) => {
      //localStorage.removeItem("token");
      state.isLoggedIn = false;
      state.claims = {
        sub: "",
        username: "",
        displayName: "",
        avatar: "",
        provider: "",
        providerSub: "",
        providerMfa: false,
        providerVerified: false,
        exp: 0,
      };
      state.accessToken = "";
      state.refreshToken = "";
    },
  },
});

export const { login, logout } = authSlice.actions;

export default authSlice.reducer;
