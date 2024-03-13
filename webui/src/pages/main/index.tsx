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
      const accessToken = searchParams.get("a"); // a for accessToken
      if (accessToken) {
        console.log("accessToken", accessToken);
        // dispatch the token to the store
        console.log("auth", auth);
        dispatch(login({ authTokens: { accessToken, refreshToken: null } }));
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
