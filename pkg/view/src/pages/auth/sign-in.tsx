import { Link as RouterLink, useNavigate, useSearchParams } from "react-router-dom";
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
  Link,
  Paper,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Typography
} from "@mui/material";
import { FormEvent, useState } from "react";
import { request } from "@/scripts/request.ts";
import { useUserinfo } from "@/stores/userinfo.tsx";
import LoginIcon from "@mui/icons-material/Login";
import SecurityIcon from "@mui/icons-material/Security";
import KeyIcon from "@mui/icons-material/Key";
import PasswordIcon from "@mui/icons-material/Password";
import EmailIcon from "@mui/icons-material/Email";

export default function SignInPage() {
  const [panel, setPanel] = useState(0);

  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const [factor, setFactor] = useState<number>();
  const [factorType, setFactorType] = useState<any>();

  const [factors, setFactors] = useState<any>(null);
  const [challenge, setChallenge] = useState<any>(null);

  const { readProfiles } = useUserinfo();

  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const handlers: any[] = [
    async (evt: FormEvent<HTMLFormElement>) => {
      evt.preventDefault();

      const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
      if (!data.id) return;

      setLoading(true);
      const res = await request("/api/auth", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
      });
      if (res.status !== 200) {
        setError(await res.text());
      } else {
        const data = await res.json();
        setFactors(data["factors"]);
        setChallenge(data["challenge"]);
        setPanel(1);
        setError(null);
      }
      setLoading(false);
    },
    async (evt: FormEvent<HTMLFormElement>) => {
      evt.preventDefault();

      if (!factor) return;

      setLoading(true);
      const res = await request(`/api/auth/factors/${factor}`, {
        method: "POST"
      });
      if (res.status !== 200 && res.status !== 204) {
        setError(await res.text());
      } else {
        const item = factors.find((item: any) => item.id === factor).type;
        setError(null);
        setPanel(2);
        setFactorType(factorTypes[item]);
      }
      setLoading(false);
    },
    async (evt: SubmitEvent) => {
      evt.preventDefault();

      const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
      if (!data.credentials) return;

      setLoading(true);
      const res = await request(`/api/auth`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          challenge_id: challenge?.id,
          factor_id: factor,
          secret: data.credentials
        })
      });
      if (res.status !== 200) {
        setError(await res.text());
      } else {
        const data = await res.json();
        if (data["is_finished"]) {
          await grantToken(data["session"]["grant_token"]);
          await readProfiles();
          callback();
        } else {
          setError(null);
          setPanel(1);
          setFactor(undefined);
          setFactorType(undefined);
          setChallenge(data["challenge"]);
        }
      }
      setLoading(false);
    }
  ];

  function callback() {
    if (searchParams.has("closable")) {
      window.close();
    } else if (searchParams.has("redirect_uri")) {
      window.open(searchParams.get("redirect_uri") ?? "/", "_self");
    } else {
      navigate("/users");
    }
  }

  function getFactorAvailable(factor: any) {
    const blacklist: number[] = challenge?.blacklist_factors ?? [];
    return blacklist.includes(factor.id);
  }

  const factorTypes = [
    { icon: <PasswordIcon />, label: "Password Verification", autoComplete: "password" },
    { icon: <EmailIcon />, label: "Email One Time Password", autoComplete: "one-time-code" }
  ];

  const elements = [
    (
      <>
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <LoginIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Welcome back
        </Typography>
        <Box component="form" onSubmit={handlers[panel]} sx={{ mt: 3, width: "100%" }}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                autoComplete="username"
                name="id"
                required
                fullWidth
                label="Account ID"
                helperText={"Use your username, email or phone number."}
                autoFocus
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
            {loading ? "Processing..." : "Next"}
          </Button>
        </Box>
      </>
    ),
    (
      <>
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <SecurityIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Verify that's you
        </Typography>
        <Box component="form" onSubmit={handlers[panel]} sx={{ mt: 3, width: "100%" }}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <ToggleButtonGroup
                exclusive
                orientation="vertical"
                color="info"
                value={factor}
                sx={{ width: "100%" }}
                onChange={(_, val) => setFactor(val)}
              >
                {factors?.map((item: any, idx: number) => (
                  <ToggleButton key={idx} value={item.id} disabled={getFactorAvailable(item)}>
                    <Grid container>
                      <Grid item xs={2}>
                        {factorTypes[item.type]?.icon}
                      </Grid>
                      <Grid item xs="auto">
                        {factorTypes[item.type]?.label}
                      </Grid>
                    </Grid>
                  </ToggleButton>
                ))}
              </ToggleButtonGroup>
            </Grid>
          </Grid>
          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={loading}
            sx={{ mt: 3, mb: 2 }}
          >
            {loading ? "Processing..." : "Next"}
          </Button>
        </Box>
      </>
    ),
    (
      <>
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <KeyIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Enter the credentials
        </Typography>
        <Box component="form" onSubmit={handlers[panel]} sx={{ mt: 3, width: "100%" }}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                autoComplete={factorType?.autoComplete ?? "password"}
                name="credentials"
                type="password"
                required
                fullWidth
                label="Credentials"
                autoFocus
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
            {loading ? "Processing..." : "Next"}
          </Button>
        </Box>
      </>
    )
  ];

  async function grantToken(tk: string) {
    const res = await request("/api/auth/token", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        code: tk,
        grant_type: "grant_token"
      })
    });
    if (res.status !== 200) {
      const err = await res.text();
      setError(err);
      throw new Error(err);
    } else {
      setError(null);
    }
  }

  return (
    <Box sx={{ height: "100vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
      <Box style={{ width: "100vw", maxWidth: "450px" }}>
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

          <Collapse in={challenge != null} unmountOnExit>
            <Box>
              <Paper square sx={{ pt: 3, px: 5, textAlign: "center" }}>
                <Typography sx={{ mb: 2 }}>
                  Risk <b className="font-mono">{challenge?.risk_level}</b>&nbsp;
                  Progress <b className="font-mono">{challenge?.progress}/{challenge?.requirements}</b>
                </Typography>
                <LinearProgress
                  variant="determinate"
                  value={challenge?.progress / challenge?.requirements * 100}
                  sx={{ width: "calc(100%+5rem)", mt: 1, mx: -5 }}
                />
              </Paper>
            </Box>
          </Collapse>
        </Card>

        <Grid container justifyContent="center" sx={{ mt: 2 }}>
          <Grid item>
            <Link component={RouterLink} to="/auth/sign-up" variant="body2">
              Haven't an account? Sign up!
            </Link>
          </Grid>
        </Grid>
      </Box>
    </Box>
  );
}