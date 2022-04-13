import { useForm } from "react-hook-form";
import { Button } from "../../../components/Button";
import { InvalidParam } from "../../../api/types";
import React from "react";
import useDeepCompareEffect from "use-deep-compare-effect";
import Input from "../../../components/Input";
import { inlineCss } from "../../../stitches.config";
import Div from "../../../components/Div";
import Text from "../../../components/Text";

interface Props {
  onSubmit: (data: SignupForm) => void;
  isLoading: boolean;
  serverErrors?: Array<InvalidParam>;
}

export interface SignupForm {
  email: string;
  password: string;
}

export const Form = ({ onSubmit, isLoading, serverErrors }: Props) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<SignupForm>();

  useDeepCompareEffect(() => {
    if (!serverErrors) return;

    for (const se of serverErrors) {
      // @ts-ignore
      setError(se.source, { type: "server", message: se.message });
    }
  }, [[serverErrors]]);

  return (
    <form
      noValidate
      className={inlineCss({
        display: "flex",
        gap: "$5",
        flexDirection: "column",
      })}
      onSubmit={handleSubmit(onSubmit)}
    >
      <Div css={{ display: "flex", flexDirection: "column", gap: "$5" }}>
        <div>
          <Input
            id={"email"}
            type={"email"}
            formNoValidate={true}
            placeholder="Email"
            {...register("email", {
              required: "An email is required.",
              validate: (value) => value.includes("@") || "Invalid email.",
            })}
          />
          {errors.email && (
            <Text size={"sm"} color={"error"}>
              {errors.email?.message}
            </Text>
          )}
        </div>

        <div>
          <Input
            id={"password"}
            type={"password"}
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
          {errors.password && (
            <Text size={"sm"} color={"error"}>
              {errors.password?.message}
            </Text>
          )}
        </div>
      </Div>

      <Button className={"w-full"} type={"submit"} isLoading={isLoading}>
        Sign up
      </Button>
    </form>
  );
};
