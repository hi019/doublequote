import { Box, Heading, Image, Text, VStack } from "@chakra-ui/react";

export const Intro = () => {
  return (
    <Box
      h={"full"}
      display={"flex"}
      alignItems={"center"}
      textAlign={"center"}
      flexDir={"column"}
      justifyContent={"space-between"}
    >
      <Box pt={16} px={4} textAlign={"center"}>
        <VStack gap={6}>
          <Heading
            fontFamily={"Playfair Display, ui-serif, Georgia"}
            fontSize={"3xl"}
            w={"90%"}
            color={"gray.800"}
            fontWeight={"bold"}
          >
            Read all of your content, in one place
          </Heading>

          <Text
            w={{ base: "80%", md: "60%", lg: "30%" }}
            fontSize={"md"}
            color={"gray.800"}
          >
            Doublequote aggregates everything you read into a single
            distraction-free user interface.
          </Text>
        </VStack>

        <VStack mt={{ base: 20, lg: 12 }}>
          <Image w={{ base: "100%", md: "60%" }} src={"/ui.png"} />
        </VStack>
      </Box>

      {/*<Box w={"100%"}>*/}
      {/*  <img src={"/wave.svg"} alt={"Wave transition"} />*/}
      {/*</Box>*/}
    </Box>
  );
};
