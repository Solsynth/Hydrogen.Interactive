import Navbar from "./shared/Navbar.tsx";
import { readProfiles, useUserinfo } from "../stores/userinfo.tsx";
import { createEffect, createMemo, createSignal, Show } from "solid-js";
import { readWellKnown } from "../stores/wellKnown.tsx";
import { BeforeLeaveEventArgs, useBeforeLeave, useLocation, useNavigate, useSearchParams } from "@solidjs/router";

export default function RootLayout(props: any) {
  const [ready, setReady] = createSignal(false);

  Promise.all([readWellKnown(), readProfiles()]).then(() => setReady(true));

  const navigate = useNavigate();
  const userinfo = useUserinfo();

  const [searchParams] = useSearchParams();
  const location = useLocation();

  createEffect(() => {
    if (ready()) {
      keepGate(location.pathname);
    }
  }, [ready, userinfo]);

  function keepGate(path: string, e?: BeforeLeaveEventArgs) {
    const whitelist = ["/auth", "/auth/callback"];

    if (!userinfo?.isLoggedIn && !whitelist.includes(path)) {
      if (!e?.defaultPrevented) e?.preventDefault();
      navigate(`/auth/login?redirect_uri=${path}`);
    }
  }

  useBeforeLeave((e: BeforeLeaveEventArgs) => keepGate(e.to.toString().split("?")[0], e));

  const mainContentStyles = createMemo(() => {
    if(!searchParams["noTitle"]) {
      return "h-[calc(100vh-64px)] mt-[64px]"
    } else {
      return "h-[100vh]"
    }
  })

  return (
    <Show when={ready()} fallback={
      <div class="h-screen w-screen flex justify-center items-center">
        <div>
          <span class="loading loading-lg loading-infinity"></span>
        </div>
      </div>
    }>
      <Show when={!searchParams["noTitle"]}>
        <Navbar />
      </Show>

      <main class={mainContentStyles()}>{props.children}</main>
    </Show>
  );
}