import { createRouter, createWebHistory } from "vue-router"
import { useUserinfo } from "@/stores/userinfo"
import MasterLayout from "@/layouts/master.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: MasterLayout,
      children: [{ path: "/", name: "dashboard", component: () => import("@/views/dashboard.vue") }],
    },
    {
      path: "/auth",
      children: [
        {
          path: "sign-in",
          name: "auth.sign-in",
          component: () => import("@/views/auth/sign-in.vue"),
          meta: { public: true },
        },
        {
          path: "sign-up",
          name: "auth.sign-up",
          component: () => import("@/views/auth/sign-up.vue"),
          meta: { public: true },
        },
      ],
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const id = useUserinfo()
  if (!id.isReady) {
    await id.readProfiles()
  }

  if (!to.meta.public && !id.userinfo.isLoggedIn) {
    next({ name: "auth.sign-in", query: { redirect_uri: to.fullPath } })
  } else {
    next()
  }
})

export default router
