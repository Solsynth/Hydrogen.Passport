import "solid-devtools";

/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import "./assets/fonts/fonts.css";
import { lazy } from "solid-js";
import { Route, Router } from "@solidjs/router";

import RootLayout from "./layouts/RootLayout.tsx";
import { UserinfoProvider } from "./stores/userinfo.tsx";
import { WellKnownProvider } from "./stores/wellKnown.tsx";

const root = document.getElementById("root");

render(() => (
  <WellKnownProvider>
    <UserinfoProvider>
      <Router root={RootLayout}>
        <Route path="/" component={lazy(() => import("./pages/dashboard.tsx"))} />
        <Route path="/security" component={lazy(() => import("./pages/security.tsx"))} />
        <Route path="/personalise" component={lazy(() => import("./pages/personalise.tsx"))} />
        <Route path="/auth/login" component={lazy(() => import("./pages/auth/login.tsx"))} />
        <Route path="/auth/register" component={lazy(() => import("./pages/auth/register.tsx"))} />
        <Route path="/auth/oauth/connect" component={lazy(() => import("./pages/auth/connect.tsx"))} />
        <Route path="/auth/oauth/callback" component={lazy(() => import("./pages/auth/callback.tsx"))} />
        <Route path="/users/me/confirm" component={lazy(() => import("./pages/users/confirm.tsx"))} />
      </Router>
    </UserinfoProvider>
  </WellKnownProvider>
), root!);
