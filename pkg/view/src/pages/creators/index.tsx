import { createMemo, createSignal, For, Show } from "solid-js";
import { getAtk } from "../../stores/userinfo.tsx";
import LoadingAnimation from "../../components/LoadingAnimation.tsx";
import styles from "../../components/posts/PostList.module.css";

export default function CreatorHub() {
  const [error, setError] = createSignal<string | null>(null);

  const [posts, setPosts] = createSignal<any[]>([]);
  const [postCount, setPostCount] = createSignal(0);

  const [page, setPage] = createSignal(1);
  const [loading, setLoading] = createSignal(false);

  const pageCount = createMemo(() => Math.ceil(postCount() / 10));

  async function readPosts(pn?: number) {
    if (pn) setPage(pn);
    setLoading(true);
    const res = await fetch("/api/creators/posts?" + new URLSearchParams({
      take: (10).toString(),
      offset: ((page() - 1) * 10).toString()
    }), { headers: { "Authorization": `Bearer ${getAtk()}` } });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setError(null);
      setPosts(data["data"]);
      setPostCount(data["count"]);
    }
    setLoading(false);
  }

  readPosts();

  function changePage(pn: number) {
    readPosts(pn).then(() => {
      setTimeout(() => window.scrollTo({ top: 0, behavior: "smooth" }), 16);
    });
  }

  return (
    <>
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

      <div class="mt-1 px-7 flex items-center justify-between border-b border-base-200">
        <h3 class="py-3 font-bold">Your posts</h3>
        <a class="btn btn-primary" href="/creators/publish">
          <i class="fa-solid fa-plus"></i>
        </a>
      </div>

      <div class="grid justify-items-strench">
        <For each={posts()}>
          {item =>
            <a href={`/creators/edit/${item.alias}`}>
              <div class="card sm:card-side hover:bg-base-200 transition-colors sm:max-w-none">
                <div class="card-body">
                  <Show when={item?.title} fallback={
                    <div class="line-clamp-3">
                      {item?.content?.replaceAll("#", "").replaceAll("*", "").trim()}
                    </div>
                  }>
                    <h2 class="text-xl">{item?.title}</h2>
                    <div class="mx-[-2px] mt-[-4px]">
                      {item?.categories?.map((category: any) => (
                        <span class="badge badge-primary">{category.name}</span>
                      ))}
                      {item?.tags?.map((tag: any) => (
                        <span class="badge badge-secondary">{tag.name}</span>
                      ))}
                    </div>
                    <div class="text-sm opacity-80 line-clamp-3">
                      {item?.content?.substring(0, 160).replaceAll("#", "").replaceAll("*", "").trim() + "……"}
                    </div>
                  </Show>

                  <div class="text-xs opacity-70 flex gap-2">
                    <span>Post #{item?.id}</span>
                    <span>Published at {new Date(item?.published_at).toLocaleString()}</span>
                  </div>
                </div>
              </div>
            </a>
          }
        </For>
      </div>

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
        <LoadingAnimation />
      </Show>
    </>
  );
}