// NavBar component
// Only loaded when the user is authenticated

import { useMemo } from "react";

const Footer = () => {
  const env: "development" | "production" = useMemo(() => {
    return import.meta.env.VITE_TEST || process.env.NODE_ENV === "test"
      ? "development"
      : "production";
  }, []);

  return env === "development" ? (
    <>
      <div className="footer-wrapper">
        <div className="dev-footer">
          <p>DEVELOPMENT</p>
        </div>
      </div>
    </>
  ) : (
    <></>
  );
};

export default Footer;
