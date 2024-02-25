import UserIcon from "@mui/icons-material/PersonAddAlt1";
import HowToRegIcon from "@mui/icons-material/HowToReg";
import { Link as RouterLink, useNavigate, useSearchParams } from "react-router-dom";
import {
  Alert,
  Avatar,
  Box,
  Button,
  Card,
  CardContent,
  Checkbox,
  Collapse,
  FormControlLabel,
  Grid,
  LinearProgress,
  Link,
  TextField,
  Typography
} from "@mui/material";
import { FormEvent, useState } from "react";
import { request } from "@/scripts/request.ts";
import { useWellKnown } from "@/stores/wellKnown.tsx";

export default function SignUpPage() {
  const [done, setDone] = useState(false);

  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const { wellKnown } = useWellKnown();

  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  async function submit(evt: FormEvent<HTMLFormElement>) {
    evt.preventDefault();

    const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
    if (!data.human_verification) return;
    if (!data.name || !data.nick || !data.email || !data.password) return;

    setLoading(true);
    const res = await request("/api/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data)
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setError(null);
      setDone(true);
    }
    setLoading(false);
  }

  function callback() {
    if (searchParams.has("closable")) {
      window.close();
    } else {
      navigate("/auth/sign-in");
    }
  }

  const elements = [
    (
      <>
        <Avatar sx={{ mb: 1, bgcolor: "secondary.main" }}>
          <UserIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Create an account
        </Typography>
        <Box component="form" onSubmit={submit} sx={{ mt: 3, width: "100%" }}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                name="name"
                required
                fullWidth
                label="Username"
                autoComplete="username"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                name="nick"
                required
                fullWidth
                label="Nickname"
                autoComplete="nickname"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                autoComplete="email"
                name="email"
                required
                fullWidth
                label="Email Address"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Password"
                name="password"
                required
                fullWidth
                type="password"
                autoComplete="new-password"
              />
            </Grid>
            {
              !wellKnown?.open_registration && <Grid item xs={12}>
                <TextField
                  label="Magic Token"
                  name="magic_token"
                  required
                  fullWidth
                  type="password"
                  autoComplete="magic-token"
                  helperText={"This server uses invitations only."}
                />
              </Grid>
            }
            <Grid item xs={12}>
              <FormControlLabel
                name="human_verification"
                control={<Checkbox value="allowExtraEmails" color="primary" />}
                label={"I'm not a robot."}
              />
            </Grid>
          </Grid>
          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={loading}
            sx={{ mt: 3, mb: 2 }}
          >
            {loading ? "Signing Now..." : "Sign Up"}
          </Button>
        </Box>
      </>
    ),
    (
      <>
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <HowToRegIcon />
        </Avatar>

        <Typography gutterBottom variant="h5" component="h1">Congratulations!</Typography>
        <Typography variant="body1">
          Your account has been created and activation email has sent to your inbox!
        </Typography>

        <Typography sx={{ my: 2 }}>
          <Link onClick={() => callback()} className="cursor-pointer">Go login</Link>
        </Typography>

        <Typography variant="body2">
          After you login, then you can take part in the entire smartsheep community.
        </Typography>
      </>
    )
  ];

  return (
    <Box sx={{ height: "100vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
      <Box style={{ width: "100vw", maxWidth: "450px" }}>
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

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
            {!done ? elements[0] : elements[1]}
          </CardContent>
        </Card>

        <Grid container justifyContent="center" sx={{ mt: 2 }}>
          <Grid item>
            <Link component={RouterLink} to="/auth/sign-in" variant="body2">
              Already have an account? Sign in!
            </Link>
          </Grid>
        </Grid>
      </Box>
    </Box>
  );
}