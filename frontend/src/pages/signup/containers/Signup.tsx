import { Message } from "../components/Message";
import { Form, SignupForm } from "../components/Form";
import { useSignupMutation } from "../../../api";
import { isErrorType, isInvalidParamError } from "../../../helpers/types";
import { toast } from "react-toastify";
import { Welcome } from "../components/Welcome";
import React from "react";

export const Signup = () => {
  const [signup, { isLoading, isSuccess, error }] = useSignupMutation();

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
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <Message />
        <Form
          isLoading={isLoading}
          onSubmit={onSubmit}
          serverErrors={isInvalidParamError(error) ? error : undefined}
        />
      </div>
    </div>
  );
};
