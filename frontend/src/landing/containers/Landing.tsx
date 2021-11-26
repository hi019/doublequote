import { Intro } from "../components/Intro";
import { Features } from "../components/Features";

import "./landing.css";
import { Nav } from "../../nav";

export const Landing = () => {
  return (
    <div>
      <Nav />
      <div className={"h-screen"}>
        <Intro />
      </div>
      <div className={"min-h-screen"}>
        <Features />
      </div>
    </div>
  );
};
