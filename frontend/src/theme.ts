import { ComponentStyleConfig, extendTheme } from "@chakra-ui/react";

// @ts-ignore
import resolveConfig from "tailwindcss/resolveConfig";

const tailwind = resolveConfig({});

const Button: ComponentStyleConfig = {
  baseStyle: ({ colorMode }) => ({
    fontWeight: "normal",
    rounded: "md",
  }),
  variants: {
    ghost: {
      bg: "transparent",
    },
    light: {
      bg: "indigo.100",
      textColor: "purple.700",
      _hover: {
        bg: "indigo.200",
      },
      _active: {
        bg: "indigo.200",
      },
    },
    solid: {
      bg: "indigo.600",
      textColor: "white",
      _hover: {
        bg: "indigo.500",
      },
    },
  },
};

const Link: ComponentStyleConfig = {
  baseStyle: ({ colorMode }) => ({
    textColor: "indigo.600",
    _hover: {
      textColor: "indigo.500",
    },
  }),
  variants: {
    ghost: {
      _hover: {
        textDecoration: "none",
      },
    },
  },
};

const theme = {
  components: {
    Button: Button,
    Link: Link,
  },

  colors: {
    ...tailwind.theme.colors,
    error: "tomato.600",
  },

  shadows: {
    outline: "0",
  },

  fonts: {
    heading: 'ui-serif, Georgia, Cambria, "Times New Roman", Times, serif', // Playfair Display, ui-serif, Georgia
  },
};

export default extendTheme(theme);
