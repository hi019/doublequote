import { Box, Button, Flex, Image } from "@chakra-ui/react";
import { Link } from "react-router-dom";
const Nav = () => {
  return (
    <Box
      // bg={"purple.50"}
      h={14}
      px={24}
      pt={8}
      display={"flex"}
      alignItems={"center"}
      justifyContent={"space-between"}
    >
      <Image
        boxSize={"35px"}
        src={"https://tailwindui.com/img/logos/workflow-mark-indigo-500.svg"}
      />

      <Flex gap={4}>
        <Button variant={"ghost"}>
          <Link to={"/signin"}>Sign in</Link>
        </Button>
        <Button variant={"light"}>
          <Link to={"/signup"}>Sign up</Link>
        </Button>
      </Flex>
    </Box>
  );
};

export default Nav;
