import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Landing } from "./pages/landing";
import { Signup } from "./pages/signup";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Signin } from "./pages/signin";
import { RequireAuth } from "./components/RequireAuth";
import { Protected } from "./components/Protected";
import { Verify } from "./pages/verify";

function App() {
  return (
    <>
      <ToastContainer />

      <BrowserRouter>
        <Routes>
          <Route path={"/signup"} element={<Signup />} />
          <Route path={"/signin"} element={<Signin />} />

          <Route element={<RequireAuth />}>
            <Route path={"/protected"} element={<Protected />} />
          </Route>

          <Route path={"/verify"} element={<Verify />} />
          <Route path={"/"} element={<Landing />} />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
