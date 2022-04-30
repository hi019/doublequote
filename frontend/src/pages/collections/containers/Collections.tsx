import { Text } from "@chakra-ui/react";
import { Link } from "react-router-dom";

const Collections = () => {
  return (
    <Text>
      <Link to={"/"} state={{ skipRedir: true }}>
        Home
      </Link>
    </Text>
  );
};

export default Collections;
