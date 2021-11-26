import { Size } from "../../types";
import classNames from "classnames";

type H1Props = JSX.IntrinsicElements["h1"];

type Font = "sans" | "serif";

interface Props extends H1Props {
  font?: Font;
  size?: Size;
}

export const Title = ({ size, children, font, className, ...rest }: Props) => {
  size = size || "md";
  font = font || "sans";

  let sizeClass;
  switch (size) {
    case "sm":
      sizeClass = "text-xl";
      break;
    case "md":
      sizeClass = "text-2xl";
      break;
    case "lg":
      sizeClass = "text-3xl";
      break;
  }

  return (
    <h1 className={classNames("font-" + font, sizeClass, className)} {...rest}>
      {children}
    </h1>
  );
};
