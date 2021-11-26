import { Title } from "../../shared/components/Title";
import { motion } from "framer-motion";
import { Checkbox } from "./Checkbox";

const variants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

export const Message = () => {
  return (
    <div>
      <div className={"mb-12"}>
        <Checkbox />
      </div>

      <motion.div initial="hidden" animate="visible" variants={variants}>
        <Title>Email verified. Welcome to Doublequote!</Title>
      </motion.div>
    </div>
  );
};
