import { createSignal, For, Show } from "solid-js";
import { getAtk, useUserinfo } from "../stores/userinfo.tsx";
import PostAttachments from "./PostAttachments.tsx";
import { SolidMarkdown } from "solid-markdown";

export default function PostItem(props: {
  post: any,
  noAuthor?: boolean,
  noControl?: boolean,
  onRepost?: (post: any) => void,
  onReply?: (post: any) => void,
  onEdit?: (post: any) => void,
  onDelete?: (post: any) => void,
  onError: (message: string | null) => void,
  onReact: () => void
}) {
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
      <Show when={!props.noAuthor}>
        <a href={`/accounts/${props.post.author.name}`}>
          <div class="flex bg-base-200">
            <div class="avatar pl-[20px]">
              <div class="w-12">
                <Show when={props.post.author.avatar}
                      fallback={<span class="text-3xl">{props.post.author.name.substring(0, 1)}</span>}>
                  <img alt="avatar" src={props.post.author.avatar} />
                </Show>
              </div>
            </div>
            <div class="flex items-center px-5">
              <div>
                <h3 class="font-bold text-sm">{props.post.author.nick}</h3>
                <p class="text-xs">{props.post.author.description}</p>
              </div>
            </div>
          </div>
        </a>
      </Show>

      <div class="px-7">
        <h2 class="card-title">{props.post.title}</h2>
        <article class="prose">
          <SolidMarkdown children={props.post.content} />
        </article>

        <div class="mt-2 flex gap-2">
          <For each={props.post.categories}>
            {item => <a class="link link-primary pb-5">
              #{item.name}
            </a>}
          </For>
          <For each={props.post.tags}>
            {item => <a class="link link-primary pb-5">
              #{item.name}
            </a>}
          </For>
        </div>

        <PostAttachments attachments={props.post.attachments ?? []} />

        <Show when={props.post.repost_to}>
          <p class="text-xs mt-3 mb-2">
            <i class="fa-solid fa-retweet me-2"></i>
            Reposted a post
          </p>
          <div class="border border-base-200 mb-5">
            <PostItem
              noControl
              post={props.post.repost_to}
              onError={props.onError}
              onReact={props.onReact}
            />
          </div>
        </Show>
        <Show when={props.post.reply_to}>
          <p class="text-xs mt-3 mb-2">
            <i class="fa-solid fa-reply me-2"></i>
            Replied a post
          </p>
          <div class="border border-base-200 mb-5">
            <PostItem
              noControl
              post={props.post.reply_to}
              onError={props.onError}
              onReact={props.onReact}
            />
          </div>
        </Show>
      </div>

      <Show when={!props.noControl}>
        <div class="relative">
          <Show when={!userinfo?.isLoggedIn}>
            <div
              class="px-7 py-2.5 h-12 w-full opacity-0 transition-opacity hover:opacity-100 bg-base-100 border-t border-base-200 z-[1] absolute top-0 left-0">
              <b>Login!</b> To access entire platform.
            </div>
          </Show>

          <div class="grid grid-cols-3 border-y border-base-200">
            <div class="grid grid-cols-2">
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
            </div>

            <div class="col-span-2 flex justify-end">
              <div class="tooltip" data-tip="Reply">
                <button type="button" class="btn btn-ghost btn-block"
                        onClick={() => props.onReply && props.onReply(props.post)}>
                  <i class="fa-solid fa-reply"></i>
                </button>
              </div>

              <div class="tooltip" data-tip="Repost">
                <button type="button" class="btn btn-ghost btn-block"
                        onClick={() => props.onRepost && props.onRepost(props.post)}>
                  <i class="fa-solid fa-retweet"></i>
                </button>
              </div>

              <div class="dropdown dropdown-end">
                <div tabIndex="0" role="button" class="btn btn-ghost w-12">
                  <i class="fa-solid fa-ellipsis-vertical"></i>
                </div>
                <ul tabIndex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                  <Show when={userinfo?.profiles?.id === props.post.author_id}>
                    <li><a onClick={() => props.onDelete && props.onDelete(props.post)}>Delete</a></li>
                  </Show>
                  <Show when={userinfo?.profiles?.id === props.post.author_id}>
                    <li><a onClick={() => props.onEdit && props.onEdit(props.post)}>Edit</a></li>
                  </Show>
                  <li><a>Report</a></li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </Show>

    </div>
  );
}