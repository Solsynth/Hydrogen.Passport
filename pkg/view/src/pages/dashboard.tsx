import { getAtk, readProfiles, useUserinfo } from "../stores/userinfo.tsx";
import { createSignal, For, Show } from "solid-js";

export default function DashboardPage() {
  const userinfo = useUserinfo();

  const [error, setError] = createSignal<string | null>(null);

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

  async function readNotification(item: any) {
    const res = await fetch(`/api/notifications/${item.id}/read`, {
      method: "PUT",
      headers: { Authorization: `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readProfiles();
      setError(null);
    }
  }

  return (
    <div class="max-w-[720px] mx-auto px-5 pt-12">
      <div id="greeting" class="px-5">
        <h1 class="text-2xl font-bold">{userinfo?.displayName}</h1>
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

      <div class="card shadow-xl mt-5">
        <div class="card-body">
          <h2 class="card-title">Notifications</h2>
          <div class="bg-base-200 mt-3 mx-[-32px]">
            <Show when={userinfo?.meta?.notifications?.length <= 0}>
              <table class="table">
                <tbody>
                <tr>
                  <td class="px-[32px]">You're done! There are no notifications unread for you.</td>
                </tr>
                </tbody>
              </table>
            </Show>
            <Show when={userinfo?.meta?.notifications?.length > 0}>
              <table class="table">
                <tbody>
                <For each={userinfo?.meta?.notifications}>
                  {item =>
                    <tr>
                      <td class="px-[32px]">
                        <h2 class="font-bold">{item.subject}</h2>
                        <p>{item.content}</p>
                        <div class="flex gap-2">
                          <Show when={item.is_important}>
                            <span class="font-bold">Important</span>
                          </Show>
                          <a class="link" onClick={() => readNotification(item)}>Mark as read</a>
                        </div>
                      </td>
                    </tr>
                  }
                </For>
                </tbody>
              </table>
            </Show>
          </div>
        </div>
      </div>

    </div>
  );
}