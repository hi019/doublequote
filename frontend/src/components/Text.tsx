import { styled } from "../stitches.config";

const StyledText = styled("p", {
  color: "$gray12",
  variants: {
    size: {
      sm: {
        fontSize: "$sm",
      },
      md: {
        fontSize: "$md",
      },
      lg: {
        fontSize: "$lg",
      },
      xl: {
        fontSize: "$xl",
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

export default StyledText;
