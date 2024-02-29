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
import LandingPage from "@/pages/landing.tsx";
import SignUpPage from "@/pages/auth/sign-up.tsx";
import SignInPage from "@/pages/auth/sign-in.tsx";
import OauthConnectPage from "@/pages/auth/connect.tsx";
import DashboardPage from "@/pages/users/dashboard.tsx";
import ErrorBoundary from "@/error.tsx";
import AppLoader from "@/components/AppLoader.tsx";
import UserLayout from "@/pages/users/layout.tsx";
import NotificationsPage from "@/pages/users/notifications.tsx";
import PersonalizePage from "@/pages/users/personalize.tsx";
import SecurityPage from "@/pages/users/security.tsx";
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
      { path: "/", element: <LandingPage /> },
      {
        path: "/",
        element: <AuthGuard />,
        children: [
          {
            path: "/users",
            element: <UserLayout />,
            children: [
              { path: "/users", element: <DashboardPage /> },
              { path: "/users/notifications", element: <NotificationsPage /> },
              { path: "/users/personalize", element: <PersonalizePage /> },
              { path: "/users/security", element: <SecurityPage /> }
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
      { path: "/auth/sign-up", element: <SignUpPage />, errorElement: <ErrorBoundary /> },
      { path: "/auth/sign-in", element: <SignInPage />, errorElement: <ErrorBoundary /> },
      { path: "/auth/o/connect", element: <OauthConnectPage />, errorElement: <ErrorBoundary /> }
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