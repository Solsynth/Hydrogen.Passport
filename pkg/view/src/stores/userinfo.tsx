import Cookie from "universal-cookie";
import { createContext, useContext } from "solid-js";
import { createStore } from "solid-js/store";
import { request } from "../scripts/request.ts";

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
  return new Cookie().get("identity_auth_key");
}

function checkLoggedIn(): boolean {
  return new Cookie().get("identity_auth_key");
}

export async function readProfiles() {
  if (!checkLoggedIn()) return;

  const res = await request("/api/users/me", {
    credentials: "include"
  });

  if (res.status !== 200) {
    clearUserinfo();
    window.location.reload();
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
  const cookies = document.cookie.split(";");
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i];
    const eqPos = cookie.indexOf("=");
    const name = eqPos > -1 ? cookie.substring(0, eqPos) : cookie;
    document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
  }

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