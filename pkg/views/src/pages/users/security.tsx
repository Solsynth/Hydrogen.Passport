import {
  Alert,
  Box,
  Card,
  CardContent,
  Collapse,
  Container,
  Grid,
  LinearProgress,
  Tab,
  Tabs,
  Typography
} from "@mui/material";
import { useUserinfo } from "@/stores/userinfo.tsx";
import { TabContext, TabPanel } from "@mui/lab";
import { useEffect, useState } from "react";
import { DataGrid, GridActionsCellItem, GridColDef, GridRowParams, GridValueGetterParams } from "@mui/x-data-grid";
import { request } from "@/scripts/request.ts";
import ExitToAppIcon from "@mui/icons-material/ExitToApp";


export function Component() {
  const dataDefinitions: { [id: string]: GridColDef[] } = {
    challenges: [
      { field: "id", headerName: "ID", width: 64 },
      { field: "ip_address", headerName: "IP Address", minWidth: 128 },
      { field: "user_agent", headerName: "User Agent", minWidth: 320 },
      {
        field: "created_at",
        headerName: "Issued At",
        minWidth: 160,
        valueGetter: (params: GridValueGetterParams) => new Date(params.row.created_at).toLocaleString()
      }
    ],
    sessions: [
      { field: "id", headerName: "ID", width: 64 },
      {
        field: "audiences",
        headerName: "Audiences",
        minWidth: 128,
        valueGetter: (params: GridValueGetterParams) => params.row.audiences.join(", ")
      },
      {
        field: "claims",
        headerName: "Claims",
        minWidth: 224,
        valueGetter: (params: GridValueGetterParams) => params.row.claims.join(", ")
      },
      {
        field: "created_at",
        headerName: "Issued At",
        minWidth: 160,
        valueGetter: (params: GridValueGetterParams) => new Date(params.row.created_at).toLocaleString()
      },
      {
        field: "actions",
        type: "actions",
        getActions: (params: GridRowParams) => [
          <GridActionsCellItem
            icon={<ExitToAppIcon />}
            onClick={() => killSession(params.row)}
            disabled={loading}
            label="Sign Out"
          />
        ]
      }
    ],
    events: [
      { field: "id", headerName: "ID", width: 64 },
      { field: "type", headerName: "Type", minWidth: 128 },
      { field: "target", headerName: "Affected Object", minWidth: 128 },
      { field: "ip_address", headerName: "IP Address", minWidth: 128 },
      { field: "user_agent", headerName: "User Agent", minWidth: 128 },
      {
        field: "created_at",
        headerName: "Performed At",
        minWidth: 160,
        valueGetter: (params: GridValueGetterParams) => new Date(params.row.created_at).toLocaleString()
      }
    ]
  };

  const { getAtk } = useUserinfo();

  const [challenges, setChallenges] = useState<any[]>([]);
  const [challengeCount, setChallengeCount] = useState(0);
  const [sessions, setSessions] = useState<any[]>([]);
  const [sessionCount, setSessionCount] = useState(0);
  const [events, setEvents] = useState<any[]>([]);
  const [eventCount, setEventCount] = useState(0);

  const [pagination, setPagination] = useState({
    challenges: { page: 0, pageSize: 5 },
    sessions: { page: 0, pageSize: 5 },
    events: { page: 0, pageSize: 5 }
  });

  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [reverting] = useState({
    challenges: true,
    sessions: true,
    events: true
  });

  const [dataPane, setDataPane] = useState("challenges");

  async function readChallenges() {
    reverting.challenges = true;
    const res = await request("/api/users/me/challenges?" + new URLSearchParams({
      take: pagination.challenges.pageSize.toString(),
      offset: (pagination.challenges.page * pagination.challenges.pageSize).toString()
    }), {
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setChallenges(data["data"]);
      setChallengeCount(data["count"]);
    }
    reverting.challenges = false;
  }

  async function readSessions() {
    reverting.sessions = true;
    const res = await request("/api/users/me/sessions?" + new URLSearchParams({
      take: pagination.sessions.pageSize.toString(),
      offset: (pagination.sessions.page * pagination.sessions.pageSize).toString()
    }), {
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setSessions(data["data"]);
      setSessionCount(data["count"]);
    }
    reverting.sessions = false;
  }

  async function readEvents() {
    reverting.events = true;
    const res = await request("/api/users/me/events?" + new URLSearchParams({
      take: pagination.events.pageSize.toString(),
      offset: (pagination.events.page * pagination.events.pageSize).toString()
    }), {
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setEvents(data["data"]);
      setEventCount(data["count"]);
    }
    reverting.events = false;
  }

  async function killSession(item: any) {
    setLoading(true);
    const res = await request(`/api/users/me/sessions/${item.id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readSessions();
      setError(null);
    }
    setLoading(false);
  }

  useEffect(() => {
    readChallenges().then(() => console.log("Refreshed challenges list."));
  }, [pagination.challenges]);

  useEffect(() => {
    readSessions().then(() => console.log("Refreshed sessions list."));
  }, [pagination.sessions]);

  useEffect(() => {
    readEvents().then(() => console.log("Refreshed events list."));
  }, [pagination.events]);

  return (
    <Container sx={{ pt: 5 }} maxWidth="md">
      <Box sx={{ px: 3 }}>
        <Typography variant="h5">Security</Typography>
        <Typography variant="body2">
          Overview and control all security details in your account.
        </Typography>
      </Box>

      <Collapse in={error != null}>
        <Alert severity="error" className="capitalize" sx={{ mt: 1.5 }}>{error}</Alert>
      </Collapse>

      <Grid container>
        <Grid item xs={12}>
          <Box sx={{ mt: 2 }}>
            <Card variant="outlined">
              <Collapse in={loading}>
                <LinearProgress />
              </Collapse>

              <TabContext value={dataPane}>
                <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
                  <Tabs centered value={dataPane} onChange={(_, val) => setDataPane(val)}>
                    <Tab label="Challenges" value="challenges" />
                    <Tab label="Sessions" value="sessions" />
                    <Tab label="Events" value="events" />
                  </Tabs>
                </Box>

                <CardContent style={{ padding: "20px 24px" }}>
                  <TabPanel value={"challenges"}>
                    <DataGrid
                      pageSizeOptions={[5, 10, 15, 20, 25]}
                      paginationMode="server"
                      loading={reverting.challenges}
                      rows={challenges}
                      rowCount={challengeCount}
                      columns={dataDefinitions.challenges}
                      paginationModel={pagination.challenges}
                      onPaginationModelChange={(val) => setPagination({ ...pagination, challenges: val })}
                      checkboxSelection
                    />
                  </TabPanel>
                  <TabPanel value={"sessions"}>
                    <DataGrid
                      pageSizeOptions={[5, 10, 15, 20, 25]}
                      paginationMode="server"
                      loading={reverting.sessions}
                      rows={sessions}
                      rowCount={sessionCount}
                      columns={dataDefinitions.sessions}
                      paginationModel={pagination.sessions}
                      onPaginationModelChange={(val) => setPagination({ ...pagination, sessions: val })}
                      checkboxSelection
                    />
                  </TabPanel>
                  <TabPanel value={"events"}>
                    <DataGrid
                      pageSizeOptions={[5, 10, 15, 20, 25]}
                      paginationMode="server"
                      loading={reverting.events}
                      rows={events}
                      rowCount={eventCount}
                      columns={dataDefinitions.events}
                      paginationModel={pagination.events}
                      onPaginationModelChange={(val) => setPagination({ ...pagination, events: val })}
                      checkboxSelection
                    />
                  </TabPanel>
                </CardContent>
              </TabContext>
            </Card>
          </Box>
        </Grid>
      </Grid>
    </Container>
  );
}