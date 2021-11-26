import { FeatureWidget } from "./FeatureWidget";
import {
  SearchIcon,
  NewspaperIcon,
  FilterIcon,
  DocumentTextIcon,
  ClockIcon,
} from "@heroicons/react/outline";

const widgets = [
  {
    title: "Keep up with the news",
    body: "Make the news come to you. With Doublequote, all of your articles are always ready for you, in one place.",
    icon: <NewspaperIcon />,
  },
  {
    title: "Find what you need",
    body: "Instantly find any article with Doublequote’s Intelligent Search. Query by title, content, publication, and more.",
    icon: <SearchIcon />,
  },
  {
    title: "Cut through the fluff",
    body: "Read what you want, remove what you won’t. Highlight or filter out specific articles based on words, topics or phrases.",
    icon: <FilterIcon />,
  },
  {
    title: "Eliminate distractions",
    body: "Use Article View to effortlessly remove popups and ads, so that you can focus on the content.",
    icon: <DocumentTextIcon />,
  },
  {
    title: "Save time",
    body: "Doublequote automatically summarizes articles, enabling you to be in the know when you're strapped for time.",
    icon: <ClockIcon />,
  },
  {
    title: "Find what you need",
    body: "Instantly find any article with Doublequote’s Intelligent Search. Query by title, content, publication, and more.",
    icon: <SearchIcon />,
  },
];

export const Features = () => {
  return (
    <div className={"bg-purple-50 pb-8"}>
      <div className={"w-100 flex items-center pt-16 flex-col"}>
        <h1 className={"text-3xl font-serif text-center text-gray-800"}>
          How Doublequote helps declutter your content
        </h1>

        <div
          className={
            "mt-20 grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
          }
        >
          {widgets.map((w) => (
            <FeatureWidget icon={w.icon} title={w.title} body={w.body} />
          ))}
        </div>
      </div>
    </div>
  );
};
