import {
  Slide,
  Toolbar,
  Typography,
  AppBar as MuiAppBar,
  AppBarProps as MuiAppBarProps,
  useScrollTrigger,
  IconButton,
  styled,
  Box,
  useMediaQuery,
} from "@mui/material";
import { ReactElement, ReactNode, useEffect, useState } from "react";
import { SITE_NAME } from "@/consts";
import { Link } from "react-router-dom";
import NavigationDrawer, { DRAWER_WIDTH, AppNavigationHeader, isMobileQuery } from "@/components/NavigationDrawer";
import MenuIcon from "@mui/icons-material/Menu";

function HideOnScroll(props: { window?: () => Window; children: ReactElement }) {
  const { children, window } = props;
  const trigger = useScrollTrigger({
    target: window ? window() : undefined,
  });

  return (
    <Slide appear={false} direction="down" in={!trigger}>
      {children}
    </Slide>
  );
}

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const ShellAppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== "open",
})<AppBarProps>(({ theme, open }) => {
  const isMobile = useMediaQuery(isMobileQuery);

  return {
    transition: theme.transitions.create(["margin", "width"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    ...(!isMobile &&
      open && {
        width: `calc(100% - ${DRAWER_WIDTH}px)`,
        transition: theme.transitions.create(["margin", "width"], {
          easing: theme.transitions.easing.easeOut,
          duration: theme.transitions.duration.enteringScreen,
        }),
        marginRight: DRAWER_WIDTH,
      }),
  };
});

const AppMain = styled("main", { shouldForwardProp: (prop) => prop !== "open" })<{
  open?: boolean;
}>(({ theme, open }) => {
  const isMobile = useMediaQuery(isMobileQuery);

  return {
    flexGrow: 1,
    transition: theme.transitions.create("margin", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    marginRight: -DRAWER_WIDTH,
    ...(!isMobile &&
      open && {
        transition: theme.transitions.create("margin", {
          easing: theme.transitions.easing.easeOut,
          duration: theme.transitions.duration.enteringScreen,
        }),
        marginRight: 0,
      }),
    position: "relative",
  };
});

export default function AppShell({ children }: { children: ReactNode }) {
  let documentWindow: Window;

  const isMobile = useMediaQuery(isMobileQuery);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    documentWindow = window;
  });

  return (
    <>
      <HideOnScroll window={() => documentWindow}>
        <ShellAppBar open={open} position="fixed">
          <Toolbar sx={{ height: 64 }}>
            <IconButton
              size="large"
              edge="start"
              color="inherit"
              aria-label="menu"
              sx={{ ml: isMobile ? 0.5 : 0, mr: 2 }}
            >
              <img src="/favicon.svg" alt="Logo" width={32} height={32} />
            </IconButton>

            <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontSize: "1.2rem" }}>
              <Link to="/">{SITE_NAME}</Link>
            </Typography>

            <IconButton
              size="large"
              edge="start"
              color="inherit"
              aria-label="menu"
              onClick={() => setOpen(true)}
              sx={{ width: 64, mr: 1, display: !isMobile && open ? "none" : "block" }}
            >
              <MenuIcon />
            </IconButton>
          </Toolbar>
        </ShellAppBar>
      </HideOnScroll>

      <Box sx={{ display: "flex" }}>
        <AppMain open={open}>
          <AppNavigationHeader />

          {children}
        </AppMain>

        <NavigationDrawer open={open} onClose={() => setOpen(false)} />
      </Box>
    </>
  );
}
