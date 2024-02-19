import "solid-devtools";

/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import "./assets/fonts/fonts.css";
import { lazy } from "solid-js";
import { Route, Router } from "@solidjs/router";

import "@fortawesome/fontawesome-free/css/all.css";

import RootLayout from "./layouts/RootLayout.tsx";
import FeedView from "./pages/view.tsx";
import Global from "./pages/global.tsx";
import PostReference from "./pages/post.tsx";
import CreatorView from "./pages/creators/view.tsx";
import { UserinfoProvider } from "./stores/userinfo.tsx";
import { WellKnownProvider } from "./stores/wellKnown.tsx";

const root = document.getElementById("root");

const router = (basename?: string) => (
  <WellKnownProvider>
    <UserinfoProvider>
      <Router root={RootLayout} base={basename}>
        <Route path="/" component={FeedView}>
          <Route path="/" component={Global} />
          <Route path="/posts/:postId" component={PostReference} />
          <Route path="/search" component={lazy(() => import("./pages/search.tsx"))} />
          <Route path="/realms" component={lazy(() => import("./pages/realms"))} />
          <Route path="/realms/:realmId" component={lazy(() => import("./pages/realms/realm.tsx"))} />
          <Route path="/accounts/:accountId" component={lazy(() => import("./pages/account.tsx"))} />
        </Route>
        <Route path="/creators" component={CreatorView}>
          <Route path="/" component={lazy(() => import("./pages/creators"))} />
          <Route path="/publish" component={lazy(() => import("./pages/creators/publish.tsx"))} />
          <Route path="/edit/:postId" component={lazy(() => import("./pages/creators/edit.tsx"))} />
        </Route>
        <Route path="/auth" component={lazy(() => import("./pages/auth/callout.tsx"))} />
        <Route path="/auth/callback" component={lazy(() => import("./pages/auth/callback.tsx"))} />
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
  if (__GARFISH_EXPORTS__.registerProvider) {
    __GARFISH_EXPORTS__.registerProvider(provider);
  } else {
    __GARFISH_EXPORTS__.provider = provider;
  }
}