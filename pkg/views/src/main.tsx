import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { CssBaseline, ThemeProvider } from "@mui/material";
import { theme } from "@/theme.ts";

import "virtual:uno.css";

import "./index.css";
import "@unocss/reset/tailwind.css";
import "@fontsource/roboto/latin.css";

import AppShell from "@/components/AppShell.tsx";
import ErrorBoundary from "@/error.tsx";
import AppLoader from "@/components/AppLoader.tsx";
import UserLayout from "@/pages/users/layout.tsx";
import { UserinfoProvider } from "@/stores/userinfo.tsx";
import { WellKnownProvider } from "@/stores/wellKnown.tsx";
import AuthLayout from "@/pages/auth/layout.tsx";
import AuthGuard from "@/pages/guard.tsx";

declare const __GARFISH_EXPORTS__: {
  provider: Object;
  registerProvider?: (provider: any) => void;
};

declare global {
  interface Window {
    __LAUNCHPAD_TARGET__?: string;
  }
}

const router = createBrowserRouter([
  {
    path: "/",
    element: <AppShell><Outlet /></AppShell>,
    errorElement: <ErrorBoundary />,
    children: [
      { path: "/", lazy: () => import("@/pages/landing.tsx") },
      {
        path: "/",
        element: <AuthGuard />,
        children: [
          {
            path: "/users",
            element: <UserLayout />,
            children: [
              { path: "/users", lazy: () => import("@/pages/users/dashboard.tsx") },
              { path: "/users/notifications", lazy: () => import("@/pages/users/notifications.tsx") },
              { path: "/users/personalize", lazy: () => import("@/pages/users/personalize.tsx") },
              { path: "/users/security", lazy: () => import("@/pages/users/security.tsx") }
            ]
          }
        ]
      }
    ]
  },
  {
    path: "/auth",
    element: <AuthLayout />,
    errorElement: <ErrorBoundary />,
    children: [
      { path: "/auth/sign-up", errorElement: <ErrorBoundary />, lazy: () => import("@/pages/auth/sign-up.tsx") },
      { path: "/auth/sign-in", errorElement: <ErrorBoundary />, lazy: () => import("@/pages/auth/sign-in.tsx") },
      { path: "/auth/o/connect", errorElement: <ErrorBoundary />, lazy: () => import("@/pages/auth/connect.tsx") }
    ]
  }
]);

const element = (
  <React.StrictMode>
    <LocalizationProvider dateAdapter={AdapterDayjs}>
      <ThemeProvider theme={theme}>
        <WellKnownProvider>
          <UserinfoProvider>
            <AppLoader>
              <CssBaseline />
              <RouterProvider router={router} />
            </AppLoader>
          </UserinfoProvider>
        </WellKnownProvider>
      </ThemeProvider>
    </LocalizationProvider>
  </React.StrictMode>
);

ReactDOM.createRoot(document.getElementById("root")!).render(element);