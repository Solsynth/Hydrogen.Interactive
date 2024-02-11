import { createSignal, Show } from "solid-js";
import { useNavigate, useParams, useSearchParams } from "@solidjs/router";
import { createStore } from "solid-js/store";
import { closeModel, openModel } from "../scripts/modals.ts";
import PostPublish from "../components/posts/PostPublish.tsx";
import PostList from "../components/posts/PostList.tsx";
import PostItem from "../components/posts/PostItem.tsx";
import { getAtk } from "../stores/userinfo.tsx";

export default function PostPage() {
  const [error, setError] = createSignal<string | null>(null);

  const [page, setPage] = createSignal(0);
  const [related, setRelated] = createSignal<any>(null);
  const [info, setInfo] = createSignal<any>(null);

  const params = useParams();
  const navigate = useNavigate();

  const [searchParams] = useSearchParams();

  async function readPost(pn?: number) {
    if (pn) setPage(pn);
    const res = await fetch(`/api/posts/${params["postId"]}?` + new URLSearchParams({
      take: (10).toString(),
      offset: ((page() - 1) * 10).toString()
    }));
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setError(null);
      const data = await res.json();
      setInfo(data["data"]);
      setRelated({
        count: data["count"],
        data: data["related"]
      });
    }
  }

  readPost();

  async function deletePost(item: any) {
    if (!confirm(`Are you sure to delete post#${item.id}?`)) return;

    const res = await fetch(`/api/posts/${item.id}`, {
      method: "DELETE",
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      back();
      setError(null);
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

  function back() {
    if (window.history.length > 0) {
      window.history.back();
    } else {
      navigate("/");
    }
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

      <div class="flex pt-1">
        <Show when={searchParams["embedded"]} fallback={
          <button class="btn btn-ghost ml-[20px] w-12 h-12" onClick={() => back()}>
            <i class="fa-solid fa-angle-left"></i>
          </button>
        }>
          <div class="w-12 h-12 ml-[20px] flex justify-center items-center">
            <i class="fa-solid fa-comments mb-1"></i>
          </div>
        </Show>
        <div class="px-5 flex items-center">
          <p>{searchParams["title"] ?? "Post details"}</p>
        </div>
      </div>

      <dialog id="post-publish" class="modal">
        <div class="modal-box p-0 w-[540px]">
          <PostPublish
            reposting={publishMeta.reposting}
            replying={publishMeta.replying}
            editing={publishMeta.editing}
            onReset={() => setMeta(null, "none", false)}
            onError={setError}
            onPost={() => readPost()}
          />
        </div>
      </dialog>

      <Show when={info()} fallback={
        <div class="w-full border-b border-base-200 pt-5 pb-7 text-center">
          <p class="loading loading-lg loading-infinity"></p>
          <p>Creating fake news...</p>
        </div>
      }>
        <PostItem
          noClick
          post={info()}
          onError={setError}
          onReact={readPost}
          onDelete={deletePost}
          noAuthor={searchParams["noAuthor"] != null}
          noContent={searchParams["noContent"] != null}
          noControl={searchParams["noControl"] != null}
          onRepost={(item) => setMeta(item, "reposting")}
          onReply={(item) => setMeta(item, "replying")}
          onEdit={(item) => setMeta(item, "editing")}
        />

        <PostList
          noRelated
          info={related()}
          onUpdate={readPost}
          onError={setError}
          onRepost={(item) => setMeta(item, "reposting")}
          onReply={(item) => setMeta(item, "replying")}
          onEdit={(item) => setMeta(item, "editing")}
        />
      </Show>
    </>
  );
}