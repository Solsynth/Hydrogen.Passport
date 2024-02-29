import { Alert, Box, Collapse, IconButton, LinearProgress, List, ListItem, ListItemText } from "@mui/material";
import { useUserinfo } from "@/stores/userinfo.tsx";
import { request } from "@/scripts/request.ts";
import { useEffect, useState } from "react";
import { TransitionGroup } from "react-transition-group";
import MarkEmailReadIcon from "@mui/icons-material/MarkEmailRead";

export function Component() {
  const { userinfo, readProfiles, getAtk } = useUserinfo();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<null | string>(null);

  const [notifications, setNotifications] = useState<any[]>([]);

  async function readNotifications() {
    const res = await request(`/api/notifications?take=100`, {
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setNotifications(data["data"]);
      setError(null);
    }
  }

  async function markNotifications(item: any) {
    setLoading(true);
    const res = await request(`/api/notifications/${item.id}/read`, {
      method: "PUT",
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      readNotifications().then(() => readProfiles());
      setError(null);
    }
    setLoading(false);
  }

  useEffect(() => {
    readNotifications().then(() => setLoading(false));
  }, []);

  return (
    <Box>
      <Collapse in={loading}>
        <LinearProgress color="info" />
      </Collapse>

      <Collapse in={error != null}>
        <Alert severity="error" variant="filled" square>{error}</Alert>
      </Collapse>

      <Collapse in={userinfo?.data?.notifications?.length <= 0}>
        <Alert severity="success" variant="filled" square>You are done! There's no unread notifications for you.</Alert>
      </Collapse>

      <List sx={{ width: "100%", bgcolor: "background.paper" }}>
        <TransitionGroup>
          {notifications.map((item, idx) => (
            <Collapse key={idx} sx={{ px: 5 }}>
              <ListItem alignItems="flex-start" secondaryAction={
                <IconButton
                  edge="end"
                  aria-label="delete"
                  title="Delete"
                  onClick={() => markNotifications(item)}
                >
                  <MarkEmailReadIcon />
                </IconButton>
              }>
                <ListItemText
                  primary={item.subject}
                  secondary={item.content}
                />
              </ListItem>
            </Collapse>
          ))}
        </TransitionGroup>
      </List>
    </Box>
  );
}