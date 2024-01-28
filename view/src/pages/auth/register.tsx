import { createSignal, Show } from "solid-js";

export default function RegisterPage() {
  const [title, setTitle] = createSignal("Create an account");
  const [subtitle, setSubtitle] = createSignal("The first step to join our community.");

  const [error, setError] = createSignal<null | string>(null);
  const [loading, setLoading] = createSignal(false);
  const [done, setDone] = createSignal(false);

  async function submit(evt: SubmitEvent) {
    evt.preventDefault();

    const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
    if (!data.name || !data.nick || !data.email || !data.password) return;

    setLoading(true);
    const res = await fetch("/api/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data)
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setTitle("Congratulations!");
      setSubtitle("Your account has been created and activation email has sent to your inbox!");
      setDone(true);
    }
    setLoading(false);
  }

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

          <Show when={!done()}>
            <form id="form" onSubmit={submit}>
              <label class="form-control w-full">
                <div class="label">
                  <span class="label-text">Username</span>
                  <span class="label-text-alt font-bold">Cannot be modify</span>
                </div>
                <input name="name" type="text" placeholder="Type here" class="input input-bordered w-full" />
                <div class="label">
                  <span class="label-text-alt">Lowercase alphabet and numbers only, maximum 16 characters</span>
                </div>
              </label>
              <label class="form-control w-full">
                <div class="label">
                  <span class="label-text">Nickname</span>
                </div>
                <input name="nick" type="text" placeholder="Type here" class="input input-bordered w-full" />
                <div class="label">
                  <span class="label-text-alt">Maximum length is 24 characters</span>
                </div>
              </label>
              <label class="form-control w-full">
                <div class="label">
                  <span class="label-text">Email Address</span>
                </div>
                <input name="email" type="email" placeholder="Type here" class="input input-bordered w-full" />
                <div class="label">
                  <span class="label-text-alt">Do not accept address with plus sign</span>
                </div>
              </label>
              <label class="form-control w-full">
                <div class="label">
                  <span class="label-text">Password</span>
                </div>
                <input name="password" type="password" placeholder="Type here" class="input input-bordered w-full" />
                <div class="label">
                  <span class="label-text-alt">Must be secure</span>
                </div>
              </label>

              <button type="submit" class="btn btn-primary btn-block mt-3" disabled={loading()}>
                <Show when={loading()} fallback={"Next"}>
                  <span class="loading loading-spinner"></span>
                </Show>
              </button>
            </form>
          </Show>

          <Show when={done()}>
            <div class="py-12 text-center">
              <h2 class="text-lg font-bold">What's next?</h2>
              <span>
                <a href="/auth/login" class="link">Go login</a>{" "}
                then you can take part in the entire smartsheep community.
              </span>
            </div>
          </Show>
        </div>

        <div class="text-sm text-center mt-3">
          <a href="/auth/login" class="link">Already had an account? Login now!</a>
        </div>
      </div>
    </div>
  );
}