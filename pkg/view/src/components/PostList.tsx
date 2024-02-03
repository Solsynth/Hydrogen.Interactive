import { createMemo, createSignal, For, Show } from "solid-js";

import styles from "./PostList.module.css";

import PostPublish from "./PostPublish.tsx";
import PostItem from "./PostItem.tsx";

export default function PostList(props: { onError: (message: string | null) => void }) {
  const [loading, setLoading] = createSignal(true);

  const [posts, setPosts] = createSignal<any[]>([]);
  const [postCount, setPostCount] = createSignal(0);

  const [page, setPage] = createSignal(1);
  const pageCount = createMemo(() => Math.ceil(postCount() / 10));

  async function readPosts() {
    setLoading(true);
    const res = await fetch("/api/posts?" + new URLSearchParams({
      take: (10).toString(),
      offset: ((page() - 1) * 10).toString()
    }));
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      const data = await res.json();
      setPosts(data["data"]);
      setPostCount(data["count"]);
      props.onError(null);
    }
    setLoading(false);
  }

  readPosts();

  function changePage(pn: number) {
    setPage(pn);
    readPosts().then(() => {
      setTimeout(() => window.scrollTo({ top: 0, behavior: "smooth" }), 16);
    });
  }

  return (
    <div id="post-list">
      <PostPublish onPost={() => readPosts()} onError={props.onError} />

      <div id="posts">
        <For each={posts()}>
          {item => <PostItem post={item} onReact={() => readPosts()} onError={props.onError} />}
        </For>

        <div class="flex justify-center">
          <div class="join">
            <button class={`join-item btn btn-ghost ${styles.paginationControl}`} disabled={page() <= 1}
                    onClick={() => changePage(page() - 1)}>
              <i class="fa-solid fa-caret-left"></i>
            </button>
            <button class="join-item btn btn-ghost">Page {page()}</button>
            <button class={`join-item btn btn-ghost ${styles.paginationControl}`} disabled={page() >= pageCount()}
                    onClick={() => changePage(page() + 1)}>
              <i class="fa-solid fa-caret-right"></i>
            </button>
          </div>
        </div>

        <Show when={loading()}>
          <div class="w-full border-b border-base-200 pt-5 pb-7 text-center">
            <p class="loading loading-lg loading-infinity"></p>
            <p>Creating fake news...</p>
          </div>
        </Show>
      </div>
    </div>
  );
}