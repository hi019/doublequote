import { Title } from "../../shared/components/Title";
import { Text } from "../../shared/components/Text";

export const Message = () => {
  return (
    <div>
      <Title className={"text-center font-extrabold"} size={"lg"}>
        Let's get you started
      </Title>
      <Text className="mt-2 text-center text-sm text-gray-600">
        Or{" "}
        <a
          href="/signin"
          className="font-medium text-indigo-600 hover:text-indigo-500"
        >
          sign in
        </a>
      </Text>
    </div>
  );
};
