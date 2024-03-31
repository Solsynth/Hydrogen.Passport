import { defineStore } from "pinia";
import { ref } from "vue";
import { checkLoggedIn, getAtk } from "@/stores/userinfo";
import { request } from "@/scripts/request";

export const useNotifications = defineStore("notifications", () => {
  let socket: WebSocket;

  const loading = ref(false);

  const notifications = ref<any[]>([]);
  const total = ref(0)

  async function list() {
    loading.value = true;
    const res = await request(
      "/api/notifications?" +
      new URLSearchParams({
        take: (25).toString(),
        offset: (0).toString()
      }),
      {
        headers: { Authorization: `Bearer ${getAtk()}` }
      }
    );
    if (res.status === 200) {
      const data = await res.json();
      notifications.value = data["data"];
      total.value = data["count"];
    }
    loading.value = false;
  }

  function remove(idx: number) {
    notifications.value.splice(idx, 1)
    total.value--;
  }

  async function connect() {
    if (!(checkLoggedIn())) return;

    const uri = `ws://${window.location.host}/api/notifications/listen`;

    socket = new WebSocket(uri + `?tk=${getAtk() as string}`);

    socket.addEventListener("open", (event) => {
      console.log("[NOTIFICATIONS] The listen websocket has been established... ", event.type);
    });
    socket.addEventListener("close", (event) => {
      console.warn("[NOTIFICATIONS] The listen websocket is disconnected... ", event.reason, event.code);
    });
    socket.addEventListener("message", (event) => {
      const data = JSON.parse(event.data);
      notifications.value.push(data);
      total.value++;
    });
  }

  function disconnect() {
    socket.close();
  }

  return { loading, notifications, total, list, remove, connect, disconnect };
});