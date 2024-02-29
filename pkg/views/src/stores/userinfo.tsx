import Cookie from "universal-cookie";
import { request } from "../scripts/request.ts";
import { createContext, useContext, useState } from "react";

export interface Userinfo {
  isReady: boolean,
  isLoggedIn: boolean,
  displayName: string,
  data: any,
}

const defaultUserinfo: Userinfo = {
  isReady: false,
  isLoggedIn: false,
  displayName: "Citizen",
  data: null
};

const UserinfoContext = createContext<any>({ userinfo: defaultUserinfo });

export function UserinfoProvider(props: any) {
  const [userinfo, setUserinfo] = useState<Userinfo>(structuredClone(defaultUserinfo));

  function getAtk(): string {
    return new Cookie().get("identity_auth_key");
  }

  function checkLoggedIn(): boolean {
    return new Cookie().get("identity_auth_key");
  }

  async function readProfiles() {
    if (!checkLoggedIn()) {
      setUserinfo((data) => {
        data.isReady = true;
        return data;
      });
    }

    const res = await request("/api/users/me", {
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });

    if (res.status !== 200) {
      clearUserinfo();
      window.location.reload();
    }

    const data = await res.json();

    setUserinfo({
      isReady: true,
      isLoggedIn: true,
      displayName: data["nick"],
      data: data
    });
  }

  function clearUserinfo() {
    const cookies = document.cookie.split(";");
    for (let i = 0; i < cookies.length; i++) {
      const cookie = cookies[i];
      const eqPos = cookie.indexOf("=");
      const name = eqPos > -1 ? cookie.substring(0, eqPos) : cookie;
      document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
    }

    setUserinfo(defaultUserinfo);
  }

  return (
    <UserinfoContext.Provider value={{ userinfo, readProfiles, checkLoggedIn, getAtk, clearUserinfo }}>
      {props.children}
    </UserinfoContext.Provider>
  );
}

export function useUserinfo() {
  return useContext(UserinfoContext);
}