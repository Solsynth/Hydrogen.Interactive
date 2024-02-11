import { createEffect, createSignal, Show } from "solid-js";
import { getAtk, useUserinfo } from "../../stores/userinfo.tsx";

import styles from "./PostPublish.module.css";
import PostEditActions from "./PostEditActions.tsx";

export default function PostPublish(props: {
  replying?: any,
  reposting?: any,
  editing?: any,
  realmId?: number,
  onReset: () => void,
  onError: (message: string | null) => void,
  onPost: () => void
}) {
  const userinfo = useUserinfo();

  if (!userinfo?.isLoggedIn) {
    return (
      <div class="py-9 flex justify-center items-center">
        <div class="text-center">
          <h2 class="text-lg font-bold">Login!</h2>
          <p>Or keep silent.</p>
        </div>
      </div>
    );
  }

  const [submitting, setSubmitting] = createSignal(false);

  const [alias, setAlias] = createSignal("");
  const [publishedAt, setPublishedAt] = createSignal("");
  const [attachments, setAttachments] = createSignal<any[]>([]);
  const [categories, setCategories] = createSignal<{ alias: string, name: string }[]>([]);
  const [tags, setTags] = createSignal<{ alias: string, name: string }[]>([]);

  createEffect(() => {
    setAttachments(props.editing?.attachments ?? []);
    setCategories(props.editing?.categories ?? []);
    setTags(props.editing?.tags ?? []);
  }, [props.editing]);

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
        alias: alias() ? alias() : crypto.randomUUID().replace(/-/g, ""),
        title: data.title,
        content: data.content,
        attachments: attachments(),
        categories: categories(),
        tags: tags(),
        realm_id: data.publish_in_realm ? props.realmId : undefined,
        published_at: publishedAt() ? new Date(publishedAt()) : new Date(),
        repost_to: props.reposting?.id,
        reply_to: props.replying?.id
      })
    });
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      form.reset();
      props.onError(null);
      props.onPost();
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
        alias: alias() ? alias() : crypto.randomUUID().replace(/-/g, ""),
        title: data.title,
        content: data.content,
        attachments: attachments(),
        categories: categories(),
        tags: tags(),
        realm_id: props.realmId,
        published_at: publishedAt() ? new Date(publishedAt()) : new Date(),
      })
    });
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      form.reset();
      props.onError(null);
      props.onPost();
    }
    setSubmitting(false);
  }

  function resetForm() {
    setAttachments([]);
    setCategories([]);
    setTags([]);
    props.onReset();
  }

  return (
    <>
      <form id="publish" onSubmit={(evt) => (props.editing ? doEdit : doPost)(evt)} onReset={() => resetForm()}>
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
                   placeholder="The describe for a long content" />
          </div>
        </div>

        <Show when={props.reposting}>
          <div role="alert" class="bg-base-200 flex justify-between">
            <div class="px-5 py-3">
              <i class="fa-solid fa-circle-info me-3"></i>
              You are reposting a post from <b>{props.reposting?.author?.nick}</b>
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
              You are replying a post from <b>{props.replying?.author?.nick}</b>
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

        <Show when={props.realmId && !props.editing}>
          <div class="border-b border-base-200 px-5 h-[48px] flex items-center">
            <div class="form-control flex-grow">
              <label class="label cursor-pointer">
                <span class="label-text">Publish in this realm</span>
                <input name="publish_in_realm" type="checkbox" checked class="checkbox checkbox-primary" />
              </label>
            </div>
          </div>
        </Show>

        <textarea required name="content" value={props.editing?.content ?? ""}
                  class={`${styles.publishInput} textarea w-full`}
                  placeholder="What's happened?! (Support markdown)" />

        <div id="publish-actions" class="flex justify-between border-y border-base-200">
          <PostEditActions
            onInputAlias={setAlias}
            onInputPublish={setPublishedAt}
            onInputAttachments={setAttachments}
            onInputCategories={setCategories}
            onInputTags={setTags}
            onError={props.onError}
          />

          <div>
            <button type="submit" class="btn btn-primary" disabled={submitting()}>
              <Show when={submitting()} fallback={props.editing ? "Save changes" : "Post a post"}>
                <span class="loading"></span>
              </Show>
            </button>
          </div>
        </div>
      </form>
    </>
  );
}