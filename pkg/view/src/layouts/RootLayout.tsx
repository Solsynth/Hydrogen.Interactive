import Navbar from "./shared/Navbar.tsx";
import { readProfiles, useUserinfo } from "../stores/userinfo.tsx";
import { createEffect, createMemo, createSignal, Show } from "solid-js";
import { readWellKnown } from "../stores/wellKnown.tsx";
import { BeforeLeaveEventArgs, useLocation, useNavigate, useSearchParams } from "@solidjs/router";

export default function RootLayout(props: any) {
  const [ready, setReady] = createSignal(false);

  Promise.all([readWellKnown(), readProfiles()]).then(() => setReady(true));

  const navigate = useNavigate();
  const userinfo = useUserinfo();

  const [searchParams] = useSearchParams();
  const location = useLocation();

  createEffect(() => {
    if (ready()) {
      keepGate(location.pathname + location.search);
    }
  }, [ready, userinfo]);

  function keepGate(path: string, e?: BeforeLeaveEventArgs) {
    const blacklist = ["/creator"];

    if (!userinfo?.isLoggedIn && blacklist.includes(path)) {
      if (!e?.defaultPrevented) e?.preventDefault();
      navigate(`/auth?redirect_uri=${path}`);
    }
  }

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