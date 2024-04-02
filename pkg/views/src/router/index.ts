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
            {
              path: "/",
              name: "dashboard",
              component: () => import("@/views/dashboard.vue"),
              meta: { title: "Your account" },
            },
            {
              path: "/me/personalize",
              name: "personalize",
              component: () => import("@/views/personalize.vue"),
              meta: { title: "Your personality" },
            },
            {
              path: "/me/personal-page",
              name: "personal-page",
              component: () => import("@/views/personal-page.vue"),
              meta: { title: "Your personal page" },
            },
            {
              path: "/me/security",
              name: "security",
              component: () => import("@/views/security.vue"),
              meta: { title: "Your security" },
            },
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
          meta: { public: true, title: "Sign in" },
        },
        {
          path: "sign-up",
          name: "auth.sign-up",
          component: () => import("@/views/auth/sign-up.vue"),
          meta: { public: true, title: "Sign up" },
        },
        {
          path: "o/connect",
          name: "openid.connect",
          component: () => import("@/views/auth/connect.vue"),
        },

        {
          path: "/me/confirm",
          name: "callback.confirm",
          component: () => import("@/views/confirm.vue"),
          meta: { public: true, title: "Confirm registration" },
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

  if (to.meta.title) {
    document.title = `Solarpass | ${to.meta.title}`
  } else {
    document.title = "Solarpass"
  }

  if (!to.meta.public && !id.userinfo.isLoggedIn) {
    next({ name: "auth.sign-in", query: { redirect_uri: to.fullPath } })
  } else {
    next()
  }
})

export default router
