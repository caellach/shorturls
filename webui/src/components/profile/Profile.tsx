// NavBar component
// Only loaded when the user is authenticated

import { logout } from "@/features/auth/authSlice";
import { StoreState } from "@/store";
import { useDispatch, useSelector } from "react-redux";
import {
  DropdownItem,
  DropdownMenu,
  DropdownToggle,
  UncontrolledDropdown,
} from "reactstrap";

const Profile = () => {
  const dispatch = useDispatch();
  const handleLogoutClick = () => {
    console.log("logout clicked");
    dispatch(logout());
  };

  const handleSettingsClick = () => {
    console.log("settings clicked");
  };

  const auth = useSelector((state: StoreState) => state.auth);

  const userAvatar =
    auth.claims.provider === "discord"
      ? `https://cdn.discordapp.com/avatars/${auth.claims.provider_sub}/${auth.claims.avatar}.png?size=128`
      : `${window.location.origin}/public/images/profile.png`;
  return (
    <>
      <UncontrolledDropdown inNavbar>
        <DropdownToggle tag="div" className="profile-toggle">
          <div className="profile-wrapper">
            <img src={userAvatar} alt="avatar" className="profile-avatar" />
            <div className="profile-info">
              <div className="username">{auth.claims.username}</div>
              <div className="auth-provider">{auth.claims.provider}</div>
            </div>
          </div>
        </DropdownToggle>
        <DropdownMenu>
          <DropdownItem onClick={handleSettingsClick}>Settings</DropdownItem>
          <DropdownItem divider />
          <DropdownItem onClick={handleLogoutClick}>Logout</DropdownItem>
        </DropdownMenu>
      </UncontrolledDropdown>
    </>
  );
};

export default Profile;
