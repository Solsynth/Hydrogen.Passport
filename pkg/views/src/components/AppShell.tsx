import {
  AppBar,
  Avatar,
  Box,
  IconButton,
  Slide,
  Toolbar,
  Typography,
  useMediaQuery,
  useScrollTrigger
} from "@mui/material";
import { ReactElement, ReactNode, useEffect, useRef, useState } from "react";
import { SITE_NAME } from "@/consts";
import { Link } from "react-router-dom";
import NavigationMenu, { AppNavigationHeader, isMobileQuery } from "@/components/NavigationMenu.tsx";
import AccountCircleIcon from "@mui/icons-material/AccountCircleOutlined";
import { useUserinfo } from "@/stores/userinfo.tsx";

function HideOnScroll(props: { window?: () => Window; children: ReactElement }) {
  const { children, window } = props;
  const trigger = useScrollTrigger({
    target: window ? window() : undefined
  });

  return (
    <Slide appear={false} direction="down" in={!trigger}>
      {children}
    </Slide>
  );
}

export default function AppShell({ children }: { children: ReactNode }) {
  let documentWindow: Window;

  const { userinfo } = useUserinfo();

  const isMobile = useMediaQuery(isMobileQuery);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    documentWindow = window;
  }, []);

  const container = useRef<HTMLDivElement>(null);

  return (
    <>
      <HideOnScroll window={() => documentWindow}>
        <AppBar position="fixed">
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
              sx={{ mr: 1 }}
            >
              <Avatar
                sx={{ width: 32, height: 32, bgcolor: "transparent" }}
                ref={container}
                alt={userinfo?.displayName}
                src={userinfo?.profiles?.avatar}
              >
                <AccountCircleIcon />
              </Avatar>
            </IconButton>
          </Toolbar>
        </AppBar>
      </HideOnScroll>

      <Box component="main">
        <AppNavigationHeader />

        {children}
      </Box>

      <NavigationMenu anchorEl={container.current} open={open} onClose={() => setOpen(false)} />
    </>
  );
}
