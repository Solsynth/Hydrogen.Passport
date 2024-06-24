import Cookie from "universal-cookie"
import { defineStore } from "pinia"
import { ref } from "vue"
import { request } from "@/scripts/request"

export interface Userinfo {
  isLoggedIn: boolean
  displayName: string
  data: any
}

const defaultUserinfo: Userinfo = {
  isLoggedIn: false,
  displayName: "Citizen",
  data: null,
}

export function getAtk(): string {
  return new Cookie().get("__hydrogen_atk")
}

export function checkLoggedIn(): boolean {
  return new Cookie().get("__hydrogen_rtk")
}

export const useUserinfo = defineStore("userinfo", () => {
  const userinfo = ref(defaultUserinfo)
  const isReady = ref(false)

  async function readProfiles() {
    if (!checkLoggedIn()) {
      isReady.value = true
    }

    const res = await request("/api/users/me", {
      headers: { Authorization: `Bearer ${getAtk()}` },
    })

    if (res.status !== 200) {
      return
    }

    const data = await res.json()

    isReady.value = true
    userinfo.value = {
      isLoggedIn: true,
      displayName: data["nick"],
      data: data,
    }
  }

  return { userinfo, isReady, readProfiles }
})
