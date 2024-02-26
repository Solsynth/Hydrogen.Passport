import { Collapse, Divider, ListItemIcon, ListItemText, Menu, MenuItem, styled } from "@mui/material";
import { theme } from "@/theme";
import { Fragment, ReactNode, useState } from "react";
import HowToRegIcon from "@mui/icons-material/HowToReg";
import LoginIcon from "@mui/icons-material/Login";
import FaceIcon from "@mui/icons-material/Face";
import LogoutIcon from "@mui/icons-material/ExitToApp";
import ExpandLess from "@mui/icons-material/ExpandLess";
import ExpandMore from "@mui/icons-material/ExpandMore";
import { useUserinfo } from "@/stores/userinfo.tsx";
import { PopoverProps } from "@mui/material/Popover";
import { Link } from "react-router-dom";

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
          <MenuItem onClick={() => setOpen(!open)} sx={{ pl: 2 + (depth ?? 0) * 2, width: 180 }}>
            <ListItemIcon>{item.icon}</ListItemIcon>
            <ListItemText primary={item.title} />
            {open ? <ExpandLess /> : <ExpandMore />}
          </MenuItem>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <AppNavigationSection items={item.children} depth={(depth ?? 0) + 1} />
          </Collapse>
        </Fragment>
      );
    } else {
      return (
        <Link key={idx} to={item.link ?? "/"}>
          <MenuItem sx={{ pl: 2 + (depth ?? 0) * 2, width: 180 }}>
            <ListItemIcon>{item.icon}</ListItemIcon>
            <ListItemText primary={item.title} />
          </MenuItem>
        </Link>
      );
    }
  });
}

export function AppNavigation() {
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

  return <AppNavigationSection items={nav} />;
}

export const isMobileQuery = theme.breakpoints.down("md");

export default function NavigationMenu({ anchorEl, open, onClose }: {
  anchorEl: PopoverProps["anchorEl"];
  open: boolean;
  onClose: () => void
}) {
  return (
    <Menu anchorEl={anchorEl} open={open} onClose={onClose}>
      <AppNavigation />
    </Menu>
  );
}
