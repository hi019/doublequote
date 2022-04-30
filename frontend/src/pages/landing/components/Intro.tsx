import { Box, Heading, Image, Skeleton, Text, VStack } from "@chakra-ui/react";

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
      <Box pt={16} px={8} textAlign={"center"}>
        <VStack gap={6}>
          <Heading
            fontFamily={"Playfair Display, ui-serif, Georgia"}
            fontSize={"3xl"}
            color={"gray.800"}
            fontWeight={"bold"}
          >
            Read all of your content, in one place
          </Heading>

          <Text fontSize={"md"} color={"gray.800"}>
            Doublequote aggregates everything you read into a single
            distraction-free user interface.
          </Text>
        </VStack>

        <VStack mt={{ base: 20, lg: 12 }}>
          <Image
            fallback={<Skeleton w={96} h={96} />}
            w={{ base: "100%", md: "60%" }}
            src={"/ui.png"}
          />
        </VStack>
      </Box>

      {/*<Box w={"100%"}>*/}
      {/*  <img src={"/wave.svg"} alt={"Wave transition"} />*/}
      {/*</Box>*/}
    </Box>
  );
};
