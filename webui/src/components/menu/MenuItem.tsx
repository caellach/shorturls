import { FC, MouseEventHandler } from "react";
import { useLocation } from "react-router-dom";

type MenuItemProps = {
  active?: boolean;
  disabled?: boolean;
  className?: string;
  displayText?: string;
  id?: string;
  onClick?: MouseEventHandler<HTMLAnchorElement>;
  pagePath?: string;
};

const MenuItem: FC<MenuItemProps> = ({
  active,
  disabled,
  className,
  displayText,
  id,
  onClick,
  pagePath,
}) => {
  const location = useLocation();
  const isCurrentPage = (path: string): boolean => {
    return location.pathname === `/${path}`;
  };

  active =
    active !== undefined
      ? active
      : pagePath !== undefined
        ? isCurrentPage(pagePath)
        : false;
  className = className ?? "menu-item";
  className = active ? `${className} active` : className;
  className = disabled ? `${className} disabled` : className;

  displayText = displayText ?? "";
  id = id ? `bm-item-${id}` : "";
  onClick =
    onClick && !active && !disabled
      ? onClick
      : () => {
          console.log(`onClick disabled for BMMenuItem ${displayText}`);
        };

  return (
    <a
      id={id}
      className={className}
      // biome-ignore lint/a11y/useValidAnchor: <explanation>
      onClick={onClick}
      type="button"
    >
      {active ? ">>" : ""}
      {displayText}
      {active ? "<<" : ""}
    </a>
  );
};

export default MenuItem;
