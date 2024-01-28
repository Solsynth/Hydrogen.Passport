import "solid-devtools";

/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import "./assets/fonts/fonts.css";
import { lazy } from "solid-js";
import { Route, Router } from "@solidjs/router";

import RootLayout from "./layouts/RootLayout.tsx";

const root = document.getElementById("root");

render(() => (
  <Router root={RootLayout}>
    <Route path="/" component={lazy(() => import("./pages/dashboard.tsx"))} />
    <Route path="/auth/login" component={lazy(() => import("./pages/auth/login.tsx"))} />
    <Route path="/auth/register" component={lazy(() => import("./pages/auth/register.tsx"))} />
    <Route path="/users/me/confirm" component={lazy(() => import("./pages/users/confirm.tsx"))} />
  </Router>
), root!);
