import { Avatar, Button, Card, CardContent, Typography } from "@mui/material";
import { useUserinfo } from "@/stores/userinfo.tsx";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";

export function Component() {
  const { clearUserinfo } = useUserinfo();

  const navigate = useNavigate();

  async function signout() {
    clearUserinfo();
    navigate("/");
  }

  return (
    <>
      <Card variant="outlined">
        <CardContent
          style={{ padding: "40px 48px 36px" }}
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center"
          }}
        >
          <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
            <LogoutIcon />
          </Avatar>

          <Typography gutterBottom variant="h5" component="h1">Sign out</Typography>
          <Typography variant="body1">
            Sign out will clear your data on this device. Also will affected those use union identification services.
            You need sign in again get access them.
          </Typography>

          <Button
            fullWidth
            variant="contained"
            color="secondary"
            sx={{ mt: 3 }}
            onClick={() => signout()}
          >
            Sign out
          </Button>
        </CardContent>
      </Card>
    </>
  );
}