import React from "react";
import { styled } from "../stitches.config";
import * as Stitches from "@stitches/react";

const StyledInput = styled("input", {
  padding: "$2 $3",
  fontSize: "$sm",
  borderRadius: "$sm",
  border: "1px solid $gray7",
  variants: {
    width: {
      full: {
        width: "100%",
      },
    },
  },
  defaultVariants: {
    width: "full",
  },
});

export default StyledInput;
