// NavBar component
// Only loaded when the user is authenticated

import { logout } from "@/features/auth/authSlice";
import { KeyboardEvent, useState } from "react";
import { useDispatch } from "react-redux";
import { scaleRotate as Menu } from "react-burger-menu";
import Profile from "./Profile";

const ProfileMenu = () => {
  const [isOpen, setIsOpen] = useState(false);

  const handleStateChange = (state: any) => {
    setIsOpen(state.isOpen);
  };

  const toggleMenu = () => {
    setIsOpen(!isOpen);
  };

  const dispatch = useDispatch();
  const handleLogoutClick = () => {
    console.log("logout clicked");
    dispatch(logout());
  };

  const handleSettingsClick = () => {
    console.log("settings clicked");
  };

  return (
    <>
      <div>
        <Profile
          onClick={toggleMenu}
          onKeyUp={(e: KeyboardEvent<HTMLDivElement>) => {
            if (e.key === "Enter" || e.key === " ") {
              toggleMenu();
            }
          }}
          style={{
            transform: isOpen ? "translateX(50vw)" : "none",
            opacity: isOpen ? 0 : 1,
            transition: isOpen
              ? "transform 0.5s ease 0s, opacity 0.3s ease 0s"
              : "transform 0.5s ease 0.2s, opacity 0.75s ease 0.3s",
          }}
        />
        <Menu
          pageWrapId="bm-outer-container"
          outerContainerId="root"
          right
          width={"300px"}
          isOpen={isOpen}
          onStateChange={handleStateChange}
          customBurgerIcon={false}
        >
          <Profile className="menu-user-profile" />
          <a id="settings" className="menu-item" onClick={handleSettingsClick}>
            Settings
          </a>
          <a id="logout" className="menu-item" onClick={handleLogoutClick}>
            Logout
          </a>
        </Menu>
      </div>
    </>
  );
};

export default ProfileMenu;
