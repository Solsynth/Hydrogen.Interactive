import { For, Match, Switch } from "solid-js";
import { clearUserinfo, useUserinfo } from "../../stores/userinfo.tsx";
import { useNavigate } from "@solidjs/router";
import { useWellKnown } from "../../stores/wellKnown.tsx";

interface MenuItem {
  icon: string;
  label: string;
  href?: string;
}

export default function Navigator() {
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
          <a href="/" class="btn btn-ghost text-xl">
            {wellKnown?.name ?? "Interactive"}
          </a>
        </div>
        <div class="navbar-center hidden md:flex">
          <ul class="menu menu-horizontal px-1">
            <For each={nav}>
              {(item) => (
                <li class="tooltip tooltip-bottom" data-tip={item.label}>
                  <a href={item.href}>
                    <i class={item.icon}></i>
                  </a>
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
              <div class="tooltip" data-tip={item.label}>
                <i class={item.icon}></i>
              </div>
            </a>
          )}
        </For>
      </div>
    </>
  );
}
