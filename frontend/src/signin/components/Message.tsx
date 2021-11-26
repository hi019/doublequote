import { Title } from "../../shared/components/Title";
import { Text } from "../../shared/components/Text";

export const Message = () => {
  return (
    <div>
      <Title className={"text-center font-extrabold"} size={"lg"}>
        Welcome back
      </Title>
      <Text className="mt-2 text-center text-sm text-gray-600">
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
