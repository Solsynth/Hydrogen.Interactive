import { createSignal, Show } from "solid-js";

import PostList from "../components/posts/PostList.tsx";
import PostPublish from "../components/posts/PostPublish.tsx";
import { createStore } from "solid-js/store";

export default function DashboardPage() {
  const [error, setError] = createSignal<string | null>(null);

  const [page, setPage] = createSignal(0);
  const [info, setInfo] = createSignal<any>(null);

  async function readPosts(pn?: number) {
    if (pn) setPage(pn);
    const res = await fetch("/api/posts?" + new URLSearchParams({
      take: (10).toString(),
      offset: ((page() - 1) * 10).toString()
    }));
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setError(null);
      setInfo(await res.json());
    }
  }

  function setMeta(data: any, field: string, scroll = true) {
    const meta: { [id: string]: any } = {
      reposting: null,
      replying: null,
      editing: null
    };
    meta[field] = data;
    setPublishMeta(meta);

    if (scroll) window.scroll({ top: 0, behavior: "smooth" });
  }

  const [publishMeta, setPublishMeta] = createStore<any>({
    replying: null,
    reposting: null,
    editing: null
  });

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

      <PostPublish
        replying={publishMeta.replying}
        reposting={publishMeta.reposting}
        editing={publishMeta.editing}
        onReset={() => setMeta(null, "none", false)}
        onPost={() => readPosts()}
        onError={setError}
      />

      <PostList
        info={info()}
        onUpdate={readPosts}
        onError={setError}
        onRepost={(item) => setMeta(item, "reposting")}
        onReply={(item) => setMeta(item, "replying")}
        onEdit={(item) => setMeta(item, "editing")}
      />
    </>
  );
}