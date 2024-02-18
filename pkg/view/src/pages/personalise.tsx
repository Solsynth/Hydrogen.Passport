import { getAtk, readProfiles, useUserinfo } from "../stores/userinfo.tsx";
import { createSignal, Show } from "solid-js";

export default function PersonalPage() {
  const userinfo = useUserinfo();

  const [error, setError] = createSignal<null | string>(null);
  const [success, setSuccess] = createSignal<null | string>(null);
  const [loading, setLoading] = createSignal(false);

  async function updateBasis(evt: SubmitEvent) {
    evt.preventDefault();

    const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));

    setLoading(true);
    const res = await fetch("/api/users/me", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getAtk()}`
      },
      body: JSON.stringify(data)
    });
    if (res.status !== 200) {
      setSuccess(null);
      setError(await res.text());
    } else {
      await readProfiles();
      setSuccess("Your basic information has been update.");
      setError(null);
    }
    setLoading(false);
  }

  async function updateAvatar(evt: SubmitEvent) {
    evt.preventDefault();

    setLoading(true);
    const data = new FormData(evt.target as HTMLFormElement);
    const res = await fetch("/api/avatar", {
      method: "PUT",
      headers: { "Authorization": `Bearer ${getAtk()}` },
      body: data
    });
    if (res.status !== 200) {
      setSuccess(null);
      setError(await res.text());
    } else {
      await readProfiles();
      setSuccess("Your avatar has been update.");
      setError(null);
    }
    setLoading(false);
  }

  return (
    <div class="max-w-[720px] mx-auto pt-12">
      <div class="px-5">
        <h1 class="text-2xl font-bold">Personalize</h1>
        <p>Customize your account and let us provide a better service to you.</p>
      </div>

      <div id="alerts">
        <Show when={error()}>
          <div role="alert" class="alert alert-error mt-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                 viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="capitalize">{error()}</span>
          </div>
        </Show>

        <Show when={success()}>
          <div role="alert" class="alert alert-success mt-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                 viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="capitalize">{success()}</span>
          </div>
        </Show>
      </div>

      <div class="card shadow-xl mt-5">

        <div class="card-body border-b border-base-200">
          <form class="grid grid-cols-1 gap-2" onSubmit={updateBasis}>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Username</span>
              </div>
              <input value={userinfo?.meta?.name} name="name" type="text" placeholder="Type here"
                     class="input input-bordered w-full" disabled />
            </label>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Nickname</span>
              </div>
              <input value={userinfo?.meta?.nick} name="nick" type="text" placeholder="Type here"
                     class="input input-bordered w-full" />
            </label>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-x-4">
              <label class="form-control w-full">
                <div class="label">
                  <span class="label-text">First Name</span>
                </div>
                <input value={userinfo?.meta?.profile?.first_name} name="first_name" type="text"
                       placeholder="Type here" class="input input-bordered w-full" />
              </label>
              <label class="form-control w-full">
                <div class="label">
                  <span class="label-text">Middle Name</span>
                </div>
                <input value={userinfo?.meta?.profile?.middle_name} name="middle_name" type="text"
                       placeholder="Type here" class="input input-bordered w-full" />
              </label>
              <label class="form-control w-full">
                <div class="label">
                  <span class="label-text">Last Name</span>
                </div>
                <input value={userinfo?.meta?.profile?.last_name} name="last_name" type="text"
                       placeholder="Type here" class="input input-bordered w-full" />
              </label>
            </div>

            <div>
              <button type="submit" class="btn btn-primary mt-5" disabled={loading()}>
                <Show when={loading()} fallback={"Save changes"}>
                  <span class="loading loading-spinner"></span>
                </Show>
              </button>
            </div>
          </form>
        </div>

        <div class="card-body">
          <form onSubmit={updateAvatar}>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Pick an avatar</span>
              </div>
              <input type="file" name="avatar" accept="image/*" class="file-input file-input-bordered w-full" />
              <div class="label">
                <span class="label-text-alt">Will took some time to apply to entire site</span>
              </div>
            </label>

            <button type="submit" class="btn btn-primary mt-5" disabled={loading()}>
              <Show when={loading()} fallback={"Save changes"}>
                <span class="loading loading-spinner"></span>
              </Show>
            </button>
          </form>
        </div>

      </div>
    </div>
  );
}