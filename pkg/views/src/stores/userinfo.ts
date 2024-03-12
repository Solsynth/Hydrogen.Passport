import Cookie from "universal-cookie"
import { defineStore } from "pinia"
import { ref } from "vue"
import { request } from "@/scripts/request"

export interface Userinfo {
  isReady: boolean
  isLoggedIn: boolean
  displayName: string
  data: any
}

const defaultUserinfo: Userinfo = {
  isReady: false,
  isLoggedIn: false,
  displayName: "Citizen",
  data: null
}

export function getAtk(): string {
  return new Cookie().get("identity_auth_key")
}

export function checkLoggedIn(): boolean {
  return new Cookie().get("identity_auth_key")
}

export const useUserinfo = defineStore("userinfo", () => {
  const userinfo = ref(defaultUserinfo)
  const isReady = ref(false)

  async function readProfiles() {
    if (!checkLoggedIn()) {
        isReady.value = true;
      }
  
      const res = await request("/api/users/me", {
        headers: { "Authorization": `Bearer ${getAtk()}` }
      });
  
      if (res.status !== 200) {
        return;
      }
  
      const data = await res.json();
  
      userinfo.value = {
        isReady: true,
        isLoggedIn: true,
        displayName: data["nick"],
        data: data
      };
  }

  return { userinfo, isReady, readProfiles }
})
