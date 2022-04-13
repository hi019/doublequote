import Title from "../../../components/Title";
import Text from "../../../components/Text";

export const Message = () => {
  return (
    <div>
      <Title className={"text-center font-extrabold"} size={"lg"}>
        Let's get you started
      </Title>
      <Text css={{ textAlign: "center", marginTop: "$2" }}>
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
