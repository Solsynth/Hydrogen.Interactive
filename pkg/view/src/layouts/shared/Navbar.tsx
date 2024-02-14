import { For, Match, Show, Switch } from "solid-js";
import { clearUserinfo, useUserinfo } from "../../stores/userinfo.tsx";
import { useNavigate } from "@solidjs/router";
import { useWellKnown } from "../../stores/wellKnown.tsx";

interface MenuItem {
  icon: string;
  label: string;
  href?: string;
  children?: MenuItem[];
}

export default function Navbar() {
  const nav: MenuItem[] = [
    { icon: "fa-solid fa-pen-nib", label: "Creators", href: "/creators" },
    { icon: "fa-solid fa-newspaper", label: "Feed", href: "/" },
    { icon: "fa-solid fa-people-group", label: "Realms", href: "/realms" },
  ];

  const wellKnown = useWellKnown();
  const userinfo = useUserinfo();
  const navigate = useNavigate();

  function logout() {
    clearUserinfo();
    navigate("/auth/login");
  }

  return (
    <>
      <div class="max-md:hidden navbar bg-base-100 shadow-md px-5 z-10 h-[64px] fixed top-0">
        <div class="navbar-start">
          <div class="dropdown">
            <div tabIndex={0} role="button" class="btn btn-ghost lg:hidden">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h8m-8 6h16" />
              </svg>
            </div>
            <ul tabIndex={0} class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
              <For each={nav}>
                {(item) => (
                  <li>
                    <a href={item.href}>{item.label}</a>
                    <Show when={item.children}>
                      <ul class="p-2">
                        <For each={item.children}>
                          {(item) => (
                            <li>
                              <a href={item.href}>{item.label}</a>
                            </li>
                          )}
                        </For>
                      </ul>
                    </Show>
                  </li>
                )}
              </For>
            </ul>
          </div>
          <a href="/" class="btn btn-ghost text-xl">
            {wellKnown?.name ?? "Interactive"}
          </a>
        </div>
        <div class="navbar-center hidden lg:flex">
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
                          {(item) => (
                            <li>
                              <a href={item.href}>{item.label}</a>
                            </li>
                          )}
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
              <button type="button" class="btn btn-sm btn-ghost" onClick={() => logout()}>
                Logout
              </button>
            </Match>
            <Match when={!userinfo?.isLoggedIn}>
              <a href="/auth" class="btn btn-sm btn-primary">
                Login
              </a>
            </Match>
          </Switch>
        </div>
      </div>

      <div class="md:hidden btm-nav fixed bottom-0 bg-base-100 border-t border-base-200 z-10 h-[64px]">
        <For each={nav}>
          {(item) => (
            <a href={item.href}>
              <i class={item.icon}></i>
            </a>
          )}
        </For>
      </div>
    </>
  );
}
