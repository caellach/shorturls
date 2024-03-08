// NavBar component
// Only loaded when the user is authenticated
import Profile from "../profile/Profile";

const NavBar = () => {
  return (
    <>
      <div className="navbar-wrapper">
        <div className="navbar-start"> </div>
        <div className="navbar-end">
          <Profile />
        </div>
      </div>
    </>
  );
};

export default NavBar;
