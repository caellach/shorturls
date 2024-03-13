import { logout } from "@/features/auth/authSlice";
import { StoreState } from "@/store";
import axios from "axios";
import { ComponentType, FC } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Navigate } from "react-router-dom";

function requireAuth(WrappedComponent: ComponentType, requireAuth = true) {
  const RequiresAuth: FC = (props) => {
    const auth = useSelector((state: StoreState) => state.auth);
    const dispatch = useDispatch();

    if (auth.accessToken) {
      // Check if the token is valid
      const token = auth.accessToken;
      const base64Url = token.split(".")[1];
      const decodedToken = JSON.parse(atob(base64Url));
      const exp = decodedToken.exp;
      const currentTime = Date.now() / 1000;
      if (exp < currentTime) {
        dispatch(logout());
      }
    }

    const isAuthorized = auth.isLoggedIn;
    if (isAuthorized) {
      if (window.location.hostname === "localhost") {
        axios.defaults.baseURL = "http://localhost:3000";
      }
      axios.defaults.headers.common.Authorization = `Bearer ${auth.accessToken}`;
    }
    if (requireAuth && !isAuthorized) {
      // Redirect to the sign-in page if not authorized
      return <Navigate to="/" replace />;
    }

    // Render the wrapped component if authorized
    return <WrappedComponent {...props} />;
  };
  return RequiresAuth;
}

export default requireAuth;
