import { createRouter, createWebHistory } from "vue-router"
import MasterLayout from "@/layouts/master.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: MasterLayout,
      children: [
        { path: "/", name: "dashboard", component: () => import("@/views/dashboard.vue") },
      ],
    },
    {
      path: "/auth",
      children: [
        { path: "sign-in", name: "auth.sign-in", component: () => import("@/views/auth/sign-in.vue") },
        { path: "sign-up", name: "auth.sign-up", component: () => import("@/views/auth/sign-up.vue") },
      ]
    }
  ],
})

export default router
