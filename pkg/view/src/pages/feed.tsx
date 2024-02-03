import { createEffect, createSignal, For, Show } from "solid-js";

import styles from "./feed.module.css";

import PostList from "../components/PostList.tsx";

export default function DashboardPage() {
  const [error, setError] = createSignal<string | null>(null);

  return (
    <div class={`${styles.wrapper} container mx-auto`}>
      <div id="trending" class="card shadow-xl h-fit"></div>

      <div id="content" class="card shadow-xl">

        <div id="alerts">
          <Show when={error()}>
            <div role="alert" class="alert alert-error">
              <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                   viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="capitalize">{error()}</span>
            </div>
          </Show>
        </div>
        
        <PostList onError={setError} />

      </div>

      <div id="well-known" class="card shadow-xl h-fit"></div>

    </div>
  );
}