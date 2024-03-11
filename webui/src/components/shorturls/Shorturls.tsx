import { ShorturlData } from "@/types/ShorturlData";
import { FC } from "react";
import ShorturlComponent from "./Shorturl";

type ShorturlsProps = {
  data: ShorturlData[];
};

const ShorturlsComponent: FC<ShorturlsProps> = ({ data }) => {
  return (
    <div className="shorturls-list">
      <ShorturlComponent header />
      {data.map((shorturlData) => (
        <ShorturlComponent data={shorturlData} />
      ))}
    </div>
  );
};

export default ShorturlsComponent;
