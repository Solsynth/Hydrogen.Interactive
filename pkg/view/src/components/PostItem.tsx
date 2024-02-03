import { createSignal, Show } from "solid-js";
import { getAtk, useUserinfo } from "../stores/userinfo.tsx";

export default function PostItem(props: { post: any, onError: (message: string | null) => void, onReact: () => void }) {
  const [reacting, setReacting] = createSignal(false);

  const userinfo = useUserinfo();

  async function reactPost(item: any, type: string) {
    setReacting(true);
    const res = await fetch(`/api/posts/${item.id}/react/${type}`, {
      method: "POST",
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });
    if (res.status !== 201 && res.status !== 204) {
      props.onError(await res.text());
    } else {
      props.onReact();
      props.onError(null);
    }
    setReacting(false);
  }

  return (
    <div class="post-item">

      <div class="flex bg-base-200">
        <div class="avatar">
          <div class="w-12">
            <Show when={props.post.author.avatar}
                  fallback={<span class="text-3xl">{props.post.author.name.substring(0, 1)}</span>}>
              <img alt="avatar" src={props.post.author.avatar} />
            </Show>
          </div>
        </div>
        <div class="flex items-center px-5">
          <div>
            <h3 class="font-bold text-sm">{props.post.author.name}</h3>
            <p class="text-xs">{props.post.author.description}</p>
          </div>
        </div>
      </div>

      <article class="py-5 px-7">
        <h2 class="card-title">{props.post.title}</h2>
        <article class="prose">{props.post.content}</article>
      </article>

      <div class="grid grid-cols-3 border-y border-base-200">

        <div class="col-span-2 grid grid-cols-4">
          <div class="tooltip" data-tip="Daisuki">
            <button type="button" class="btn btn-ghost btn-block" disabled={reacting()}
                    onClick={() => reactPost(props.post, "like")}>
              <i class="fa-solid fa-thumbs-up"></i>
              <code class="font-mono">{props.post.like_count}</code>
            </button>
          </div>

          <div class="tooltip" data-tip="Daikirai">
            <button type="button" class="btn btn-ghost btn-block" disabled={reacting()}
                    onClick={() => reactPost(props.post, "dislike")}>
              <i class="fa-solid fa-thumbs-down"></i>
              <code class="font-mono">{props.post.dislike_count}</code>
            </button>
          </div>

          <div class="tooltip" data-tip="Reply">
            <button type="button" class="btn btn-ghost btn-block">
              <i class="fa-solid fa-reply"></i>
            </button>
          </div>

          <div class="tooltip" data-tip="Repost">
            <button type="button" class="btn btn-ghost btn-block">
              <i class="fa-solid fa-retweet"></i>
            </button>
          </div>
        </div>

        <div class="flex justify-end">
          <div class="dropdown dropdown-end">
            <div tabIndex="0" role="button" class="btn btn-ghost w-12">
              <i class="fa-solid fa-ellipsis-vertical"></i>
            </div>
            <ul tabIndex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
              <Show when={userinfo?.profiles?.id === props.post.author_id}>
                <li><a>Edit</a></li>
              </Show>
              <li><a>Report</a></li>
            </ul>
          </div>
        </div>

      </div>

    </div>
  );
}