import { createSignal, Show } from "solid-js";
import { getAtk, useUserinfo } from "../stores/userinfo.tsx";

import styles from "./PostPublish.module.css";

export default function PostPublish(props: {
  replying?: any,
  reposting?: any,
  editing?: any,
  onError: (message: string | null) => void,
  onPost: () => void
}) {
  const userinfo = useUserinfo();

  const [submitting, setSubmitting] = createSignal(false);

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
        content: data.content,
        repost_to: props.reposting?.id,
        reply_to: props.replying?.id,
      })
    });
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      form.reset();
      props.onPost();
      props.onError(null);
    }
    setSubmitting(false);
  }

  return (
    <form id="publish" onSubmit={doPost}>
      <div id="publish-identity" class="flex border-y border-base-200">
        <div class="avatar pl-[20px]">
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

      <Show when={props.reposting}>
        <div role="alert" class="px-5 py-3 bg-base-200">
          <i class="fa-solid fa-circle-info me-3"></i>
          You are reposting a post from <b>{props.reposting?.author?.name}</b>
        </div>
      </Show>
      <Show when={props.replying}>
        <div role="alert" class="px-5 py-3 bg-base-200">
          <i class="fa-solid fa-circle-info me-3"></i>
          You are replying a post from <b>{props.replying?.author?.name}</b>
        </div>
      </Show>
      <Show when={props.editing}>
        <div role="alert" class="px-5 py-3 bg-base-200">
          <i class="fa-solid fa-circle-info me-3"></i>
          You are editing a post published at{" "}
          <b>{new Date(props.editing?.created_at).toLocaleString()}</b>
        </div>
      </Show>

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
  );
}