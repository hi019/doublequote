import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Landing } from "./pages/landing";
import { Signup } from "./pages/signup";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Signin } from "./pages/signin";
import { RequireAuth } from "./components/RequireAuth";
import { Verify } from "./pages/verify";
import Collections from "./pages/collections";
import { Heading } from "@chakra-ui/react";

function App() {
  return (
    <>
      <ToastContainer />

      <BrowserRouter>
        <Routes>
          <Route path={"/signup"} element={<Signup />} />
          <Route path={"/signin"} element={<Signin />} />

          <Route element={<RequireAuth />}>
            <Route path={"/collections"} element={<Collections />} />
          </Route>

          <Route path={"/verify"} element={<Verify />} />
          <Route path={"/"} element={<Landing />} />

          <Route path="*" element={<Heading>Not found</Heading>} />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
