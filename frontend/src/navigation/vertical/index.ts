// ** Icon imports
import HomeOutline from "mdi-material-ui/HomeOutline";
import CogOutline from "mdi-material-ui/CogOutline";

// ** Type import
import { VerticalNavItemsType } from "src/@core/layouts/types";

const navigation = (): VerticalNavItemsType => {
  return [
    {
      title: "Monitor",
      icon: HomeOutline,
      path: "/",
    },
    {
      title: "Resources",
      icon: CogOutline,
      path: "/resources-grid",
    },
  ];
};

export default navigation;
