import "solid-devtools";

/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import "./assets/fonts/fonts.css";
import { lazy } from "solid-js";
import { Route, Router } from "@solidjs/router";

import "@fortawesome/fontawesome-free/css/all.min.css";

import RootLayout from "./layouts/RootLayout.tsx";
import { UserinfoProvider } from "./stores/userinfo.tsx";
import { WellKnownProvider } from "./stores/wellKnown.tsx";

const root = document.getElementById("root");

const router = (basename?: string) => (
  <WellKnownProvider>
    <UserinfoProvider>
      <Router root={RootLayout} base={basename}>
        <Route path="/" component={lazy(() => import("./pages/dashboard.tsx"))} />
        <Route path="/security" component={lazy(() => import("./pages/security.tsx"))} />
        <Route path="/personalise" component={lazy(() => import("./pages/personalise.tsx"))} />
        <Route path="/auth/login" component={lazy(() => import("./pages/auth/login.tsx"))} />
        <Route path="/auth/register" component={lazy(() => import("./pages/auth/register.tsx"))} />
        <Route path="/auth/o/connect" component={lazy(() => import("./pages/auth/connect.tsx"))} />
        <Route path="/auth/o/callback" component={lazy(() => import("./pages/auth/callback.tsx"))} />
        <Route path="/users/me/confirm" component={lazy(() => import("./pages/users/confirm.tsx"))} />
      </Router>
    </UserinfoProvider>
  </WellKnownProvider>
);

declare const __GARFISH_EXPORTS__: {
  provider: Object;
  registerProvider?: (provider: any) => void;
};

declare global {
  interface Window {
    __GARFISH__: boolean;
    __LAUNCHPAD_TARGET__?: string;
  }
}

export const provider = () => ({
  render: ({ dom, basename }: { dom: any, basename: string }) => {
    render(
      () => router(basename),
      dom.querySelector("#root")
    );
  },
  destroy: () => {
  }
});

if (!window.__GARFISH__) {
  console.log("Running directly!")
  render(router, root!);
} else if (typeof __GARFISH_EXPORTS__ !== "undefined") {
  console.log("Running in launchpad container!")
  console.log("Launchpad target:", window.__LAUNCHPAD_TARGET__)
  if (__GARFISH_EXPORTS__.registerProvider) {
    __GARFISH_EXPORTS__.registerProvider(provider);
  } else {
    __GARFISH_EXPORTS__.provider = provider;
  }
}
