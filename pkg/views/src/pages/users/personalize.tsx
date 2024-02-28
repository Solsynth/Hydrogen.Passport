import {
  Alert,
  Avatar,
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Collapse,
  Container,
  Divider,
  Grid,
  LinearProgress,
  Snackbar,
  styled,
  TextField,
  Typography
} from "@mui/material";
import { useUserinfo } from "@/stores/userinfo.tsx";
import { ChangeEvent, FormEvent, useState } from "react";
import { DatePicker } from "@mui/x-date-pickers";
import { request } from "@/scripts/request.ts";
import SaveIcon from "@mui/icons-material/Save";
import PublishIcon from "@mui/icons-material/Publish";
import NoAccountsIcon from "@mui/icons-material/NoAccounts";
import dayjs from "dayjs";

const VisuallyHiddenInput = styled("input")({
  clip: "rect(0 0 0 0)",
  clipPath: "inset(50%)",
  height: 1,
  overflow: "hidden",
  position: "absolute",
  bottom: 0,
  left: 0,
  whiteSpace: "nowrap",
  width: 1
});

export default function PersonalizePage() {
  const { userinfo, readProfiles, getAtk } = useUserinfo();

  const [done, setDone] = useState(false);
  const [error, setError] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  async function submit(evt: FormEvent<HTMLFormElement>) {
    evt.preventDefault();

    const data: any = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
    if (data.birthday) data.birthday = new Date(data.birthday);

    setLoading(true);
    const res = await request("/api/users/me", {
      method: "PUT",
      headers: { "Content-Type": "application/json", "Authorization": `Bearer ${getAtk()}` },
      body: JSON.stringify(data)
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readProfiles();
      setDone(true);
      setError(null);
    }
    setLoading(false);
  }

  async function changeAvatar(evt: ChangeEvent<HTMLInputElement>) {
    if (!evt.target.files) return;

    const file = evt.target.files[0];
    const payload = new FormData();
    payload.set("avatar", file);

    setLoading(true);
    const res = await request("/api/avatar", {
      method: "PUT",
      headers: { "Authorization": `Bearer ${getAtk()}` },
      body: payload
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readProfiles();
      setDone(true);
      setError(null);
    }
    setLoading(false);
  }

  function getBirthday() {
    return userinfo?.data?.profile?.birthday ? dayjs(userinfo?.data?.profile?.birthday) : undefined;
  }

  const basisForm = (
    <Box component="form" onSubmit={submit} sx={{ mt: 3 }}>
      <Grid container spacing={2}>
        <Grid item xs={6}>
          <TextField
            name="name"
            required
            disabled
            fullWidth
            label="Username"
            autoComplete="username"
            defaultValue={userinfo?.data?.name}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>
        <Grid item xs={6}>
          <TextField
            name="nick"
            required
            fullWidth
            label="Nickname"
            autoComplete="nickname"
            defaultValue={userinfo?.data?.nick}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>
        <Grid item xs={12}>
          <TextField
            name="description"
            multiline
            fullWidth
            label="Description"
            autoComplete="bio"
            defaultValue={userinfo?.data?.description}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>
        <Grid item xs={6} md={4}>
          <TextField
            name="first_name"
            fullWidth
            label="First Name"
            autoComplete="given_name"
            defaultValue={userinfo?.data?.profile?.first_name}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>
        <Grid item xs={6} md={4}>
          <TextField
            name="last_name"
            fullWidth
            label="Last Name"
            autoComplete="famliy_name"
            defaultValue={userinfo?.data?.profile?.last_name}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>
        <Grid item xs={12} md={4}>
          <DatePicker
            name="birthday"
            label="Birthday"
            defaultValue={getBirthday()}
            sx={{ width: "100%" }}
          />
        </Grid>
      </Grid>

      <Button
        type="submit"
        variant="contained"
        disabled={loading}
        startIcon={<SaveIcon />}
        sx={{ mt: 2, width: "180px" }}
      >
        Save changes
      </Button>

      <Divider sx={{ my: 2, mx: -3 }} />

      <Box sx={{ mt: 2.5, display: "flex", gap: 1, alignItems: "center" }}>
        <Box>
          <Avatar
            sx={{ width: 32, height: 32 }}
            alt={userinfo?.displayName}
            src={`/api/avatar/${userinfo?.data?.avatar}`}
          >
            <NoAccountsIcon />
          </Avatar>
        </Box>
        <Box>
          {/* @ts-ignore */}
          <Button
            type="button"
            color="info"
            component="label"
            tabIndex={-1}
            disabled={loading}
            startIcon={<PublishIcon />}
            sx={{ width: "180px" }}
          >
            Change avatar
            <VisuallyHiddenInput type="file" accept="image/*" onChange={changeAvatar} />
          </Button>
        </Box>
      </Box>
    </Box>
  );

  return (
    <Container sx={{ pt: 5 }} maxWidth="md">
      <Box sx={{ px: 3 }}>
        <Typography variant="h5">Personalize</Typography>
        <Typography variant="body2">
          Customize your appearance and name card across all Goatworks information.
        </Typography>
      </Box>

      <Collapse in={error}>
        <Alert severity="error" className="capitalize" sx={{ mt: 1.5 }}>{error}</Alert>
      </Collapse>

      <Box sx={{ mt: 2 }}>
        <Card variant="outlined">
          <Collapse in={loading}>
            <LinearProgress />
          </Collapse>

          <CardContent style={{ padding: "20px 24px" }}>
            <Box sx={{ px: 1, my: 1 }}>
              <Typography variant="h6">Information</Typography>
              <Typography variant="subtitle2">
                The information for public. Let us and others better to know who you are.
              </Typography>
            </Box>

            {
              userinfo?.data != null ? basisForm :
                <Box sx={{ pt: 1, px: 1 }}>
                  <CircularProgress />
                </Box>
            }

          </CardContent>
        </Card>
      </Box>

      <Snackbar
        open={done}
        autoHideDuration={1000 * 10}
        onClose={() => setDone(false)}
        message="Your profile has been updated. Some settings maybe need sometime to apply across site."
      />
    </Container>
  );
}