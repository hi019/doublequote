import { Size } from "../../types";
import classNames from "classnames";

type PProps = JSX.IntrinsicElements["p"];

interface Props extends PProps {
  size?: Size;
}

export const Text = ({ children, size, className, ...rest }: Props) => {
  size = size || "md";

  const cls = "text-gray-900";

  return (
    <p className={classNames("text-" + size, cls, className)} {...rest}>
      {children}
    </p>
  );
};
