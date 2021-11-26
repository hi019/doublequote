import classNames from "classnames";
import React from "react";
import { Text } from "./Text";

type InputProps = JSX.IntrinsicElements["input"];

interface Props extends Omit<InputProps, "ref"> {
  error?: string;
}

export const Input = React.forwardRef<HTMLInputElement, Props>(
  ({ className, error, ...rest }: Props, ref) => {
    return (
      <div>
        <input
          className={classNames(
            (className =
              "appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"),
            { "border-red-400": error },
            className
          )}
          ref={ref}
          {...rest}
        />

        {error && (
          <Text role={"alert"} size={"sm"} className={"mt-1"}>
            {error}
          </Text>
        )}
      </div>
    );
  }
);
