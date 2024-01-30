import { createSignal, Show } from "solid-js";
import { useLocation, useSearchParams } from "@solidjs/router";
import { getAtk, useUserinfo } from "../../stores/userinfo.tsx";

export default function OauthConnectPage() {
  const [title, setTitle] = createSignal("Connect Third-party");
  const [subtitle, setSubtitle] = createSignal("Via your Goatpass account");

  const [error, setError] = createSignal<string | null>(null);
  const [status, setStatus] = createSignal("Handshaking...");
  const [loading, setLoading] = createSignal(true);

  const [client, setClient] = createSignal<any>(null);

  const [searchParams] = useSearchParams();

  const userinfo = useUserinfo();
  const location = useLocation();

  async function preConnect() {
    const res = await fetch(`/api/auth/oauth/connect${location.search}`, {
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });

    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();

      if (data["session"]) {
        setStatus("Redirecting...");
        redirect(data["session"]);
      } else {
        setTitle(`Connect ${data["client"].name}`);
        setSubtitle(`Via ${userinfo?.displayName}`);
        setClient(data["client"]);
        setLoading(false);
      }
    }
  }

  function decline() {
    if (window.history.length > 0) {
      window.history.back();
    } else {
      window.close();
    }
  }

  async function approve() {
    setLoading(true);
    setStatus("Approving...");

    const res = await fetch("/api/auth/oauth/connect?" + new URLSearchParams({
      client_id: searchParams["client_id"] as string,
      redirect_uri: encodeURIComponent(searchParams["redirect_uri"] as string),
      response_type: "code",
      scope: searchParams["scope"] as string
    }), {
      method: "POST",
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });

    if (res.status !== 200) {
      setError(await res.text());
      setLoading(false);
    } else {
      const data = await res.json();
      setStatus("Redirecting...");
      setTimeout(() => redirect(data["session"]), 1850);
    }
  }

  function redirect(session: any) {
    const url = `${searchParams["redirect_uri"]}?code=${session["grant_token"]}&state=${searchParams["state"]}`;
    window.open(url, "_self");
  }

  preConnect();

  return (
    <div class="w-full h-full flex justify-center items-center">
      <div class="card w-[480px] max-w-screen shadow-xl">
        <div class="card-body">
          <div id="header" class="text-center mb-5">
            <h1 class="text-xl font-bold">{title()}</h1>
            <p>{subtitle()}</p>
          </div>

          <Show when={error()}>
            <div id="alerts" class="mt-1">
              <div role="alert" class="alert alert-error">
                <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                     viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="capitalize">{error()}</span>
              </div>
            </div>
          </Show>

          <Show when={loading()}>
            <div class="py-16 text-center">
              <div class="text-center">
                <div>
                  <span class="loading loading-lg loading-bars"></span>
                </div>
                <span>{status()}</span>
              </div>
            </div>
          </Show>

          <Show when={!loading()}>
            <div class="mb-3">
              <h2 class="font-bold">About who you connecting to</h2>
              <p>{client().description}</p>
            </div>
            <div class="mb-3">
              <h2 class="font-bold">Make sure you trust them</h2>
              <p>You may share your personal information after connect them. Learn about their privacy policy and user
                agreement to keep your personal information in safe.</p>
            </div>
            <div class="mb-5">
              <h2 class="font-bold">After approve this request</h2>
              <p>
                You will be redirect to{" "}
                <span class="link link-primary cursor-not-allowed">{searchParams["redirect_uri"]}</span>
              </p>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2">
              <button class="btn btn-accent" onClick={() => decline()}>Decline</button>
              <button class="btn btn-primary" onClick={() => approve()}>Approve</button>
            </div>
          </Show>
        </div>
      </div>
    </div>
  );
}