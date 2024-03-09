// NavBar component
// Only loaded when the user is authenticated

import { StoreState } from "@/store";
import {
  CSSProperties,
  FC,
  KeyboardEventHandler,
  MouseEventHandler,
} from "react";
import { useSelector } from "react-redux";

type ProfileProps = {
  className?: string;
  onClick?: MouseEventHandler<HTMLDivElement>;
  onKeyUp?: KeyboardEventHandler<HTMLDivElement>;
  style?: CSSProperties;
};

const Profile: FC<ProfileProps> = ({ className, onClick, onKeyUp, style }) => {
  const auth = useSelector((state: StoreState) => state.auth);

  const userAvatar =
    auth.claims.provider === "discord"
      ? `https://cdn.discordapp.com/avatars/${auth.claims.provider_sub}/${auth.claims.avatar}.png?size=128`
      : `${window.location.origin}/public/images/profile.png`;
  return (
    <>
      <div
        className={className ?? "user-profile"}
        onClick={onClick ?? (() => {})}
        onKeyUp={onKeyUp ?? (() => {})}
        style={style ?? {}}
      >
        <img src={userAvatar} alt="avatar" className="profile-avatar" />
        <div className="profile-info">
          <div className="username">{auth.claims.username}</div>
          <div className="auth-provider">{auth.claims.provider}</div>
        </div>
      </div>
    </>
  );
};

export default Profile;
