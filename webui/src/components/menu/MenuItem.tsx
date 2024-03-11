import { FC, MouseEventHandler } from "react";

type ProfileProps = {
  className?: string;
  displayText?: string;
  id?: string;
  onClick?: MouseEventHandler<HTMLAnchorElement>;
};

const MenuItem: FC<ProfileProps> = ({
  className,
  displayText,
  id,
  onClick,
}) => {
  className = className ?? "menu-item";
  displayText = displayText ?? "";
  id = id ? `bm-item-${id}` : "";
  onClick =
    onClick ??
    (() => {
      console.log(`onClick not implemented for BMMenuItem ${displayText}`);
    });

  return (
    <a
      id={id}
      className={className}
      // biome-ignore lint/a11y/useValidAnchor: <explanation>
      onClick={onClick}
      type="button"
    >
      {displayText}
    </a>
  );
};

export default MenuItem;
