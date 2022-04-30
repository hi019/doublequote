import { Box, Button, Flex, Image } from "@chakra-ui/react";
import { Link } from "react-router-dom";
import { useAppSelector } from "../../store/hooks";

const Nav = () => {
  const isSignedIn = useAppSelector((s) => s.user.isSignedIn);

  return (
    <Box
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
        {!isSignedIn && (
          <>
            <Button as={Link} to={"/signin"} variant={"ghost"}>
              Sign in
            </Button>
            <Button as={Link} to={"/signup"} variant={"light"}>
              Sign up
            </Button>
          </>
        )}
        {isSignedIn && (
          <Button as={Link} to={"/collections"} variant={"ghost"}>
            App
          </Button>
        )}
      </Flex>
    </Box>
  );
};

export default Nav;
