import { StoreState } from "@/store";
import { ComponentType, FC } from "react";
import { useSelector } from "react-redux";
import { Navigate } from "react-router-dom";

function checkAuthorization() {
  const auth = useSelector((state: StoreState) => state.auth);
  return auth.isLoggedIn;
}

function requireAuth(WrappedComponent: ComponentType, requireAuth = true) {
  const RequiresAuth: FC = (props) => {
    const isAuthorized = checkAuthorization(); // Replace this with your actual authorization check

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
