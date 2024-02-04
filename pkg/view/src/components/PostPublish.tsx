import { createSignal, Show } from "solid-js";
import { getAtk, useUserinfo } from "../stores/userinfo.tsx";

import styles from "./PostPublish.module.css";
import { closeModel, openModel } from "../scripts/modals.ts";

export default function PostPublish(props: {
  replying?: any,
  reposting?: any,
  editing?: any,
  onReset: () => void,
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
        published_at: data.published_at ? new Date(data.published_at as string) : new Date(),
        repost_to: props.reposting?.id,
        reply_to: props.replying?.id
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

  async function doEdit(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));
    if (!data.content) return;

    setSubmitting(true);
    const res = await fetch(`/api/posts/${props.editing?.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getAtk()}`
      },
      body: JSON.stringify({
        alias: data.alias ?? crypto.randomUUID().replace(/-/g, ""),
        title: data.title,
        content: data.content
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
    <form id="publish" onSubmit={props.editing ? doEdit : doPost} onReset={props.onReset}>
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
          <input name="title" value={props.editing?.title ?? ""}
                 class={`${styles.publishInput} input w-full`}
                 placeholder="The describe for a long content (Optional)" />
        </div>
      </div>

      <Show when={props.reposting}>
        <div role="alert" class="bg-base-200 flex justify-between">
          <div class="px-5 py-3">
            <i class="fa-solid fa-circle-info me-3"></i>
            You are reposting a post from <b>{props.reposting?.author?.name}</b>
          </div>
          <button type="reset" class="btn btn-ghost w-12" disabled={submitting()}>
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
      </Show>
      <Show when={props.replying}>
        <div role="alert" class="bg-base-200 flex justify-between">
          <div class="px-5 py-3">
            <i class="fa-solid fa-circle-info me-3"></i>
            You are replying a post from <b>{props.replying?.author?.name}</b>
          </div>
          <button type="reset" class="btn btn-ghost w-12" disabled={submitting()}>
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
      </Show>
      <Show when={props.editing}>
        <div role="alert" class="bg-base-200 flex justify-between">
          <div class="px-5 py-3">
            <i class="fa-solid fa-circle-info me-3"></i>
            You are editing a post published at{" "}
            <b>{new Date(props.editing?.created_at).toLocaleString()}</b>
          </div>
          <button type="reset" class="btn btn-ghost w-12" disabled={submitting()}>
            <i class="fa-solid fa-xmark"></i>
          </button>
        </div>
      </Show>

      <textarea name="content" value={props.editing?.content ?? ""}
                class={`${styles.publishInput} textarea w-full`}
                placeholder="What's happend?!" />

      <div id="publish-actions" class="flex justify-between border-y border-base-200">
        <div class="flex">
          <button type="button" class="btn btn-ghost w-12">
            <i class="fa-solid fa-paperclip"></i>
          </button>
          <button type="button" class="btn btn-ghost w-12" onClick={() => openModel("#planning-publish")}>
            <i class="fa-solid fa-calendar-day"></i>
          </button>
        </div>

        <div>
          <button type="submit" class="btn btn-primary" disabled={submitting()}>
            <Show when={submitting()} fallback={props.editing ? "Save changes" : "Post a post"}>
              <span class="loading"></span>
            </Show>
          </button>
        </div>
      </div>

      <dialog id="planning-publish" class="modal">
        <div class="modal-box">
          <h3 class="font-bold text-lg mx-1">Planning Publish</h3>
          <label class="form-control w-full mt-3">
            <div class="label">
              <span class="label-text">Published At</span>
            </div>
            <input name="published_at" type="datetime-local" placeholder="Pick a date"
                   class="input input-bordered w-full" />
            <div class="label">
              <span class="label-text-alt">
                Before this time, your post will not be visible for everyone.
                You can modify this plan on Creator Hub.
              </span>
            </div>
          </label>
          <div class="modal-action">
            <button type="button" class="btn" onClick={() => closeModel("#planning-publish")}>Close</button>
          </div>
        </div>
      </dialog>
    </form>
  );
}