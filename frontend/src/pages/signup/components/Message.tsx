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
        Let's get you started
      </Heading>
      <Text size={"md"} textAlign={"center"} mt={2}>
        Or{" "}
        <Link as={RouterLink} to="/signin" fontWeight={"medium"}>
          sign in
        </Link>
      </Text>
    </Box>
  );
};
