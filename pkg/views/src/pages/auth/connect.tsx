import { useEffect, useState } from "react";
import {
  Alert,
  Avatar,
  Box,
  Button,
  Card,
  CardContent,
  Collapse,
  Grid,
  LinearProgress,
  Typography
} from "@mui/material";
import { request } from "@/scripts/request.ts";
import { useUserinfo } from "@/stores/userinfo.tsx";
import { useSearchParams } from "react-router-dom";
import OutletIcon from "@mui/icons-material/Outlet";
import WhatshotIcon from "@mui/icons-material/Whatshot";

export function Component() {
  const { getAtk } = useUserinfo();

  const [panel, setPanel] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const [client, setClient] = useState<any>(null);

  const [searchParams] = useSearchParams();

  async function preconnect() {
    const res = await request(`/api/auth/o/connect${location.search}`, {
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });

    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();

      if (data["session"]) {
        setPanel(1);
        redirect(data["session"]);
      } else {
        setClient(data["client"]);
        setLoading(false);
      }
    }
  }

  useEffect(() => {
    preconnect().then(() => console.log("Fetched metadata"));
  }, []);

  function decline() {
    if (window.history.length > 0) {
      window.history.back();
    } else {
      window.close();
    }
  }

  async function approve() {
    setLoading(true);

    const res = await request("/api/auth/o/connect?" + new URLSearchParams({
      client_id: searchParams.get("client_id") as string,
      redirect_uri: encodeURIComponent(searchParams.get("redirect_uri") as string),
      response_type: "code",
      scope: searchParams.get("scope") as string
    }), {
      method: "POST",
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });

    if (res.status !== 200) {
      setError(await res.text());
      setLoading(false);
    } else {
      const data = await res.json();
      setPanel(1);
      setTimeout(() => redirect(data["session"]), 1850);
    }
  }

  function redirect(session: any) {
    const url = `${searchParams.get("redirect_uri")}?code=${session["grant_token"]}&state=${searchParams.get("state")}`;
    window.open(url, "_self");
  }

  const elements = [
    (
      <>
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <OutletIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign in to {client?.name}
        </Typography>
        <Box sx={{ mt: 3, width: "100%" }}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Typography fontWeight="bold">About this app</Typography>
              <Typography variant="body2">{client?.description}</Typography>
            </Grid>
            <Grid item xs={12}>
              <Typography fontWeight="bold">Make you trust this app</Typography>
              <Typography variant="body2">
                After you click Approve button, you will share your basic personal information to this application
                developer. Some of them will leak your data. Think twice.
              </Typography>
            </Grid>
            <Grid item xs={12} md={6}>
              <Button
                fullWidth
                color="info"
                variant="outlined"
                disabled={loading}
                sx={{ mt: 3 }}
                onClick={() => decline()}
              >
                Decline
              </Button>
            </Grid>
            <Grid item xs={12} md={6}>
              <Button
                fullWidth
                variant="outlined"
                disabled={loading}
                sx={{ mt: 3 }}
                onClick={() => approve()}
              >
                Approve
              </Button>
            </Grid>
          </Grid>
        </Box>
      </>
    ),
    (
      <>
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <WhatshotIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Authorized
        </Typography>
        <Box sx={{ mt: 3, width: "100%", textAlign: "center" }}>
          <Grid container spacing={2}>
            <Grid item xs={12} sx={{ my: 8 }}>
              <Typography variant="h6">Now Redirecting...</Typography>
              <Typography>Hold on a second, we are going to redirect you to the target.</Typography>
            </Grid>
          </Grid>
        </Box>
      </>
    )
  ];

  return (
    <>
      {error && <Alert severity="error" className="capitalize" sx={{ mb: 2 }}>{error}</Alert>}

      <Card variant="outlined">
        <Collapse in={loading}>
          <LinearProgress />
        </Collapse>

        <CardContent
          style={{ padding: "40px 48px 36px" }}
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center"
          }}
        >
          {elements[panel]}
        </CardContent>
      </Card>
    </>
  );
}