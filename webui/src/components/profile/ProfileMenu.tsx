// Only loaded when the user is authenticated

import { logout } from "@/features/auth/authSlice";
import { KeyboardEvent, useState } from "react";
import { useDispatch } from "react-redux";
import { scaleRotate as Menu } from "react-burger-menu";
import Profile from "./Profile";
import MenuHeader from "../menu/MenuHeader";
import MenuItem from "../menu/MenuItem";
import MenuFooter from "../menu/MenuFooter";

const ProfileMenu = () => {
  const [isOpen, setIsOpen] = useState(false);

  const handleStateChange = (state: { isOpen: boolean }) => {
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
      <div className="profile-menu-wrapper">
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
          pageWrapId="content-container"
          outerContainerId="root"
          right
          isOpen={isOpen}
          onStateChange={handleStateChange}
          customBurgerIcon={false}
        >
          <MenuHeader>
            <Profile
              className="menu-user-profile"
              onClick={toggleMenu}
              onKeyUp={(e: KeyboardEvent<HTMLDivElement>) => {
                if (e.key === "Enter" || e.key === " ") {
                  toggleMenu();
                }
              }}
            />
          </MenuHeader>
          <MenuItem
            disabled
            id="settings"
            displayText="Settings"
            onClick={handleSettingsClick}
          />
          <MenuItem displayText="Logout" onClick={handleLogoutClick} />
          <MenuFooter />
        </Menu>
      </div>
    </>
  );
};

export default ProfileMenu;
