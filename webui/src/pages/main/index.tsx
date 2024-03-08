import React, { useEffect, useMemo } from "react";
import Auth from "@/components/auth/Auth";
import { useDispatch, useSelector } from "react-redux";
import { StoreState } from "@/store";
import { login } from "@/features/auth/authSlice";
import { useLocation } from "react-router-dom";
import NavBar from "@/components/navbar/NavBar";
import Footer from "@/components/footer/Footer";

const Main = () => {
  const dispatch = useDispatch();
  const location = useLocation();
  const auth = useSelector((state: StoreState) => state.auth);

  useEffect(() => {
    const handleAuth = async () => {
      const searchParams = new URLSearchParams(location.search);
      // clear the code from the url
      window.history.replaceState({}, document.title, location.pathname);

      const access_token = searchParams.get("a"); // a for access_token
      if (access_token) {
        console.log("access_token", access_token);
        // dispatch the token to the store
        dispatch(login({ authTokens: { access_token, refresh_token: null } }));
      }
    };

    handleAuth();
  }, [dispatch, location]);

  const env: "development" | "production" = useMemo(() => {
    return import.meta.env.VITE_TEST || process.env.NODE_ENV === "test"
      ? "development"
      : "production";
  }, []);

  console.log(auth);

  return (
    <>
      {!auth.isLoggedIn ? (
        <div className="auth-wrapper">
          <Auth name="discord" />
        </div>
      ) : (
        <>
          <NavBar />
          <div className="main-wrapper">Content goes here</div>
          <Footer />
        </>
      )}
    </>
  );
};

export default Main;
