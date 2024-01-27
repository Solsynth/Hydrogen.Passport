import "solid-devtools";

/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import "./assets/fonts/fonts.css";
import { Route, Router } from "@solidjs/router";

import RootLayout from "./layouts/RootLayout.tsx";
import DashboardPage from "./pages/dashboard.tsx";
import LoginPage from "./pages/auth/login.tsx";
import RegisterPage from "./pages/auth/register.tsx";

const root = document.getElementById("root");

render(() => (
  <Router root={RootLayout}>
    <Route path="/" component={DashboardPage} />
    <Route path="/auth/login" component={LoginPage} />
    <Route path="/auth/register" component={RegisterPage} />
  </Router>
), root!);
