import { FC } from "react";

type ProfileProps = {
  className?: string;
  displayText?: string;
};

const MenuFooter: FC<ProfileProps> = ({ className, displayText }) => {
  className = `${className} menu-footer`;
  displayText = displayText ?? "";

  return <div className={className}>{displayText}</div>;
};

export default MenuFooter;
