import { Heading, Text, VStack } from "@chakra-ui/react";

export const Welcome = () => {
  return (
    <VStack
      w={"100vw"}
      h={"100vh"}
      gap={6}
      justifyContent={"center"}
      alignItems={"center"}
      textAlign={"center"}
      px={6}
    >
      <Heading fontWeight={"normal"} size={"lg"}>
        Welcome to Doublequote. We're delighted to have you.
      </Heading>
      <Text size={"2xl"} w={{ base: "100%", md: "80%", lg: "40%" }}>
        We've sent you a confirmation email. You'll need to complete it before
        you can sign in.
      </Text>
    </VStack>
  );
};
