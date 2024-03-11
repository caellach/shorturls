import AppMenu from "@/components/appmenu/AppMenu";
import Footer from "@/components/footer/Footer";
import ProfileMenu from "@/components/profile/ProfileMenu";
import { Outlet } from "react-router-dom";

const LoggedInLayout = () => {
  return (
    <>
      <div id="content-container">
        <div className="main-wrapper">
          <Outlet />
        </div>
        <Footer />
      </div>
      <AppMenu />
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
