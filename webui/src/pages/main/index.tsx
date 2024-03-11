import { useEffect } from "react";
import Auth from "@/components/auth/Auth";
import { useDispatch, useSelector } from "react-redux";
import { StoreState } from "@/store";
import { login } from "@/features/auth/authSlice";
import { Navigate, useLocation } from "react-router-dom";

const Main = () => {
  const auth = useSelector((state: StoreState) => state.auth);
  const dispatch = useDispatch();
  const location = useLocation();

  useEffect(() => {
    const handleAuth = async () => {
      const searchParams = new URLSearchParams(location.search);
      // clear the code from the url
      const access_token = searchParams.get("a"); // a for access_token
      if (access_token) {
        console.log("access_token", access_token);
        // dispatch the token to the store
        console.log("auth", auth);
        dispatch(login({ authTokens: { access_token, refresh_token: null } }));
        window.history.replaceState({}, document.title, location.pathname);
      }
    };

    handleAuth();
  }, [auth, dispatch, location]);

  if (auth.isLoggedIn) {
    return <Navigate to="/shorturls" replace />;
  }

  // we should never reach this point if the user is logged in
  return (
    <>
      <div className="auth-wrapper">
        <Auth name="discord" />
      </div>
    </>
  );
};

export default Main;
