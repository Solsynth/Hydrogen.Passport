import { Link as RouterLink, useRouteError } from "react-router-dom";
import { Box, Container, Link, Typography } from "@mui/material";

export default function ErrorBoundary() {
  const error = useRouteError() as any;

  return (
    <Container sx={{
      height: "100vh",
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      textAlign: "center"
    }}>
      <Box>
        <Typography variant="h1">{error.status}</Typography>
        <Typography variant="h6" sx={{ mb: 2 }}>{error?.message ?? "Something went wrong"}</Typography>

        <Link component={RouterLink} to="/">Back to homepage</Link>
      </Box>
    </Container>
  );
}