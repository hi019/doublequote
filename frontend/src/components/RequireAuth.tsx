import React from "react";
import { Outlet, useLocation } from "react-router";
import { Navigate } from "react-router-dom";
import { useAppSelector } from "../store/hooks";
import { Spinner } from "@chakra-ui/react";

export const RequireAuth = () => {
  const isSignedIn = useAppSelector((s) => s.user.isSignedIn);
  const { pathname } = useLocation();

  if (isSignedIn === null) {
    return <Spinner />;
  }

  if (!isSignedIn) {
    return <Navigate to={"/signin"} replace state={{ path: pathname }} />;
  }

  return <Outlet />;
};
