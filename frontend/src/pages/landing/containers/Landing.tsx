import { Intro } from "../components/Intro";
import { Features } from "../components/Features";

import "./landing.css";
import Nav from "../../../components/Nav";
import { Box } from "@chakra-ui/react";
import { useAppSelector } from "../../../store/hooks";
import { Navigate, useLocation } from "react-router-dom";
import React from "react";
import { Location } from "history";

export const Landing = () => {
  // @ts-ignore
  const { state } = useLocation() as Location<{ skipRedir: boolean }>;
  const isSignedIn = useAppSelector((s) => s.user.isSignedIn);

  if (isSignedIn && !state?.skipRedir) {
    return <Navigate to={"/collections"} />;
  }

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
