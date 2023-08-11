// ** Type Imports
import { PaletteMode } from "@mui/material";
import { ThemeColor } from "src/@core/layouts/types";

const DefaultPalette = (mode: PaletteMode, themeColor: ThemeColor) => {
  // ** Vars
  const lightColor = "58, 53, 65";
  const darkColor = "231, 227, 252";
  const mainColor = mode === "light" ? lightColor : darkColor;

  const primaryGradient = () => {
    if (themeColor === "primary") {
      return "#364156";
    } else if (themeColor === "secondary") {
      return "#A7C4C2";
    } else if (themeColor === "success") {
      return "#8FBC94";
    } else if (themeColor === "error") {
      return "#D16666";
    } else if (themeColor === "warning") {
      return "#FFB563";
    } else {
      return "#94C5CC";
    }
  };

  return {
    customColors: {
      main: mainColor,
      primaryGradient: primaryGradient(),
      tableHeaderBg: mode === "light" ? "#F9FAFC" : "#3D3759",
    },
    common: {
      black: "#000",
      white: "#FFF",
    },
    mode: mode,
    primary: {
      light: "#868d9a",
      main: "#364156",
      dark: "#262e3c",
      contrastText: "#FFF",
    },
    secondary: {
      light: "#cadcda",
      main: "#A7C4C2",
      dark: "#647674",
      contrastText: "#FFF",
    },
    success: {
      light: "#bcd7bf",
      main: "#8FBC94",
      dark: "#567159",
      contrastText: "#FFF",
    },
    error: {
      light: "#e3a3a3",
      main: "#D16666",
      dark: "#7d3d3d",
      contrastText: "#FFF",
    },
    warning: {
      light: "#ffd3a1",
      main: "#FFB563",
      dark: "#b37f45",
      contrastText: "#FFF",
    },
    info: {
      light: "#bfdce0",
      main: "#94C5CC",
      dark: "#59767a",
      contrastText: "#FFF",
    },
    grey: {
      50: "#FAFAFA",
      100: "#F5F5F5",
      200: "#EEEEEE",
      300: "#E0E0E0",
      400: "#BDBDBD",
      500: "#9E9E9E",
      600: "#757575",
      700: "#616161",
      800: "#424242",
      900: "#212121",
      A100: "#D5D5D5",
      A200: "#AAAAAA",
      A400: "#616161",
      A700: "#303030",
    },
    text: {
      primary: `rgba(${mainColor}, 0.87)`,
      secondary: `rgba(${mainColor}, 0.68)`,
      disabled: `rgba(${mainColor}, 0.38)`,
    },
    divider: `rgba(${mainColor}, 0.12)`,
    background: {
      paper: mode === "light" ? "#FFF" : "#312D4B",
      default: mode === "light" ? "#efefef" : "#625d79",
    },
    action: {
      active: `rgba(${mainColor}, 0.54)`,
      hover: `rgba(${mainColor}, 0.04)`,
      selected: `rgba(${mainColor}, 0.08)`,
      disabled: `rgba(${mainColor}, 0.3)`,
      disabledBackground: `rgba(${mainColor}, 0.18)`,
      focus: `rgba(${mainColor}, 0.12)`,
    },
  };
};

export default DefaultPalette;
