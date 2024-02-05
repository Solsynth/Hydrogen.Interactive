import { createEffect, createSignal, For, Match, Show, Switch } from "solid-js";
import { getAtk, useUserinfo } from "../stores/userinfo.tsx";

import styles from "./PostPublish.module.css";
import { closeModel, openModel } from "../scripts/modals.ts";

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
  const [uploading, setUploading] = createSignal(false);

  const [attachments, setAttachments] = createSignal<any[]>([]);
  const [categories, setCategories] = createSignal<{ alias: string, name: string }[]>([]);
  const [tags, setTags] = createSignal<{ alias: string, name: string }[]>([]);

  const [attachmentMode, setAttachmentMode] = createSignal(0);

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
        alias: data.alias ?? crypto.randomUUID().replace(/-/g, ""),
        title: data.title,
        content: data.content,
        attachments: attachments(),
        categories: categories(),
        tags: tags(),
        realm_id: props.realmId,
        published_at: data.published_at ? new Date(data.published_at as string) : new Date(),
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
    if (uploading()) return;

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
        content: data.content,
        attachments: attachments(),
        categories: categories(),
        tags: tags(),
        realm_id: props.realmId,
        published_at: data.published_at ? new Date(data.published_at as string) : new Date()
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

  async function uploadAttachment(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = new FormData(form);
    if (!data.get("attachment")) return;

    setUploading(true);
    const res = await fetch("/api/attachments", {
      method: "POST",
      headers: { "Authorization": `Bearer ${getAtk()}` },
      body: data
    });
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      const data = await res.json();
      setAttachments(attachments().concat([data.info]));
      props.onError(null);
      form.reset();
    }
    setUploading(false);
  }

  function addAttachment(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));

    setAttachments(attachments().concat([{
      ...data,
      author_id: userinfo?.profiles?.id
    }]));
    form.reset();
  }

  function addCategory(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));
    if (!data.alias) data.alias = crypto.randomUUID().replace(/-/g, "");
    if (!data.name) return;

    setCategories(categories().concat([data as any]));
    form.reset();
  }

  function removeCategory(target: any) {
    setCategories(categories().filter(item => item.alias !== target.alias));
  }

  function addTag(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
    if (!data.alias) data.alias = crypto.randomUUID().replace(/-/g, "");
    if (!data.name) return;

    setTags(tags().concat([data as any]));
    form.reset();
  }

  function removeTag(target: any) {
    setTags(tags().filter(item => item.alias !== target.alias));
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
                   placeholder="The describe for a long content (Optional)" />
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

        <textarea name="content" value={props.editing?.content ?? ""}
                  class={`${styles.publishInput} textarea w-full`}
                  placeholder="What's happend?!" />

        <div id="publish-actions" class="flex justify-between border-y border-base-200">
          <div class="flex pl-[20px]">
            <button type="button" class="btn btn-ghost w-12" onClick={() => openModel("#attachments")}>
              <i class="fa-solid fa-paperclip"></i>
            </button>
            <button type="button" class="btn btn-ghost w-12" onClick={() => openModel("#planning-publish")}>
              <i class="fa-solid fa-calendar-day"></i>
            </button>
            <button type="button" class="btn btn-ghost w-12" onClick={() => openModel("#categories-and-tags")}>
              <i class="fa-solid fa-tag"></i>
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

      <dialog id="attachments" class="modal">
        <div class="modal-box">
          <h3 class="font-bold text-lg mx-1">Attachments</h3>

          <div role="tablist" class="tabs tabs-boxed mt-3">
            <input type="radio" name="attachment" role="tab" class="tab" aria-label="File picker"
                   checked={attachmentMode() === 0} onClick={() => setAttachmentMode(0)} />
            <input type="radio" name="attachment" role="tab" class="tab" aria-label="External link"
                   checked={attachmentMode() === 1} onClick={() => setAttachmentMode(1)} />
          </div>

          <Switch>
            <Match when={attachmentMode() === 0}>
              <form class="w-full mt-2" onSubmit={uploadAttachment}>
                <label class="form-control">
                  <div class="label">
                    <span class="label-text">Pick a file</span>
                  </div>
                  <div class="join">
                    <input required type="file" name="attachment"
                           class="join-item file-input file-input-bordered w-full" />
                    <button type="submit" class="join-item btn btn-primary" disabled={uploading()}>
                      <i class="fa-solid fa-upload"></i>
                    </button>
                  </div>
                  <div class="label">
                    <span class="label-text-alt">Click upload to add this file into list</span>
                  </div>
                </label>
              </form>
            </Match>
            <Match when={attachmentMode() === 1}>
              <form class="w-full mt-2" onSubmit={addAttachment}>
                <label class="form-control">
                  <div class="label">
                    <span class="label-text">Attach an external file</span>
                  </div>
                  <div class="join">
                    <input required type="text" name="mimetype" class="join-item input input-bordered w-full"
                           placeholder="Mimetype" />
                    <input required type="text" name="filename" class="join-item input input-bordered w-full"
                           placeholder="Name" />
                  </div>
                  <div class="join">
                    <input required type="text" name="external_url" class="join-item input input-bordered w-full"
                           placeholder="External URL" />
                    <button type="submit" class="join-item btn btn-primary">
                      <i class="fa-solid fa-plus"></i>
                    </button>
                  </div>
                  <div class="label">
                    <span class="label-text-alt">Click add button to add it into list</span>
                  </div>
                </label>
              </form>
            </Match>
          </Switch>

          <Show when={attachments().length > 0}>
            <h3 class="font-bold mt-3 mx-1">Attachment list</h3>
            <ol class="mt-2 mx-1 text-sm">
              <For each={attachments()}>
                {item => <li>
                  <i class="fa-regular fa-file me-2"></i>
                  {item.filename}
                </li>}
              </For>
            </ol>
          </Show>

          <div class="modal-action">
            <button type="button" class="btn" onClick={() => closeModel("#attachments")}>Close</button>
          </div>
        </div>
      </dialog>

      <dialog id="categories-and-tags" class="modal">
        <div class="modal-box">
          <h3 class="font-bold text-lg mx-1">Categories & Tags</h3>
          <form class="w-full mt-3" onSubmit={addCategory}>
            <label class="form-control">
              <div class="label">
                <span class="label-text">Add a category</span>
              </div>
              <div class="join">
                <input type="text" name="alias" placeholder="Alias" class="join-item input input-bordered w-full" />
                <input type="text" name="name" placeholder="Name" class="join-item input input-bordered w-full" />
                <button type="submit" class="join-item btn btn-primary">
                  <i class="fa-solid fa-plus"></i>
                </button>
              </div>
              <div class="label">
                <span class="label-text-alt">
                  Alias is the url key of this category. Lowercase only, required length 4-24.
                  Leave blank for auto generate.
                </span>
              </div>
            </label>
          </form>

          <Show when={categories().length > 0}>
            <h3 class="font-bold mt-3 mx-1">Category list</h3>
            <ol class="mt-2 mx-1 text-sm">
              <For each={categories()}>
                {item => <li>
                  <i class="fa-solid fa-layer-group me-2"></i>
                  {item.name} <span class={styles.description}>#{item.alias}</span>
                  <button class="ml-2" onClick={() => removeCategory(item)}>
                    <i class="fa-solid fa-delete-left"></i>
                  </button>
                </li>}
              </For>
            </ol>
          </Show>

          <form class="w-full mt-3" onSubmit={addTag}>
            <label class="form-control">
              <div class="label">
                <span class="label-text">Add a tag</span>
              </div>
              <div class="join">
                <input type="text" name="alias" placeholder="Alias" class="join-item input input-bordered w-full" />
                <input type="text" name="name" placeholder="Name" class="join-item input input-bordered w-full" />
                <button type="submit" class="join-item btn btn-primary">
                  <i class="fa-solid fa-plus"></i>
                </button>
              </div>
              <div class="label">
                <span class="label-text-alt">
                  Alias is the url key of this tag. Lowercase only, required length 4-24.
                  Leave blank for auto generate.
                </span>
              </div>
            </label>
          </form>

          <Show when={tags().length > 0}>
            <h3 class="font-bold mt-3 mx-1">Category list</h3>
            <ol class="mt-2 mx-1 text-sm">
              <For each={tags()}>
                {item => <li>
                  <i class="fa-solid fa-tag me-2"></i>
                  {item.name} <span class={styles.description}>#{item.alias}</span>
                  <button class="ml-2" onClick={() => removeTag(item)}>
                    <i class="fa-solid fa-delete-left"></i>
                  </button>
                </li>}
              </For>
            </ol>
          </Show>

          <div class="modal-action">
            <button type="button" class="btn" onClick={() => closeModel("#categories-and-tags")}>Close</button>
          </div>
        </div>
      </dialog>
    </>
  );
}