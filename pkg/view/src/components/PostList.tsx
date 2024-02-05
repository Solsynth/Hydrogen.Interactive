import { createMemo, createSignal, For, Show } from "solid-js";

import styles from "./PostList.module.css";
import PostItem from "./PostItem.tsx";
import { getAtk } from "../stores/userinfo.tsx";

export default function PostList(props: {
  noRelated?: boolean,
  info: { data: any[], count: number } | null,
  onRepost?: (post: any) => void,
  onReply?: (post: any) => void,
  onEdit?: (post: any) => void,
  onUpdate: (pn: number) => Promise<void>,
  onError: (message: string | null) => void
}) {
  const [loading, setLoading] = createSignal(true);

  const posts = createMemo(() => props.info?.data);
  const postCount = createMemo<number>(() => props.info?.count ?? 0);

  const [page, setPage] = createSignal(1);
  const pageCount = createMemo(() => Math.ceil(postCount() / 10));

  async function readPosts() {
    setLoading(true);
    await props.onUpdate(page());
    setLoading(false);
  }

  readPosts();

  async function deletePost(item: any) {
    if (!confirm(`Are you sure to delete post#${item.id}?`)) return;

    setLoading(true);
    const res = await fetch(`/api/posts/${item.id}`, {
      method: "DELETE",
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      await readPosts();
      props.onError(null);
    }
    setLoading(false);
  }

  function changePage(pn: number) {
    setPage(pn);
    readPosts().then(() => {
      setTimeout(() => window.scrollTo({ top: 0, behavior: "smooth" }), 16);
    });
  }

  return (
    <div id="post-list">
      <div id="posts">
        <For each={posts()}>
          {item => <PostItem
            post={item}
            noRelated={props.noRelated}
            onRepost={props.onRepost}
            onReply={props.onReply}
            onEdit={props.onEdit}
            onDelete={deletePost}
            onReact={() => readPosts()}
            onError={props.onError}
          />}
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