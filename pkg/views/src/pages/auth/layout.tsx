import { Box } from "@mui/material";
import { Outlet } from "react-router-dom";

export default function AuthLayout() {
  return (
    <Box sx={{ height: "100vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
      <Box style={{ width: "100vw", maxWidth: "450px" }}>
        <Outlet />
      </Box>
    </Box>
  )
}