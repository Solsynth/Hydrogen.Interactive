import { createSignal, For, Show } from "solid-js";
import { getAtk, useUserinfo } from "../../stores/userinfo.tsx";
import { request } from "../../scripts/request.ts";
import PostAttachments from "./PostAttachments.tsx";
import * as marked from "marked";
import DOMPurify from "dompurify";

export default function PostItem(props: {
  post: any;
  noClick?: boolean;
  noAuthor?: boolean;
  noControl?: boolean;
  noRelated?: boolean;
  noContent?: boolean;
  onRepost?: (post: any) => void;
  onReply?: (post: any) => void;
  onEdit?: (post: any) => void;
  onDelete?: (post: any) => void;
  onSearch?: (filter: any) => void;
  onError: (message: string | null) => void;
  onReact: () => void;
}) {
  const [reacting, setReacting] = createSignal(false);

  const userinfo = useUserinfo();

  async function reactPost(item: any, type: string) {
    setReacting(true);
    const res = await request(`/api/posts/${item.id}/react/${type}`, {
      method: "POST",
      headers: { Authorization: `Bearer ${getAtk()}` },
    });
    if (res.status !== 201 && res.status !== 204) {
      props.onError(await res.text());
    } else {
      props.onReact();
      props.onError(null);
    }
    setReacting(false);
  }

  const content = <article class="prose" innerHTML={DOMPurify.sanitize(marked.parse(props.post.content) as string)} />;

  return (
    <div class="post-item">
      <Show when={!props.noAuthor}>
        <a href={`/accounts/${props.post.author.name}`}>
          <div class="flex bg-base-200">
            <div class="avatar pl-[20px]">
              <div class="w-12">
                <Show
                  when={props.post.author.avatar}
                  fallback={<span class="text-3xl">{props.post.author.name.substring(0, 1)}</span>}
                >
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

      <Show when={!props.noContent}>
        <div class="px-7 py-5">
          <h2 class="card-title">{props.post.title}</h2>
          <Show when={!props.noClick} fallback={content}>
            <a href={`/posts/${props.post.alias}`}>{content}</a>
          </Show>

          <div class="mt-2 flex gap-2">
            <For each={props.post.categories}>
              {(item) => (
                <a href={`/search?category=${item.alias}`} class="badge badge-primary">
                  <i class="fa-solid fa-layer-group me-1.5"></i>
                  {item.name}
                </a>
              )}
            </For>
            <For each={props.post.tags}>
              {(item) => (
                <a href={`/search?tag=${item.alias}`} class="badge badge-accent">
                  <i class="fa-solid fa-tag me-1.5"></i>
                  {item.name}
                </a>
              )}
            </For>
          </div>

          <Show when={props.post.attachments?.length > 0}>
            <div>
              <PostAttachments attachments={props.post.attachments ?? []} />
            </div>
          </Show>

          <Show when={!props.noRelated && props.post.repost_to}>
            <p class="text-xs mt-3 mb-2">
              <i class="fa-solid fa-retweet me-2"></i>
              Reposted a post
            </p>
            <div class="border border-base-200 mb-5">
              <PostItem noControl post={props.post.repost_to} onError={props.onError} onReact={props.onReact} />
            </div>
          </Show>
          <Show when={!props.noRelated && props.post.reply_to}>
            <p class="text-xs mt-3 mb-2">
              <i class="fa-solid fa-reply me-2"></i>
              Replied a post
            </p>
            <div class="border border-base-200 mb-5">
              <PostItem noControl post={props.post.reply_to} onError={props.onError} onReact={props.onReact} />
            </div>
          </Show>
        </div>
      </Show>

      <Show when={!props.noControl}>
        <div class="relative">
          <Show when={!userinfo?.isLoggedIn}>
            <div class="px-7 py-2.5 h-12 w-full opacity-0 transition-opacity hover:opacity-100 bg-base-100 border-t border-base-200 z-[1] absolute top-0 left-0">
              <b>Login!</b> To access entire platform.
            </div>
          </Show>

          <div class="grid grid-cols-3 border-y border-base-200">
            <div class="max-md:col-span-2 md:col-span-1 grid grid-cols-2">
              <div class="tooltip" data-tip="Daisuki">
                <button
                  type="button"
                  class="btn btn-ghost btn-block"
                  disabled={reacting()}
                  onClick={() => reactPost(props.post, "like")}
                >
                  <i class="fa-solid fa-thumbs-up"></i>
                  <code class="font-mono">{props.post.like_count}</code>
                </button>
              </div>

              <div class="tooltip" data-tip="Daikirai">
                <button
                  type="button"
                  class="btn btn-ghost btn-block"
                  disabled={reacting()}
                  onClick={() => reactPost(props.post, "dislike")}
                >
                  <i class="fa-solid fa-thumbs-down"></i>
                  <code class="font-mono">{props.post.dislike_count}</code>
                </button>
              </div>
            </div>

            <div class="max-md:col-span-1 md:col-span-2 flex justify-end">
              <section class="max-md:hidden">
                <div class="tooltip" data-tip="Reply">
                  <button
                    type="button"
                    class="indicator btn btn-ghost btn-block"
                    onClick={() => props.onReply && props.onReply(props.post)}
                  >
                    <span class="indicator-item badge badge-sm badge-primary">{props.post.reply_count}</span>
                    <i class="fa-solid fa-reply"></i>
                  </button>
                </div>

                <div class="tooltip" data-tip="Repost">
                  <button
                    type="button"
                    class="indicator btn btn-ghost btn-block"
                    onClick={() => props.onRepost && props.onRepost(props.post)}
                  >
                    <span class="indicator-item badge badge-sm badge-secondary">{props.post.repost_count}</span>
                    <i class="fa-solid fa-retweet"></i>
                  </button>
                </div>
              </section>

              <div class="dropdown dropdown-end">
                <div tabIndex="0" role="button" class="btn btn-ghost w-12">
                  <i class="fa-solid fa-ellipsis-vertical"></i>
                </div>
                <ul tabIndex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                  <li class="md:hidden">
                    <a class="flex justify-between" onClick={() => props.onReply && props.onReply(props.post)}>
                      <span>Reply</span>
                      <span class="badge badge-primary">{props.post.reply_count}</span>
                    </a>
                  </li>
                  <li class="md:hidden">
                    <a class="flex justify-between" onClick={() => props.onRepost && props.onRepost(props.post)}>
                      <span>Repost</span>
                      <span class="badge badge-secondary">{props.post.repost_count}</span>
                    </a>
                  </li>
                  <Show when={userinfo?.profiles?.id === props.post.author_id}>
                    <li>
                      <a onClick={() => props.onDelete && props.onDelete(props.post)}>Delete</a>
                    </li>
                  </Show>
                  <Show when={userinfo?.profiles?.id === props.post.author_id}>
                    <li>
                      <a onClick={() => props.onEdit && props.onEdit(props.post)}>Edit</a>
                    </li>
                  </Show>
                  <li>
                    <a>Report</a>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </Show>
    </div>
  );
}
