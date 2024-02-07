import Navbar from "./shared/Navbar.tsx";
import { readProfiles, useUserinfo } from "../stores/userinfo.tsx";
import { createEffect, createMemo, createSignal, Show } from "solid-js";
import { readWellKnown } from "../stores/wellKnown.tsx";
import { BeforeLeaveEventArgs, useLocation, useNavigate, useSearchParams } from "@solidjs/router";

export default function RootLayout(props: any) {
  const [ready, setReady] = createSignal(false);

  Promise.all([readWellKnown(), readProfiles()]).then(() => setReady(true));

  const navigate = useNavigate();
  const userinfo = useUserinfo();

  const location = useLocation();
  const [searchParams] = useSearchParams();

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

  const mainContentStyles = createMemo(() => {
    if (searchParams["embedded"]) {
      return "h-screen";
    } else {
      return "h-[calc(100vh-64px)] mt-[64px]";
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
      <Show when={!searchParams["embedded"]}>
        <Navbar />
      </Show>

      <main class={`${mainContentStyles()} px-5`}>{props.children}</main>
    </Show>
  );
}