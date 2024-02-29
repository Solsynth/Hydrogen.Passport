import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { Box, Tab, Tabs, useMediaQuery } from "@mui/material";
import { useEffect, useState } from "react";
import { theme } from "@/theme.ts";
import DashboardIcon from "@mui/icons-material/Dashboard";
import InboxIcon from "@mui/icons-material/Inbox";
import DrawIcon from "@mui/icons-material/Draw";
import SecurityIcon from "@mui/icons-material/Security";

export default function UserLayout() {
  const [focus, setFocus] = useState(0);

  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const locations = ["/users", "/users/notifications", "/users/personalize", "/users/security"];
  const tabs = [
    { icon: <DashboardIcon />, label: "Dashboard" },
    { icon: <InboxIcon />, label: "Notifications" },
    { icon: <DrawIcon />, label: "Personalize" },
    { icon: <SecurityIcon />, label: "Security" }
  ];

  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    const idx = locations.indexOf(location.pathname);
    setFocus(idx);
  }, []);

  function swap(idx: number) {
    navigate(locations[idx]);
    setFocus(idx);
  }

  return (
    <Box sx={{ display: "flex", flexDirection: isMobile ? "column" : "row", height: "calc(100vh - 64px)" }}>
      <Box sx={{ width: isMobile ? "100%" : 280 }}>
        <Tabs
          orientation={isMobile ? "horizontal" : "vertical"}
          variant="scrollable"
          value={focus}
          onChange={(_, val) => swap(val)}
          sx={{
            borderRight: isMobile ? 0 : 1,
            borderBottom: isMobile ? 1 : 0,
            borderColor: "divider",
            height: isMobile ? "fit-content" : "100%",
            py: isMobile ? 0 : 1,
            px: isMobile ? 1 : 0
          }}
        >
          {tabs.map((tab, idx) => (
            <Tab key={idx} icon={tab.icon} iconPosition={isMobile ? "top" : "start"} label={tab.label}
                 sx={{ px: 5, justifyContent: isMobile ? "center" : "left" }} />
          ))}
        </Tabs>
      </Box>

      <Box sx={{ flexGrow: 1 }}>
        <Outlet />
      </Box>
    </Box>
  );
}