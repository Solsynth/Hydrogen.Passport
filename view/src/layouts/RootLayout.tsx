import Navbar from "./shared/Navbar.tsx";
import { readProfiles, UserinfoProvider } from "../stores/userinfo.tsx";
import { createSignal, Show } from "solid-js";

export default function RootLayout(props: any) {
  const [ready, setReady] = createSignal(false);

  readProfiles().then(() => setReady(true));

  return (
    <Show when={ready()} fallback={
      <div class="h-screen w-screen flex justify-center items-center">
        <div>
          <span class="loading loading-lg loading-infinity"></span>
        </div>
      </div>
    }>
      <UserinfoProvider>
        <Navbar />

        <main class="h-[calc(100vh-68px)]">{props.children}</main>
      </UserinfoProvider>
    </Show>
  );
}