import { motion } from "framer-motion";
import { Checkbox } from "./Checkbox";
import { Box, Heading } from "@chakra-ui/react";

const variants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

export const Message = () => {
  return (
    <Box>
      <Box mb={12}>
        <Checkbox />
      </Box>

      <motion.div initial="hidden" animate="visible" variants={variants}>
        <Heading size={"sm"} fontWeight={"normal"}>
          Email verified. Welcome to Doublequote!
        </Heading>
      </motion.div>
    </Box>
  );
};
