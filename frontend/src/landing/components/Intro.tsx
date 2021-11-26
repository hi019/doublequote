export const Intro = () => {
  return (
    <div className={"bg-yellow-50 h-full"}>
      <div className={"w-100 flex items-center text-center pt-16 flex-col"}>
        <h1 className={"text-3xl font-serif text-gray-800 font-bold"}>
          Read all of your content, in one place
        </h1>

        <p className={"pt-8 w-4/5 md:w-1/3 text-center text-md text-gray-800"}>
          Doublequote aggregates everything you read into a single
          distraction-free user interface.
        </p>

        <div className={"w-3/5 h-96 bg-gray-200 mt-20 rounded-sm"} />
      </div>

      <div className={"pt-6"}>
        <img src={"/wave.svg"} alt={"Wave transition"} />
      </div>
    </div>
  );
};
