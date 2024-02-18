import { For, Match, Show, Switch } from "solid-js";
import { clearUserinfo, useUserinfo } from "../../stores/userinfo.tsx";
import { useNavigate } from "@solidjs/router";

interface MenuItem {
  label: string;
  href?: string;
  children?: MenuItem[];
}

export default function Navigatior() {
  const nav: MenuItem[] = [
    {
      label: "You", children: [
        { label: "Dashboard", href: "/" },
        { label: "Security", href: "/security" },
        { label: "Personalise", href: "/personalise" }
      ]
    }
  ];

  const userinfo = useUserinfo();
  const navigate = useNavigate();

  function logout() {
    clearUserinfo();
    navigate("/auth/login");
  }

  return (
    <div class="navbar bg-base-100 shadow-md px-5 z-10 fixed top-0">
      <div class="navbar-start">
        <a class="btn btn-ghost text-xl p-2 w-[48px] h-[48px] max-lg:ml-2.5" href="/">
          <img width="40" height="40" src="/favicon.svg" alt="Logo" />
        </a>
      </div>
      <div class="navbar-center flex">
        <ul class="menu menu-horizontal px-1">
          <For each={nav}>
            {(item) => (
              <li>
                <Show when={item.children} fallback={<a href={item.href}>{item.label}</a>}>
                  <details>
                    <summary>
                      <a href={item.href}>{item.label}</a>
                    </summary>
                    <ul class="p-2">
                      <For each={item.children}>
                        {(item) =>
                          <li>
                            <a href={item.href}>{item.label}</a>
                          </li>
                        }
                      </For>
                    </ul>
                  </details>
                </Show>
              </li>
            )}
          </For>
        </ul>
      </div>
      <div class="navbar-end pe-5">
        <Switch>
          <Match when={userinfo?.isLoggedIn}>
            <button type="button" class="btn btn-sm btn-ghost" onClick={() => logout()}>Logout</button>
          </Match>
          <Match when={!userinfo?.isLoggedIn}>
            <a href="/auth/login" class="btn btn-sm btn-primary">Login</a>
          </Match>
        </Switch>
      </div>
    </div>
  );
}
