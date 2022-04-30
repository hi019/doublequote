import { Box, Heading, Link, Text } from "@chakra-ui/react";
import { Link as RouterLink } from "react-router-dom";

export const Message = () => {
  return (
    <Box>
      <Heading
        fontSize={"3xl"}
        textAlign={"center"}
        fontWeight={"extrabold"}
        size={"lg"}
      >
        Welcome back
      </Heading>
      <Text size={"md"} textAlign={"center"} mt={2}>
        Looking to{" "}
        <Link as={RouterLink} to="/signup" fontWeight={"medium"}>
          sign up
        </Link>
        ?
      </Text>
    </Box>
  );
};
