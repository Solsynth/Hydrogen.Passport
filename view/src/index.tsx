/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import "./assets/fonts/fonts.css";
import { Route, Router } from "@solidjs/router";

import RootLayout from "./layouts/RootLayout.tsx";
import Dashboard from "./pages/dashboard.tsx";
import Login from "./pages/auth/login.tsx";

const root = document.getElementById("root");

render(() => (
  <Router root={RootLayout}>
    <Route path="/" component={Dashboard} />
    <Route path="/auth/login" component={Login} />
  </Router>
), root!);
