import { readProfiles } from "../../stores/userinfo.tsx";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { createSignal, For, Match, Show, Switch } from "solid-js";
import Cookie from "universal-cookie";

export default function LoginPage() {
  const [title, setTitle] = createSignal("Sign in");
  const [subtitle, setSubtitle] = createSignal("Via your Goatpass account");

  const [error, setError] = createSignal<null | string>(null);
  const [loading, setLoading] = createSignal(false);

  const [factor, setFactor] = createSignal<number>();
  const [factors, setFactors] = createSignal<any[]>([]);
  const [challenge, setChallenge] = createSignal<any>();
  const [stage, setStage] = createSignal("starting");

  const [searchParams] = useSearchParams();

  const navigate = useNavigate();

  const handlers: { [id: string]: any } = {
    "starting": async (evt: SubmitEvent) => {
      evt.preventDefault();

      const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
      if (!data.id) return;

      setLoading(true);
      const res = await fetch("/api/auth", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
      });
      if (res.status !== 200) {
        setError(await res.text());
      } else {
        const data = await res.json();
        setTitle(`Welcome, ${data["display_name"]}`);
        setSubtitle("Before continue, we need verify that's you");
        setFactors(data["factors"]);
        setChallenge(data["challenge"]);
        setError(null);
        setStage("choosing");
      }
      setLoading(false);
    },
    "choosing": async (evt: SubmitEvent) => {
      evt.preventDefault();

      const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
      if (!data.factor) return;

      setLoading(true);
      const res = await fetch(`/api/auth/factors/${data.factor}`, {
        method: "POST"
      });
      if (res.status !== 200 && res.status !== 204) {
        setError(await res.text());
      } else {
        setTitle(`Enter the code`);
        setSubtitle(res.status === 204 ? "Enter your credentials" : "Code has been sent to your inbox");
        setError(null);
        setFactor(parseInt(data.factor as string));
        setStage("verifying");
      }
      setLoading(false);
    },
    "verifying": async (evt: SubmitEvent) => {
      evt.preventDefault();

      const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
      if (!data.credentials) return;

      setLoading(true);
      const res = await fetch(`/api/auth`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          challenge_id: challenge().id,
          factor_id: factor(),
          secret: data.credentials
        })
      });
      if (res.status !== 200) {
        setError(await res.text());
      } else {
        const data = await res.json();
        if (data["is_finished"]) {
          await grantToken(data["session"]["grant_token"]);
          await readProfiles();
          navigate(searchParams["redirect_uri"] ?? "/");
        } else {
          setError(null);
          setStage("choosing");
          setTitle("Continue verifying");
          setSubtitle("You passed one check, but that's not enough.");
          setChallenge(data["challenge"]);
        }
      }
      setLoading(false);
    }
  };

  async function grantToken(tk: string) {
    const res = await fetch("/api/auth/token", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        code: tk,
        grant_type: "grant_token"
      })
    });
    if (res.status !== 200) {
      const err = await res.text();
      setError(err);
      throw new Error(err);
    } else {
      const data = await res.json();
      new Cookie().set("access_token", data["access_token"], { path: "/" });
      new Cookie().set("refresh_token", data["refresh_token"], { path: "/" });
      setError(null);
    }
  }

  function getFactorAvailable(factor: any) {
    const blacklist: number[] = challenge()?.blacklist_factors ?? [];
    return blacklist.includes(factor.id);
  }

  function getFactorName(factor: any) {
    switch (factor.type) {
      case 0:
        return "Password Verification";
      case 1:
        return "Email Verification Code";
      default:
        return "Unknown";
    }
  }

  return (
    <div class="w-full h-full flex justify-center items-center">
      <div>
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

            <form id="form" onSubmit={(e) => handlers[stage()](e)}>
              <Switch>
                <Match when={stage() === "starting"}>
                  <label class="form-control w-full">
                    <div class="label">
                      <span class="label-text">Account ID</span>
                    </div>
                    <input name="id" type="text" placeholder="Type here" class="input input-bordered w-full" />
                    <div class="label">
                      <span class="label-text-alt">Your username, email or phone number.</span>
                    </div>
                  </label>
                </Match>
                <Match when={stage() === "choosing"}>
                  <div class="join join-vertical w-full">
                    <For each={factors()}>
                      {item =>
                        <input class="join-item btn" type="radio" name="factor"
                               value={item.id}
                               disabled={getFactorAvailable(item)}
                               aria-label={getFactorName(item)}
                        />
                      }
                    </For>
                  </div>
                  <p class="text-center text-sm mt-2">Choose a way to verify that's you</p>
                </Match>
                <Match when={stage() === "verifying"}>
                  <label class="form-control w-full">
                    <div class="label">
                      <span class="label-text">Credentials</span>
                    </div>
                    <input name="credentials" type="password" placeholder="Type here"
                           class="input input-bordered w-full" />
                    <div class="label">
                      <span class="label-text-alt">Password or one time password.</span>
                    </div>
                  </label>
                </Match>
              </Switch>

              <button type="submit" class="btn btn-primary btn-block mt-3" disabled={loading()}>
                <Show when={loading()} fallback={"Next"}>
                  <span class="loading loading-spinner"></span>
                </Show>
              </button>
            </form>
          </div>
        </div>

        <Show when={searchParams["redirect_uri"]}>
          <div id="redirect-info" class="mt-3">
            <div role="alert" class="alert shadow-xl">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                   class="stroke-info shrink-0 w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span>You need to login before access that.</span>
            </div>
          </div>
        </Show>

        <div class="text-sm text-center mt-3">
          <a target="_blank" href="/auth/register?closable=yes" class="link">Haven't an account? Click here to create
            one!</a>
        </div>
      </div>
    </div>
  );
}