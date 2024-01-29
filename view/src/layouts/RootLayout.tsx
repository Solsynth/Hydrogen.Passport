import Navbar from "./shared/Navbar.tsx";
import { readProfiles } from "../stores/userinfo.tsx";
import { createSignal, Show } from "solid-js";
import { readWellKnown } from "../stores/wellKnown.tsx";
import { BeforeLeaveEventArgs, useBeforeLeave, useNavigate } from "@solidjs/router";

export default function RootLayout(props: any) {
  const [ready, setReady] = createSignal(false);

  Promise.all([readWellKnown(), readProfiles()]).then(() => setReady(true));

  const navigate = useNavigate()

  useBeforeLeave((e: BeforeLeaveEventArgs) => {
    const whitelist = ["/auth/login", "/auth/register", "/users/me/confirm"]

    if (!whitelist.includes(e.to.toString()) && !e.defaultPrevented) {
      e.preventDefault();
      navigate(`/auth/login?redirect_uri=${e.to.toString()}`)
    }
  });

  return (
    <Show when={ready()} fallback={
      <div class="h-screen w-screen flex justify-center items-center">
        <div>
          <span class="loading loading-lg loading-infinity"></span>
        </div>
      </div>
    }>
      <Navbar />
      <main class="h-[calc(100vh-68px)] mt-[68px]">{props.children}</main>
    </Show>
  );
}