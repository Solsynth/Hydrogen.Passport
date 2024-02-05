import Navbar from "./shared/Navbar.tsx";
import { readProfiles, useUserinfo } from "../stores/userinfo.tsx";
import { createEffect, createSignal, Show } from "solid-js";
import { readWellKnown } from "../stores/wellKnown.tsx";
import { BeforeLeaveEventArgs, useLocation, useNavigate } from "@solidjs/router";

export default function RootLayout(props: any) {
  const [ready, setReady] = createSignal(false);

  Promise.all([readWellKnown(), readProfiles()]).then(() => setReady(true));

  const navigate = useNavigate();
  const userinfo = useUserinfo();

  const location = useLocation();

  createEffect(() => {
    if (ready()) {
      keepGate(location.pathname + location.search);
    }
  }, [ready, userinfo, location]);

  function keepGate(path: string, e?: BeforeLeaveEventArgs) {
    const pathname = path.split("?")[0];
    const whitelist = ["/auth/login", "/auth/register", "/users/me/confirm"];

    if (!userinfo?.isLoggedIn && !whitelist.includes(pathname)) {
      if (!e?.defaultPrevented) e?.preventDefault();
      navigate(`/auth/login?redirect_uri=${encodeURIComponent(path)}`);
    }
  }

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