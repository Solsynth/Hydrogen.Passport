import { createRouter, createWebHistory } from "vue-router"
import { useUserinfo } from "@/stores/userinfo"
import MasterLayout from "@/layouts/master.vue"
import UserCenterLayout from "@/layouts/user-center.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: MasterLayout,
      children: [
        {
          path: "/",
          component: UserCenterLayout,
          children: [
            { path: "/", name: "dashboard", component: () => import("@/views/dashboard.vue") },
            { path: "/me/personalize", name: "personalize", component: () => import("@/views/personalize.vue") },
            { path: "/me/security", name: "security", component: () => import("@/views/security.vue") },
          ],
        },
      ],
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
        {
          path: "o/connect",
          name: "openid.connect",
          component: () => import("@/views/auth/connect.vue"),
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
