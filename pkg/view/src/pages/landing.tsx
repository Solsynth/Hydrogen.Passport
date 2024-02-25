import { Button, Container, Grid, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

export default function LandingPage() {
  return (
    <Container sx={{ height: "calc(100vh - 64px)", display: "flex", alignItems: "center", textAlign: "center" }}>
      <Grid padding={5} spacing={8} container>
        <Grid item xs={12} md={6}>
          <Typography variant="h3">All Goatworks<sup>®</sup> Services</Typography>
          <Typography variant="h3">In a single account</Typography>

          <Typography variant="body2" sx={{ mt: 8 }}>That's</Typography>
          <Typography variant="h1">Goatpass</Typography>
          <Button component={RouterLink} to="/auth/sign-up" variant="contained" sx={{ mt: 2 }}>Getting Start</Button>
        </Grid>
        <Grid item xs={12} md={6} sx={{ order: { xs: -100, md: 0 } }}>
          <img src="/favicon.svg" alt="Logo" width={256} height={256} className="block mx-auto" />
        </Grid>
      </Grid>
    </Container>
  );
}