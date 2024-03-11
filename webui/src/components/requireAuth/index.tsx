import { StoreState } from "@/store";
import axios from "axios";
import { ComponentType, FC } from "react";
import { useSelector } from "react-redux";
import { Navigate } from "react-router-dom";

function requireAuth(WrappedComponent: ComponentType, requireAuth = true) {
  const RequiresAuth: FC = (props) => {
    const auth = useSelector((state: StoreState) => state.auth);
    const isAuthorized = auth.isLoggedIn; // Replace this with your actual authorization check

    if (isAuthorized) {
      axios.defaults.baseURL = "http://localhost:3000";
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
