import { Size } from "../types";
import classNames from "classnames";
import { Loader } from "./Loader";

type Scheme = "primary" | "secondary";

type ButtonProps = JSX.IntrinsicElements["button"];

interface Props extends ButtonProps {
  isLoading?: boolean;
  onClick?: () => void;
  scheme?: Scheme;
  size?: Size;
}

export const Button = ({
  size,
  scheme,
  onClick,
  className,
  isLoading,
  children,
  ...rest
}: Props) => {
  size = size || "md";
  scheme = scheme || "primary";

  const classes = {
    // Size
    "py-2 px-3": size === "sm",
    "py-2 px-4 text-md": size === "md",
    "py-3 px-5 text-lg": size === "lg",

    // Type
    "bg-indigo-100 text-purple-700 hover:bg-indigo-200": scheme === "secondary",
    "bg-indigo-600 text-white hover:bg-indigo-500": scheme === "primary",
  };

  return (
    <button
      className={classNames(
        "flex justify-center text-sm rounded-sm text-gray-100 font-medium",
        classes,
        className
      )}
      onClick={onClick}
      {...rest}
    >
      {isLoading && <Loader />}
      {children}
    </button>
  );
};
