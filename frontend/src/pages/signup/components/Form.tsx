import { useForm } from "react-hook-form";
import { InvalidParam } from "../../../api/types";
import React from "react";
import useDeepCompareEffect from "use-deep-compare-effect";
import { Box, VStack, Text, Button, Input } from "@chakra-ui/react";

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
    <VStack as={"form"} noValidate gap={3} onSubmit={handleSubmit(onSubmit)}>
      <VStack w={"full"} gap={3}>
        <Box w={"full"}>
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
        </Box>

        <Box w={"full"}>
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
        </Box>
      </VStack>

      <Button w={"full"} type={"submit"} isLoading={isLoading}>
        Sign up
      </Button>
    </VStack>
  );
};
