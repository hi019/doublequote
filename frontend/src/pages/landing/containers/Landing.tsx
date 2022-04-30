import { Intro } from "../components/Intro";
import { Features } from "../components/Features";

import "./landing.css";
import Nav from "../../../components/Nav";
import { Box } from "@chakra-ui/react";

export const Landing = () => {
  return (
    <Box>
      <Nav />
      <Box h={"100vh"}>
        <Intro />
      </Box>
      <Box minH={"100vh"}>
        <Features />
      </Box>
    </Box>
  );
};
