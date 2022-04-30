import { Message } from "../components/Message";
import { useVerifyEmailMutation } from "../../../api";
import React from "react";
import { useLocation } from "react-router-dom";
import { Heading, Spinner, VStack } from "@chakra-ui/react";

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
    element = <Heading textColor={"error"}>Error: token not found</Heading>;
  } else if (isLoading) {
    element = <Spinner />;
  } else if (isSuccess) {
    element = <Message />;
  } else if (error) {
    element = <Heading textColor={"error"}>Error</Heading>;
  }

  return (
    <VStack
      minH={"100vh"}
      alignItems={"center"}
      justifyContent={"center"}
      bg={"gray.50"}
      py={12}
      px={{ base: 4, sm: 6, lg: 8 }}
    >
      {element}
    </VStack>
  );
};
