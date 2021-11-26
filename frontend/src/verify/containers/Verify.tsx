import { Message } from "../components/Message";
import { useVerifyEmailMutation } from "../../api";
import React from "react";
import { useLocation } from "react-router-dom";
import { Title } from "../../shared/components/Title";
import { Loader } from "../../shared/components/Loader";

export const Verify = () => {
  const [verifyEmail, { isLoading, isSuccess, error }] =
    useVerifyEmailMutation();

  const location = useLocation();

  console.log(location);

  const token = location.search;
  React.useEffect(() => {
    token && verifyEmail({ token });
  }, [token, verifyEmail]);

  let element: React.ReactNode;

  if (!token) {
    element = <Title>Error: token not found</Title>;
  } else if (isLoading) {
    element = <Loader className={"text-black"} />;
  } else if (isSuccess) {
    element = <Message />;
  } else if (error) {
    element = <Title>Error</Title>;
  }

  return (
    <div className="min-h-screen flex items-center flex-col justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      {element}
    </div>
  );
};
