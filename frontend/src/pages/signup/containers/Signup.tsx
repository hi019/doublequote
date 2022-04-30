import { isErrorType, isInvalidParamError } from "../../../helpers/types";
import { toast } from "react-toastify";
import React from "react";
import { Message } from "../components/Message";
import { Form, SignupForm } from "../components/Form";
import { useLocation, useNavigate } from "react-router-dom";
import { Location } from "history";
import { useSignupMutation } from "../../../api";
import { useAppDispatch } from "../../../store/hooks";
import { setIsSignedIn } from "../../../slices/user";
import { Box } from "@chakra-ui/react";
import { Welcome } from "../components/Welcome";

export const Signup = () => {
  const [signup, { isLoading, isSuccess, error }] = useSignupMutation();
  const navigate = useNavigate();
  // @ts-ignore
  const { state } = useLocation() as Location<{ path: string }>;
  const dispatch = useAppDispatch();

  const onSubmit = (data: SignupForm) => {
    signup(data);
  };

  if (isSuccess) {
    return <Welcome />;
  }

  if (error !== undefined && !isInvalidParamError(error)) {
    const toString = isErrorType(error)
      ? error.title
      : "An unknown error occurred";
    toast(toString);
  }

  return (
    <Box
      h={"100vh"}
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      bg={"gray.50"}
      py={12}
      px={4}
    >
      <Box maxW={"md"} w={"full"} experimental_spaceY={10}>
        <Message />
        <Form
          isLoading={isLoading}
          onSubmit={onSubmit}
          serverErrors={isInvalidParamError(error) ? error : undefined}
        />
      </Box>
    </Box>
  );
};
