import Title from "../../../components/Title";

export const Welcome = () => {
  return (
    <div
      className={
        "w-screen h-screen flex-col flex space-y-8 justify-center items-center text-center px-6"
      }
    >
      <Title size={"lg"}>
        Welcome to Doublequote. We're delighted to have you.
      </Title>
      <Title size={"sm"}>
        We've sent you a confirmation email. You'll need to complete it before
        you can sign in.
      </Title>
    </div>
  );
};
