import { FeatureWidget } from "./FeatureWidget";
import {
  SearchIcon,
  NewspaperIcon,
  FilterIcon,
  DocumentTextIcon,
  ClockIcon,
} from "@heroicons/react/outline";
import { Box, Heading } from "@chakra-ui/react";

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
    <Box bg={"purple.50"} pb={8}>
      <Box
        w={"100%"}
        display={"flex"}
        alignItems={"center"}
        pt={16}
        flexDir={"column"}
      >
        <Heading fontSize={"3xl"} fontWeight={"normal"}>
          How Doublequote helps declutter your content
        </Heading>

        <div
          className={
            "mt-20 grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
          }
        >
          {widgets.map((w) => (
            <FeatureWidget icon={w.icon} title={w.title} body={w.body} />
          ))}
        </div>
      </Box>
    </Box>
  );
};
