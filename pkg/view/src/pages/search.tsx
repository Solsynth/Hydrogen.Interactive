import { useNavigate, useSearchParams } from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import { createStore } from "solid-js/store";
import PostPublish from "../components/posts/PostPublish.tsx";
import PostList from "../components/posts/PostList.tsx";
import { closeModel, openModel } from "../scripts/modals.ts";

export default function SearchPage() {
  const [searchParams] = useSearchParams();

  const [error, setError] = createSignal<string | null>(null);

  const [page, setPage] = createSignal(0);
  const [info, setInfo] = createSignal<any>(null);

  const navigate = useNavigate();

  async function readPosts(pn?: number) {
    if (pn) setPage(pn);
    const res = await request("/api/posts?" + new URLSearchParams({
      take: (10).toString(),
      offset: ((page() - 1) * 10).toString(),
      ...searchParams
    }));
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setError(null);
      setInfo(await res.json());
    }
  }

  function setMeta(data: any, field: string, open = true) {
    const meta: { [id: string]: any } = {
      reposting: null,
      replying: null,
      editing: null
    };
    meta[field] = data;
    setPublishMeta(meta);

    if (open) openModel("#post-publish");
    else closeModel("#post-publish");
  }

  const [publishMeta, setPublishMeta] = createStore<any>({
    replying: null,
    reposting: null,
    editing: null
  });

  function getDescribe() {
    let builder = [];
    if (searchParams["category"]) {
      builder.push("category is #" + searchParams["category"]);
    } else if (searchParams["tag"]) {
      builder.push("tag is #" + searchParams["tag"]);
    }

    return builder.join(" and ");
  }

  function back() {
    if (window.history.length > 0) {
      window.history.back();
    } else {
      navigate("/");
    }
  }

  return (
    <>
      <div class="flex pt-1">
        <button class="btn btn-ghost ml-[20px] w-12 h-12" onClick={() => back()}>
          <i class="fa-solid fa-angle-left"></i>
        </button>
        <div class="px-5 flex items-center">
          <p>Search</p>
        </div>
      </div>

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

      <div role="alert" class="alert alert-info px-[20px]">
        <i class="fa-solid fa-magnifying-glass pl-[13px]"></i>
        <span>You will only see posts with <b>{getDescribe()}</b></span>
      </div>

      <dialog id="post-publish" class="modal">
        <div class="modal-box p-0 w-[540px]">
          <PostPublish
            reposting={publishMeta.reposting}
            replying={publishMeta.replying}
            editing={publishMeta.editing}
            onReset={() => setMeta(null, "none", false)}
            onError={setError}
            onPost={() => readPosts()}
          />
        </div>
      </dialog>

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