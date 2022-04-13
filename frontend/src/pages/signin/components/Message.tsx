import Title from "../../../components/Title";
import Text from "../../../components/Text";
import { inlineCss } from "../../../stitches.config";

export const Message = () => {
  return (
    <div>
      <Title className={"text-center font-extrabold"} size={"lg"}>
        Welcome back
      </Title>
      <Text size={"md"} css={{ textAlign: "center", marginTop: "$2" }}>
        Looking to{" "}
        <a
          href="/signup"
          className="font-medium text-indigo-600 hover:text-indigo-500"
        >
          sign up
        </a>
        ?
      </Text>
    </div>
  );
};
