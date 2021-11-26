import { useForm } from "react-hook-form";
import { Input } from "../../shared/components/Input";
import { Button } from "../../shared/components/Button";
import { InvalidParam } from "../../api/types";
import React from "react";
import useDeepCompareEffect from "use-deep-compare-effect";

interface Props {
  onSubmit: (data: SigninForm) => void;
  isLoading: boolean;
  serverErrors?: Array<InvalidParam>;
}

export interface SigninForm {
  email: string;
  password: string;
}

export const Form = ({ onSubmit, isLoading, serverErrors }: Props) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<SigninForm>();

  useDeepCompareEffect(() => {
    if (!serverErrors) return;

    for (const se of serverErrors) {
      // @ts-ignore
      setError(se.source, { type: "server", message: se.message });
    }
  }, [[serverErrors]]);

  return (
    <form noValidate className={"space-y-5"} onSubmit={handleSubmit(onSubmit)}>
      <div className={"space-y-3"}>
        <Input
          id={"email"}
          type={"email"}
          error={errors.email?.message}
          formNoValidate={true}
          placeholder="Email"
          {...register("email", {
            required: "An email is required.",
            validate: (value) => value.includes("@") || "Invalid email.",
          })}
        />

        <Input
          id={"password"}
          type={"password"}
          error={errors.password?.message}
          placeholder="Password"
          {...register("password", {
            required: "A password is required.",
            minLength: {
              value: 6,
              message: "Your password must be at least 6 characters.",
            },
            maxLength: {
              value: 64,
              message: "Your password must be at most 64 characters.",
            },
          })}
        />
      </div>

      <Button className={"w-full"} type={"submit"} isLoading={isLoading}>
        Sign in
      </Button>
    </form>
  );
};
