import { getAtk, useUserinfo } from "../stores/userinfo.tsx";
import { createSignal, For, Show } from "solid-js";

import styles from "./feed.module.css";

export default function DashboardPage() {
  const userinfo = useUserinfo();

  const [error, setError] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(true);
  const [submitting, setSubmitting] = createSignal(false);

  const [posts, setPosts] = createSignal<any[]>([]);
  const [postCount, setPostCount] = createSignal(0);

  const [page, setPage] = createSignal(1);

  async function readPosts() {
    setLoading(true);
    const res = await fetch("/api/posts?" + new URLSearchParams({
      take: (10).toString(),
      skip: ((page() - 1) * 10).toString()
    }));
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setPosts(data["data"]);
      setPostCount(data["count"]);
      setError(null);
    }
    setLoading(false);
  }

  async function doPost(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));
    if (!data.content) return;

    setSubmitting(true);
    const res = await fetch("/api/posts", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getAtk()}`
      },
      body: JSON.stringify({
        alias: data.alias ?? crypto.randomUUID().replace(/-/g, ""),
        title: data.title,
        content: data.content
      })
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readPosts();
      form.reset();
      setError(null);
    }
    setSubmitting(false);
  }

  readPosts();

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

        <form id="publish" onSubmit={doPost}>
          <div id="publish-identity" class="flex border-y border-base-200">
            <div class="avatar">
              <div class="w-12">
                <Show when={userinfo?.profiles?.avatar}
                      fallback={<span class="text-3xl">{userinfo?.displayName.substring(0, 1)}</span>}>
                  <img alt="avatar" src={userinfo?.profiles?.avatar} />
                </Show>
              </div>
            </div>
            <div class="flex flex-grow">
              <input name="title" class={`${styles.publishInput} input w-full`}
                     placeholder="The describe for a long content (Optional)" />
            </div>
          </div>

          <textarea name="content" class={`${styles.publishInput} textarea w-full`}
                    placeholder="What's happend?!" />

          <div id="publish-actions" class="flex justify-between border-y border-base-200">
            <div>
              <button type="button" class="btn btn-ghost">
                <i class="fa-solid fa-paperclip"></i>
              </button>
            </div>

            <button type="submit" class="btn btn-primary" disabled={submitting()}>
              <Show when={submitting()} fallback={"Post a post"}>
                <span class="loading"></span>
              </Show>
            </button>
          </div>
        </form>

        <div id="posts">
          <Show when={!loading()} fallback={<span class="loading loading-lg loading-infinity"></span>}>
            <For each={posts()}>
              {item => <div class="post-item">

                <div class="flex bg-base-200">
                  <div class="avatar">
                    <div class="w-12">
                      <Show when={item.author.avatar}
                            fallback={<span class="text-3xl">{item.author.name.substring(0, 1)}</span>}>
                        <img alt="avatar" src={item.author.avatar} />
                      </Show>
                    </div>
                  </div>
                  <div class="flex items-center px-5">
                    <div>
                      <h3 class="font-bold text-sm">{item.author.name}</h3>
                      <p class="text-xs">{item.author.description}</p>
                    </div>
                  </div>
                </div>

                <article class="py-5 px-7">
                  <h2 class="card-title">{item.title}</h2>
                  <article class="prose">{item.content}</article>
                </article>

                <div class="grid grid-cols-4 border-y border-base-200">
                  <div class="tooltip" data-tip="Daisuki">
                    <button type="button" class="btn btn-ghost btn-block">
                      <i class="fa-solid fa-thumbs-up"></i>
                      <code class="font-mono">0</code>
                    </button>
                  </div>
                  <div class="tooltip" data-tip="Daikirai">
                    <button type="button" class="btn btn-ghost btn-block">
                      <i class="fa-solid fa-thumbs-down"></i>
                      <code class="font-mono">0</code>
                    </button>
                  </div>
                  <button type="button" class="btn btn-ghost">
                    <i class="fa-solid fa-reply"></i>
                    <span>Reply</span>
                  </button>
                  <button type="button" class="btn btn-ghost">
                    <i class="fa-solid fa-retweet"></i>
                    <span>Forward</span>
                  </button>
                </div>

              </div>}
            </For>
          </Show>
        </div>

      </div>

      <div id="well-known" class="card shadow-xl h-fit"></div>

    </div>
  );
}