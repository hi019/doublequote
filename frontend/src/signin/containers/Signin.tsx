import { isErrorType, isErrorWithInvalidParamsType } from "../../utils";
import { toast } from "react-toastify";
import React from "react";
import { Message } from "../components/Message";
import { Form, SigninForm } from "../components/Form";
import { useLocation, useNavigate } from "react-router-dom";
import { Location } from "history";
import { useSigninMutation } from "../../api";
import { useAppDispatch } from "../../store/hooks";
import { setIsSignedIn } from "../../shared/slices/user";

export const Signin = () => {
  const [login, { isLoading, isSuccess, error }] = useSigninMutation();
  const navigate = useNavigate();
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

  if (error !== undefined && !isErrorWithInvalidParamsType(error)) {
    const toString = isErrorType(error)
      ? error.title
      : "An unknown error occurred";
    toast(toString);
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <Message />
        <Form
          isLoading={isLoading}
          onSubmit={onSubmit}
          serverErrors={isErrorWithInvalidParamsType(error) ? error : undefined}
        />
      </div>
    </div>
  );
};
