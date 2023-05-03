// ** Icon imports
import HomeOutline from "mdi-material-ui/HomeOutline";
import AlertCircleOutline from "mdi-material-ui/AlertCircleOutline";
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
    // {
    //   title: "Feed Settings",
    //   icon: CogOutline,
    //   path: "/feed-settings",
    // },
    // {
    //   title: "Feed Details",
    //   icon: AlertCircleOutline,
    //   path: "/feed-details",
    // },
  ];
};

export default navigation;
