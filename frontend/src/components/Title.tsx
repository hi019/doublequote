import { styled } from "../stitches.config";

const StyledH1 = styled("h1", {
  fontFamily: "$serif",

  variants: {
    size: {
      sm: {
        fontSize: "$xl",
      },
      md: {
        fontSize: "$xl2",
      },
      lg: {
        fontSize: "$xl3",
      },
      xl: {
        fontSize: "$xl4",
      },
    },
    color: {
      primary: {
        color: "$gray12",
      },
      error: {
        color: "$error",
      },
    },
  },

  defaultVariants: {
    size: "md",
    color: "primary",
  },
});

export default StyledH1;
