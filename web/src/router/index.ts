import { createRouter, createWebHistory } from "vue-router"
import { useUserinfo } from "@/stores/userinfo"
import UserCenterLayout from "@/layouts/user-center.vue"
import AdministratorLayout from "@/layouts/administrator.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: { name: "dashboard" },
      meta: { public: true },
    },
    {
      path: "/users",
      component: UserCenterLayout,
      children: [
        {
          path: "/me",
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
          path: "/me/security",
          name: "security",
          component: () => import("@/views/security.vue"),
          meta: { title: "Your security" },
        },
      ],
    },
    {
      path: "/",
      children: [
        {
          path: "/sign-in",
          alias: ["/mfa"],
          name: "auth.sign-in",
          component: () => import("@/views/auth/sign-in.vue"),
          meta: { public: true, title: "Sign in" },
        },
        {
          path: "/sign-up",
          name: "auth.sign-up",
          component: () => import("@/views/auth/sign-up.vue"),
          meta: { public: true, title: "Sign up" },
        },
        {
          path: "/authorize",
          name: "oauth.authorize",
          component: () => import("@/views/auth/authorize.vue"),
        },
      ],
    },
    {
      path: "/flow",
      children: [
        {
          path: "confirm",
          name: "callback.confirm",
          component: () => import("@/views/flow/confirm.vue"),
          meta: { public: true, title: "Confirm registration" },
        },
        {
          path: "password-reset",
          name: "callback.password-reset",
          component: () => import("@/views/flow/password-reset.vue"),
          meta: { public: true, title: "Reset password" },
        },
      ],
    },
    {
      path: "/admin",
      component: AdministratorLayout,
      children: [
        {
          path: "",
          name: "admin.dashboard",
          component: () => import("@/views/admin/dashboard.vue"),
        },
        {
          path: "users",
          name: "admin.users",
          component: () => import("@/views/admin/users.vue"),
        },
      ]
    }
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
