import { getAtk } from "../stores/userinfo.tsx";
import { createSignal, For, Match, Show, Switch } from "solid-js";

export default function DashboardPage() {
  const [challenges, setChallenges] = createSignal<any[]>([]);
  const [challengeCount, setChallengeCount] = createSignal(0);
  const [sessions, setSessions] = createSignal<any[]>([]);
  const [sessionCount, setSessionCount] = createSignal(0);
  const [events, setEvents] = createSignal<any[]>([]);
  const [eventCount, setEventCount] = createSignal(0);

  const [error, setError] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const [contentTab, setContentTab] = createSignal(0);

  async function readChallenges() {
    const res = await fetch("/api/users/me/challenges?take=10", {
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setChallenges(data["data"]);
      setChallengeCount(data["count"]);
    }
  }

  async function readSessions() {
    const res = await fetch("/api/users/me/sessions?take=10", {
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setSessions(data["data"]);
      setSessionCount(data["count"]);
    }
  }

  async function readEvents() {
    const res = await fetch("/api/users/me/events?take=10", {
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setEvents(data["data"]);
      setEventCount(data["count"]);
    }
  }

  async function killSession(item: any) {
    setSubmitting(true);
    const res = await fetch(`/api/users/me/sessions/${item.id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readSessions();
      setError(null);
    }
    setSubmitting(false);
  }

  readChallenges();
  readSessions();
  readEvents();

  return (
    <div class="max-w-[720px] mx-auto pt-12">
      <div id="greeting" class="px-5">
        <h1 class="text-2xl font-bold">Security</h1>
        <p>Here is your account status of security.</p>
      </div>

      <div id="alerts">
        <Show when={error()}>
          <div role="alert" class="alert alert-error mt-5">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                 viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="capitalize">{error()}</span>
          </div>
        </Show>
      </div>

      <div id="overview" class="mt-5">
        <div class="stats shadow w-full">
          <div class="stat">
            <div class="stat-figure text-secondary">
              <i class="fa-solid fa-door-open inline-block text-[28px] w-8 h-8 stroke-current"></i>
            </div>
            <div class="stat-title">Challenges</div>
            <div class="stat-value">{challengeCount()}</div>
          </div>

          <div class="stat">
            <div class="stat-figure text-secondary">
              <i class="fa-solid fa-key inline-block text-[28px] w-8 h-8 stroke-current"></i>
            </div>
            <div class="stat-title">Sessions</div>
            <div class="stat-value">{sessionCount()}</div>
          </div>

          <div class="stat">
            <div class="stat-figure text-secondary">
              <i class="fa-solid fa-person-walking inline-block text-[28px] w-8 h-8 stroke-current"></i>
            </div>
            <div class="stat-title">Events</div>
            <div class="stat-value">{eventCount()}</div>
          </div>
        </div>
      </div>

      <div id="switch-area" class="mt-5">
        <div role="tablist" class="tabs tabs-boxed">
          <input type="radio" name="content-switch" role="tab" class="tab" aria-label="Challenges"
                 checked={contentTab() === 0} onChange={() => setContentTab(0)} />
          <input type="radio" name="content-switch" role="tab" class="tab" aria-label="Sessions"
                 checked={contentTab() === 1} onChange={() => setContentTab(1)} />
          <input type="radio" name="content-switch" role="tab" class="tab" aria-label="Events"
                 checked={contentTab() === 2} onChange={() => setContentTab(2)} />
        </div>
      </div>

      <div id="data-area" class="mt-5 shadow">
        <Switch>
          <Match when={contentTab() === 0}>
            <div class="overflow-x-auto">
              <table class="table">
                <thead>
                <tr>
                  <th></th>
                  <th>State</th>
                  <th>IP Address</th>
                  <th>User Agent</th>
                  <th>Date</th>
                </tr>
                </thead>
                <tbody>
                <For each={challenges()}>
                  {item => <tr>
                    <th>{item.id}</th>
                    <td>{item.state}</td>
                    <td>{item.ip_address}</td>
                    <td>
                        <span class="tooltip" data-tip={item.user_agent}>
                          {item.user_agent.substring(0, 10) + "..."}
                        </span>
                    </td>
                    <td>{new Date(item.created_at).toLocaleString()}</td>
                  </tr>}
                </For>
                </tbody>
              </table>
            </div>
          </Match>

          <Match when={contentTab() === 1}>
            <div class="overflow-x-auto">
              <table class="table">
                <thead>
                <tr>
                  <th></th>
                  <th>Third Client</th>
                  <th>Audiences</th>
                  <th>Date</th>
                  <th></th>
                </tr>
                </thead>
                <tbody>
                <For each={sessions()}>
                  {item => <tr>
                    <th>{item.id}</th>
                    <td>{item.client_id ? "Linked" : "Non-linked"}</td>
                    <td>{item.audiences?.join(", ")}</td>
                    <td>{new Date(item.created_at).toLocaleString()}</td>
                    <td class="py-0">
                      <button disabled={submitting()} onClick={() => killSession(item)}>
                        <i class="fa-solid fa-right-from-bracket"></i>
                      </button>
                    </td>
                  </tr>}
                </For>
                </tbody>
              </table>
            </div>
          </Match>

          <Match when={contentTab() === 2}>
            <div class="overflow-x-auto">
              <table class="table">
                <thead>
                <tr>
                  <th></th>
                  <th>Type</th>
                  <th>Target</th>
                  <th>IP Address</th>
                  <th>User Agent</th>
                  <th>Date</th>
                </tr>
                </thead>
                <tbody>
                <For each={events()}>
                  {item => <tr>
                    <th>{item.id}</th>
                    <td>{item.type}</td>
                    <td>{item.target}</td>
                    <td>{item.ip_address}</td>
                    <td>
                        <span class="tooltip" data-tip={item.user_agent}>
                          {item.user_agent.substring(0, 10) + "..."}
                        </span>
                    </td>
                    <td>{new Date(item.created_at).toLocaleString()}</td>
                  </tr>}
                </For>
                </tbody>
              </table>
            </div>
          </Match>
        </Switch>
      </div>
    </div>
  );
}