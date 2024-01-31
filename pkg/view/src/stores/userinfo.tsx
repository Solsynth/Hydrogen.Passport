import Cookie from "universal-cookie";
import { createContext, useContext } from "solid-js";
import { createStore } from "solid-js/store";

export interface Userinfo {
  isLoggedIn: boolean,
  displayName: string,
  profiles: any,
  meta: any
}

const UserinfoContext = createContext<Userinfo>();

const defaultUserinfo: Userinfo = {
  isLoggedIn: false,
  displayName: "Citizen",
  profiles: null,
  meta: null
};

const [userinfo, setUserinfo] = createStore<Userinfo>(structuredClone(defaultUserinfo));

export function getAtk(): string {
  return new Cookie().get("access_token");
}

export async function refreshAtk() {
  const rtk = new Cookie().get("refresh_token");

  const res = await fetch("/api/auth/token", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      refresh_token: rtk,
      grant_type: "refresh_token"
    })
  });
  if (res.status !== 200) {
    console.error(await res.text())
  } else {
    const data = await res.json();
    new Cookie().set("access_token", data["access_token"], { path: "/", maxAge: undefined });
    new Cookie().set("refresh_token", data["refresh_token"], { path: "/", maxAge: undefined });
  }
}

function checkLoggedIn(): boolean {
  return new Cookie().get("access_token");
}

export async function readProfiles(recovering = true) {
  if (!checkLoggedIn()) return;

  const res = await fetch("/api/users/me", {
    headers: { "Authorization": `Bearer ${getAtk()}` }
  });

  if (res.status !== 200) {
    if (recovering) {
      // Auto retry after refresh access token
      await refreshAtk();
      return await readProfiles(false);
    } else {
      clearUserinfo();
      window.location.reload();
    }
  }

  const data = await res.json();

  setUserinfo({
    isLoggedIn: true,
    displayName: data["nick"],
    profiles: null,
    meta: data
  });
}

export function clearUserinfo() {
  new Cookie().remove("access_token", { path: "/", maxAge: undefined });
  new Cookie().remove("refresh_token", { path: "/", maxAge: undefined });
  setUserinfo(defaultUserinfo);
}

export function UserinfoProvider(props: any) {
  return (
    <UserinfoContext.Provider value={userinfo}>
      {props.children}
    </UserinfoContext.Provider>
  );
}

export function useUserinfo() {
  return useContext(UserinfoContext);
}