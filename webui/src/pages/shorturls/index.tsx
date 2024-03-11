import { ShorturlData } from "@/types/ShorturlData";
import ShorturlsComponent from "@/components/shorturls/Shorturls";

const Shorturls = () => {
  const shorturlData: ShorturlData[] = [
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
    {
      id: "CorrectHorseBatteryStaple",
      shortUrl: "http://localhost:3000/CorrectHorseBatteryStaple",
      url: "https://www.google.com",
      uses: Math.floor(Math.random() * 1000),
      lastUsed: new Date(),
    },
  ];

  return (
    <>
      <div className="shorturls-content-wrapper">
        <h1>Shorturls</h1>
        <ShorturlsComponent data={shorturlData} />
      </div>
    </>
  );
};

export default Shorturls;
