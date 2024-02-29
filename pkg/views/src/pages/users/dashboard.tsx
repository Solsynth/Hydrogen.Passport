import { Alert, Box, Card, CardContent, Container, Typography } from "@mui/material";
import { useUserinfo } from "@/stores/userinfo.tsx";

export function Component() {
  const { userinfo } = useUserinfo();

  return (
    <Container sx={{ pt: 5 }} maxWidth="md">
      <Box sx={{ px: 3 }}>
        <Typography variant="h5">Welcome, {userinfo?.displayName}</Typography>
        <Typography variant="body2">What can I help you today?</Typography>
      </Box>

      {
        !userinfo?.profiles?.confirmed_at &&
        <Alert severity="warning" sx={{ mt: 3, mx: 1 }}>
          Your account haven't confirmed yet. Go to your linked email
          inbox and check out our registration confirm email.
        </Alert>
      }

      <Box sx={{ px: 1, mt: 3 }}>
        <Typography variant="h6" sx={{ px: 2 }}>Frequently Asked Questions</Typography>

        <Card variant="outlined" sx={{ mt: 1 }}>
          <CardContent style={{ padding: "40px" }}>
            <Typography>没有人有问题。没有人敢有问题。鲁迅曾经说过：</Typography>
            <Typography sx={{ pl: 4 }} fontWeight="bold">解决不了问题，就解决提问题的人。 —— 鲁迅</Typography>
            <Typography>所以，我们的客诉率是 0% 哦～</Typography>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
}