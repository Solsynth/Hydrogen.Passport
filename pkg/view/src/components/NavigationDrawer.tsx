import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import {
  Box, Collapse,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  styled,
  useMediaQuery
} from "@mui/material";
import { theme } from "@/theme";
import { Fragment, ReactNode, useState } from "react";
import HomeIcon from "@mui/icons-material/Home";
import ArticleIcon from "@mui/icons-material/Article";
import FeedIcon from "@mui/icons-material/RssFeed";
import InfoIcon from "@mui/icons-material/Info";
import GavelIcon from "@mui/icons-material/Gavel";
import PolicyIcon from "@mui/icons-material/Policy";
import SupervisedUserCircleIcon from "@mui/icons-material/SupervisedUserCircle";
import ExpandLess from "@mui/icons-material/ExpandLess";
import ExpandMore from "@mui/icons-material/ExpandMore";

export interface NavigationItem {
  icon?: ReactNode;
  title?: string;
  link?: string;
  divider?: boolean;
  children?: NavigationItem[];
}

export const DRAWER_WIDTH = 320;
export const NAVIGATION_ITEMS: NavigationItem[] = [
  { icon: <HomeIcon />, title: "首页", link: "/" },
  { icon: <ArticleIcon />, title: "博客", link: "/posts" },
  {
    icon: <InfoIcon />, title: "信息中心", children: [
      { icon: <GavelIcon />, title: "用户协议", link: "/i/user-agreement" },
      { icon: <PolicyIcon />, title: "隐私协议", link: "/i/privacy-policy" },
      { icon: <SupervisedUserCircleIcon />, title: "社区准则", link: "/i/community-guidelines" }
    ]
  },
  { divider: true },
  { icon: <FeedIcon />, title: "订阅源", link: "/feed" }
];

export const AppNavigationHeader = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  padding: theme.spacing(0, 1),
  justifyContent: "flex-start",
  height: 64,
  ...theme.mixins.toolbar
}));

export function AppNavigationSection({ items, depth }: { items: NavigationItem[], depth?: number }) {
  const [open, setOpen] = useState(false);

  return items.map((item, idx) => {
    if (item.divider) {
      return <Divider key={idx} sx={{ my: 1 }} />;
    } else if (item.children) {
      return (
        <Fragment key={idx}>
          <ListItemButton onClick={() => setOpen(!open)} sx={{ pl: 2 + (depth ?? 0) * 2 }}>
            <ListItemIcon>{item.icon}</ListItemIcon>
            <ListItemText primary={item.title} />
            {open ? <ExpandLess /> : <ExpandMore />}
          </ListItemButton>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <List component="div" disablePadding>
              <AppNavigationSection items={item.children} depth={(depth ?? 0) + 1} />
            </List>
          </Collapse>
        </Fragment>
      );
    } else {
      return (
        <a key={idx} href={item.link ?? "/"}>
          <ListItemButton sx={{ pl: 2 + (depth ?? 0) * 2 }}>
            <ListItemIcon>{item.icon}</ListItemIcon>
            <ListItemText primary={item.title} />
          </ListItemButton>
        </a>
      );
    }
  });
}

export function AppNavigation({ showClose, onClose }: { showClose?: boolean; onClose: () => void }) {
  return (
    <>
      <AppNavigationHeader>
        {showClose && (
          <IconButton onClick={onClose}>
            {theme.direction === "rtl" ? <ChevronLeftIcon /> : <ChevronRightIcon />}
          </IconButton>
        )}
      </AppNavigationHeader>
      <Divider />
      <List>
        <AppNavigationSection items={NAVIGATION_ITEMS} />
      </List>
    </>
  );
}

export const isMobileQuery = theme.breakpoints.down("md");

export default function NavigationDrawer({ open, onClose }: { open: boolean; onClose: () => void }) {
  const isMobile = useMediaQuery(isMobileQuery);

  return isMobile ? (
    <>
      <Box sx={{ flexShrink: 0, width: DRAWER_WIDTH }} />
      <Drawer
        keepMounted
        anchor="right"
        variant="temporary"
        open={open}
        onClose={onClose}
        sx={{
          "& .MuiDrawer-paper": {
            boxSizing: "border-box",
            width: DRAWER_WIDTH
          }
        }}
      >
        <AppNavigation onClose={onClose} />
      </Drawer>
    </>
  ) : (
    <Drawer
      variant="persistent"
      anchor="right"
      open={open}
      sx={{
        width: DRAWER_WIDTH,
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          width: DRAWER_WIDTH
        }
      }}
    >
      <AppNavigation showClose onClose={onClose} />
    </Drawer>
  );
}
