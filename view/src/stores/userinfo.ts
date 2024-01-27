import { createStore } from "solid-js/store";
import Cookie from "universal-cookie";

const [userinfo, setUserinfo] = createStore({
  isLoggedIn: false,
  displayName: "Citizen",
  profiles: null,
  meta: null
});

function checkLoggedIn(): boolean {
  return new Cookie().get("access_token");
}

export function getAtk(): string {
  return new Cookie().get("access_token");
}

export async function refreshAtk() {
  const rtk = new Cookie().get("refresh_token");

  const res = await fetch("/api/auth/token", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      code: rtk,
      grant_type: "refresh_token"
    })
  });
  if (res.status !== 200) {
    throw new Error(await res.text());
  } else {
    const data = await res.json();
    new Cookie().set("access_token", data["access_token"], { path: "/" });
    new Cookie().set("refresh_token", data["refresh_token"], { path: "/" });
  }
}

export async function readProfiles() {
  if (!checkLoggedIn()) return;

  const res = await fetch("/api/users/me", {
    headers: { "Authorization": `Bearer ${getAtk()}` }
  });

  if (res.status !== 200) {
    // Auto retry after refresh access token
    await refreshAtk();
    return await readProfiles();
  }

  const data = await res.json();

  setUserinfo({
    isLoggedIn: true,
    displayName: data["nick"],
    profiles: null,
    meta: data
  });
}

export { userinfo };