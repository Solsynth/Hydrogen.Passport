import { useUserinfo } from "../stores/userinfo.tsx";
import { Show } from "solid-js";

export default function DashboardPage() {
  const userinfo = useUserinfo();

  return (
    <div class="max-w-[720px] mx-auto px-5 pt-12">
      <h1 class="text-2xl font-bold">Welcome, {userinfo?.displayName}</h1>
      <p>What's a nice day!</p>

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
      </div>
    </div>
  );
}