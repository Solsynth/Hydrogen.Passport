import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import {
  Box,
  Collapse,
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
import HowToRegIcon from "@mui/icons-material/HowToReg";
import LoginIcon from "@mui/icons-material/Login";
import FaceIcon from "@mui/icons-material/Face";
import LogoutIcon from "@mui/icons-material/Logout";
import ExpandLess from "@mui/icons-material/ExpandLess";
import ExpandMore from "@mui/icons-material/ExpandMore";
import { useUserinfo } from "@/stores/userinfo.tsx";

export interface NavigationItem {
  icon?: ReactNode;
  title?: string;
  link?: string;
  divider?: boolean;
  children?: NavigationItem[];
}

export const DRAWER_WIDTH = 320;

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
  const { checkLoggedIn } = useUserinfo();

  const nav: NavigationItem[] = [
    ...(
      checkLoggedIn() ?
        [
          { icon: <FaceIcon />, title: "Account", link: "/users" },
          { divider: true },
          { icon: <LogoutIcon />, title: "Sign out", link: "/auth/sign-out" }
        ] :
        [
          { icon: <HowToRegIcon />, title: "Sign up", link: "/auth/sign-up" },
          { icon: <LoginIcon />, title: "Sign in", link: "/auth/sign-in" }
        ]
    )
  ];

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
        <AppNavigationSection items={nav} />
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
