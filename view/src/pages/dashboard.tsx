import { getAtk, useUserinfo } from "../stores/userinfo.tsx";
import { createSignal, For, Show } from "solid-js";

export default function DashboardPage() {
  const userinfo = useUserinfo();

  function getGreeting() {
    const currentHour = new Date().getHours();

    if (currentHour >= 0 && currentHour < 12) {
      return "Good morning! Wishing you a day filled with joy and success. ☀️";
    } else if (currentHour >= 12 && currentHour < 18) {
      return "Afternoon! Hope you have a productive and joyful afternoon! ☀️";
    } else {
      return "Good evening! Wishing you a relaxing and pleasant evening. 🌙";
    }
  }

  const [events, setEvents] = createSignal<any[]>([]);
  const [eventCount, setEventCount] = createSignal(0);

  const [error, setError] = createSignal<string | null>(null);

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

  readEvents();

  return (
    <div class="max-w-[720px] mx-auto px-5 pt-12">
      <div id="greeting" class="px-5">
        <h1 class="text-2xl font-bold">Welcome, {userinfo?.displayName}</h1>
        <p>{getGreeting()}</p>
      </div>

      <div id="alerts">
        <Show when={!userinfo?.meta?.confirmed_at}>
          <div role="alert" class="alert alert-warning mt-5">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                 viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div>
              <span>Your account isn't confirmed yet. Please check your inbox and confirm your account.</span> <br />
              <span>Otherwise your account will be deactivate after 48 hours.</span>
            </div>
          </div>
        </Show>
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
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                   class="inline-block w-8 h-8 stroke-current">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <div class="stat-title">Challenges</div>
            <div class="stat-value">{userinfo?.meta?.challenges?.length}</div>
          </div>

          <div class="stat">
            <div class="stat-figure text-secondary">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                   class="inline-block w-8 h-8 stroke-current">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"></path>
              </svg>
            </div>
            <div class="stat-title">Sessions</div>
            <div class="stat-value">{userinfo?.meta?.sessions?.length}</div>
          </div>

          <div class="stat">
            <div class="stat-figure text-secondary">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                   class="inline-block w-8 h-8 stroke-current">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path>
              </svg>
            </div>
            <div class="stat-title">Events</div>
            <div class="stat-value">{eventCount()}</div>
          </div>
        </div>
      </div>

      <div id="data-area" class="mt-5 shadow">
        <div class="join join-vertical w-full">

          <details class="collapse collapse-plus join-item border border-base-300">
            <summary class="collapse-title text-lg font-medium">
              Challenges
            </summary>
            <div class="collapse-content mx-[-16px]">
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
                  <For each={userinfo?.meta?.challenges ?? []}>
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
            </div>
          </details>

          <details class="collapse collapse-plus join-item border border-base-300">
            <summary class="collapse-title text-lg font-medium">
              Sessions
            </summary>
            <div class="collapse-content mx-[-16px]">
              <table class="table">
                <thead>
                <tr>
                  <th></th>
                  <th>Third Client</th>
                  <th>Audiences</th>
                  <th>Date</th>
                </tr>
                </thead>
                <tbody>
                <For each={userinfo?.meta?.sessions ?? []}>
                  {item => <tr>
                    <th>{item.id}</th>
                    <td>{item.client_id ? "Linked" : "Non-linked"}</td>
                    <td>{item.audiences?.join(", ")}</td>
                    <td>{new Date(item.created_at).toLocaleString()}</td>
                  </tr>}
                </For>
                </tbody>
              </table>
            </div>
          </details>

          <details class="collapse collapse-plus join-item border border-base-300">
            <summary class="collapse-title text-lg font-medium">
              Events
            </summary>
            <div class="collapse-content mx-[-16px]">
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
            </div>
          </details>

        </div>
      </div>
    </div>
  );
}