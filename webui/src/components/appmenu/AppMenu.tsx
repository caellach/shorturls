// Only loaded when the user is authenticated

import { KeyboardEvent, useState } from "react";
import { scaleRotate as Menu } from "react-burger-menu";
import MenuItem from "../menu/MenuItem";
import MenuFooter from "../menu/MenuFooter";
import MenuHeader from "../menu/MenuHeader";

const AppMenu = () => {
  const [isOpen, setIsOpen] = useState(false);

  const handleStateChange = (state: { isOpen: boolean }) => {
    setIsOpen(state.isOpen);
  };

  const toggleMenu = () => {
    setIsOpen(!isOpen);
  };

  const handleSettingsClick = () => {
    // navigate
  };

  return (
    <>
      <div
        className="app-menu-wrapper"
        // set z-index:1 if isOpen
        style={{
          zIndex: isOpen ? 1 : 0,
        }}
      >
        <Menu
          pageWrapId="content-container"
          outerContainerId="root"
          isOpen={isOpen}
          onStateChange={handleStateChange}
          //customBurgerIcon={false}
        >
          <MenuHeader
            onClick={toggleMenu}
            onKeyUp={(e: KeyboardEvent<HTMLDivElement>) => {
              if (e.key === "Enter" || e.key === " ") {
                toggleMenu();
              }
            }}
          >
            Apps
          </MenuHeader>
          <MenuItem
            pagePath="shorturls"
            id="menu-shorturls"
            displayText="Short URLs"
            onClick={handleSettingsClick}
          />
          <MenuFooter />
        </Menu>
      </div>
    </>
  );
};

export default AppMenu;
