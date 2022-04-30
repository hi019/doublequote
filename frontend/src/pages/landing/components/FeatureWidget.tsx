import { Box, Heading, Text } from "@chakra-ui/react";

interface Props {
  title: string;
  body: string;
  icon: JSX.Element;
}

export const FeatureWidget = ({ title, body, icon }: Props) => {
  return (
    <Box
      bg={"white"}
      p={8}
      w={80}
      h={60}
      boxShadow={"0px 0px 3px -0.1px rgba(0, 0, 0, 0.25)"}
    >
      <Box h={"full"}>
        <Box w={7} mb={3} textColor={"purple.500"}>
          {icon}
        </Box>

        <Heading fontSize={"xl"} fontWeight={"normal"} textColor={"gray.600"}>
          {title}
        </Heading>

        <Text mt={3} textColor={"gray.700"}>
          {body}
        </Text>
      </Box>
    </Box>
  );
};
