import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom";
import { CssBaseline, ThemeProvider } from "@mui/material";
import { theme } from "@/theme.ts";

import "virtual:uno.css";

import "./index.css";
import "@unocss/reset/tailwind.css";

import AppShell from "@/components/AppShell.tsx";
import LandingPage from "@/pages/landing.tsx";
import SignUpPage from "@/pages/auth/sign-up.tsx";
import SignInPage from "@/pages/auth/sign-in.tsx";
import ErrorBoundary from "@/error.tsx";
import AppLoader from "@/components/AppLoader.tsx";
import { UserinfoProvider } from "@/stores/userinfo.tsx";
import { WellKnownProvider } from "@/stores/wellKnown.tsx";

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
      { path: "/", element: <LandingPage /> }
    ]
  },
  { path: "/auth/sign-up", element: <SignUpPage />, errorElement: <ErrorBoundary /> },
  { path: "/auth/sign-in", element: <SignInPage />, errorElement: <ErrorBoundary /> }
]);

const element = (
  <React.StrictMode>
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
  </React.StrictMode>
);

ReactDOM.createRoot(document.getElementById("root")!).render(element);