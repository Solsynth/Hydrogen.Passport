import { useEffect } from "react";
import { Box, CircularProgress } from "@mui/material";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { useUserinfo } from "@/stores/userinfo.tsx";

export default function AuthGuard() {
  const { userinfo } = useUserinfo();

  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    console.log(userinfo)
    if (userinfo?.isReady) {
      if (!userinfo?.isLoggedIn) {
        const callback = location.pathname + location.search;
        navigate({ pathname: "/auth/sign-in", search: `redirect_uri=${callback}` });
      }
    }
  }, [userinfo]);

  return !userinfo?.isReady ? (
    <Box sx={{ pt: 32, display: "flex", justifyContent: "center", alignItems: "center" }}>
      <Box>
        <CircularProgress />
      </Box>
    </Box>
  ) : <Outlet />;
}