import Footer from "@/components/footer/Footer";
import ProfileMenu from "@/components/profile/ProfileMenu";
import { Outlet } from "react-router-dom";

const LoggedInLayout = () => {
  return (
    <>
      <div id="bm-outer-container">
        <div className="main-wrapper">
          <Outlet />
        </div>
        <Footer />
      </div>
      <ProfileMenu />
    </>
  );
};

const AnonLayout = () => {
  return (
    <div>
      <Outlet />
    </div>
  );
};

export default { LoggedInLayout, AnonLayout };
