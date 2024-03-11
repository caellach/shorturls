import { FC, KeyboardEventHandler, MouseEventHandler, ReactNode } from "react";

type ProfileProps = {
  className?: string;
  children?: ReactNode;
  id?: string;
  onClick?: MouseEventHandler<HTMLDivElement>;
  onKeyUp?: KeyboardEventHandler<HTMLDivElement>;
};

const MenuHeader: FC<ProfileProps> = ({
  className,
  children,
  id,
  onClick,
  onKeyUp,
}) => {
  className =
    className != null && className.length > 0 ? `${className} menu-header` : "";
  onClick =
    onClick ??
    (() => {
      console.log(`onClick not implemented for BMMenuItem ${id}`);
    });

  return (
    <div id={id} className={className} onClick={onClick} onKeyUp={onKeyUp}>
      {children}
    </div>
  );
};

export default MenuHeader;
