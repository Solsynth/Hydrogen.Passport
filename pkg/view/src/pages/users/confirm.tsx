import { createSignal, Show } from "solid-js";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { readProfiles } from "../../stores/userinfo.tsx";

export default function ConfirmRegistrationPage() {
  const [error, setError] = createSignal<string | null>(null);
  const [status, setStatus] = createSignal("Confirming your account...");

  const [searchParams] = useSearchParams();

  const navigate = useNavigate();

  async function doConfirm() {
    if (!searchParams["tk"]) {
      setError("Bad Request: Code was not exists");
    }

    const res = await fetch("/api/users/me/confirm", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        code: searchParams["tk"]
      })
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setStatus("Confirmed. Redirecting to dashboard...");
      await readProfiles();
      navigate("/");
    }
  }

  doConfirm();

  return (
    <div class="w-full h-full flex justify-center items-center">
      <div class="card w-[480px] max-w-screen shadow-xl">
        <div class="card-body">
          <div id="header" class="text-center mb-5">
            <h1 class="text-xl font-bold">Confirm your account</h1>
            <p>Hold on, we are working on it. Almost finished.</p>
          </div>

          <div class="pt-16 text-center">
            <div class="text-center">
              <div>
                <span class="loading loading-lg loading-bars"></span>
              </div>
              <span>{status()}</span>
            </div>
          </div>

          <Show when={error()} fallback={<div class="mt-16"></div>}>
            <div id="alerts" class="mt-16">
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
        </div>
      </div>
    </div>
  );
}