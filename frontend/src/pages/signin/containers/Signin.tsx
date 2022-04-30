import { isErrorType, isInvalidParamError } from "../../../helpers/types";
import { toast } from "react-toastify";
import React from "react";
import { Message } from "../components/Message";
import { Form, SigninForm } from "../components/Form";
import { useLocation, useNavigate } from "react-router-dom";
import { Location } from "history";
import { useSigninMutation } from "../../../api";
import { useAppDispatch } from "../../../store/hooks";
import { setIsSignedIn } from "../../../slices/user";
import { Box } from "@chakra-ui/react";

export const Signin = () => {
  const [login, { isLoading, isSuccess, error }] = useSigninMutation();
  const navigate = useNavigate();
  // @ts-ignore
  const { state } = useLocation() as Location<{ path: string }>;
  const dispatch = useAppDispatch();

  const onSubmit = (data: SigninForm) => {
    login(data);
  };

  if (isSuccess) {
    dispatch(setIsSignedIn(true));
    navigate(state?.path || "/");
    return null;
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
