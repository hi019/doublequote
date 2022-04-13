interface Props {
  title: string;
  body: string;
  icon: JSX.Element;
}

export const FeatureWidget = ({ title, body, icon }: Props) => {
  return (
    <div
      className={"bg-white p-8 w-80 h-60"}
      style={{ boxShadow: "0px 0px 3px -0.1px rgba(0, 0, 0, 0.25)" }}
    >
      <div className={"h-full"}>
        <div className={"w-7 mb-3 text-purple-500"}>{icon}</div>

        <h2 className={"font-serif text-gray-800 text-xl"}>{title}</h2>

        <p className={"text-sm mt-3 text-gray-700"}>{body}</p>
      </div>
    </div>
  );
};
