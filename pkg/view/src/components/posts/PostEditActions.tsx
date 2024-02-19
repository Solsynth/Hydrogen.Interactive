import { closeModel, openModel } from "../../scripts/modals.ts";
import { createSignal, For, Match, Show, Switch } from "solid-js";
import { getAtk, useUserinfo } from "../../stores/userinfo.tsx";
import { request } from "../../scripts/request.ts";

import styles from "./PostPublish.module.css";

export default function PostEditActions(props: {
  editing?: any;
  onInputAlias: (value: string) => void;
  onInputPublish: (value: string) => void;
  onInputAttachments: (value: any[]) => void;
  onInputCategories: (categories: any[]) => void;
  onInputTags: (tags: any[]) => void;
  onError: (message: string | null) => void;
}) {
  const userinfo = useUserinfo();

  const [uploading, setUploading] = createSignal(false);

  const [attachments, setAttachments] = createSignal<any[]>(props.editing?.attachments ?? []);
  const [categories, setCategories] = createSignal<{ alias: string; name: string }[]>(props.editing?.categories ?? []);
  const [tags, setTags] = createSignal<{ alias: string; name: string }[]>(props.editing?.tags ?? []);

  const [availableCategories, setAvailableCategories] = createSignal<any[]>([]);
  const [attachmentMode, setAttachmentMode] = createSignal(0);

  async function readCategories() {
    const res = await request("/api/categories");
    if (res.status === 200) {
      setAvailableCategories(await res.json());
    }
  }

  readCategories();

  async function uploadAttachment(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = new FormData(form);
    if (!data.get("attachment")) return;

    setUploading(true);
    const res = await request("/api/attachments", {
      method: "POST",
      headers: { Authorization: `Bearer ${getAtk()}` },
      body: data,
    });
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      const data = await res.json();
      setAttachments(attachments().concat([data.info]));
      props.onInputAttachments(attachments());
      props.onError(null);
      form.reset();
    }
    setUploading(false);
  }

  function addAttachment(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));

    setAttachments(
      attachments().concat([
        {
          ...data,
          author_id: userinfo?.profiles?.id,
        },
      ]),
    );
    props.onInputAttachments(attachments());
    form.reset();
  }

  function removeAttachment(idx: number) {
    const data = attachments().slice();
    data.splice(idx, 1);
    setAttachments(data);
    props.onInputAttachments(attachments());
  }

  function addCategory(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));
    if (!data.category) return;

    const item = availableCategories().find((item) => item.alias === data.category);

    setCategories(categories().concat([item]));
    props.onInputCategories(categories());
    form.reset();
  }

  function removeCategory(idx: number) {
    const data = categories().slice();
    data.splice(idx, 1);
    setCategories(data);
    props.onInputCategories(categories());
  }

  function addTag(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(evt.target as HTMLFormElement));
    if (!data.alias) data.alias = crypto.randomUUID().replace(/-/g, "");
    if (!data.name) return;

    setTags(tags().concat([data as any]));
    props.onInputTags(tags());
    form.reset();
  }

  function removeTag(idx: number) {
    const data = tags().slice();
    data.splice(idx, 1);
    setTags(data);
    props.onInputTags(tags());
  }

  return (
    <>
      <div class="flex pl-[20px]">
        <button type="button" class="btn btn-ghost w-12" onClick={() => openModel("#alias")}>
          <i class="fa-solid fa-link"></i>
        </button>
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

      <dialog id="alias" class="modal">
        <div class="modal-box">
          <h3 class="font-bold text-lg mx-1">Permalink</h3>
          <label class="form-control w-full mt-3">
            <div class="label">
              <span class="label-text">Alias</span>
            </div>
            <input
              name="alias"
              type="text"
              placeholder="Type here"
              class="input input-bordered w-full"
              value={props.editing?.alias ?? ""}
              onInput={(evt) => props.onInputAlias(evt.target.value)}
            />
            <div class="label">
              <span class="label-text-alt">Leave blank to generate a random string.</span>
            </div>
          </label>
          <div class="modal-action">
            <button type="button" class="btn" onClick={() => closeModel("#alias")}>
              Close
            </button>
          </div>
        </div>
      </dialog>

      <dialog id="planning-publish" class="modal">
        <div class="modal-box">
          <h3 class="font-bold text-lg mx-1">Planning Publish</h3>
          <label class="form-control w-full mt-3">
            <div class="label">
              <span class="label-text">Published At</span>
            </div>
            <input
              name="published_at"
              type="datetime-local"
              placeholder="Pick a date"
              class="input input-bordered w-full"
              value={props.editing?.published_at ?? ""}
              onInput={(evt) => props.onInputAlias(evt.target.value)}
            />
            <div class="label">
              <span class="label-text-alt">
                Before this time, your post will not be visible for everyone. You can modify this plan on Creator Hub.
              </span>
            </div>
          </label>
          <div class="modal-action">
            <button type="button" class="btn" onClick={() => closeModel("#planning-publish")}>
              Close
            </button>
          </div>
        </div>
      </dialog>

      <dialog id="attachments" class="modal">
        <div class="modal-box">
          <h3 class="font-bold text-lg mx-1">Attachments</h3>

          <div role="tablist" class="tabs tabs-boxed mt-3">
            <input
              type="radio"
              name="attachment"
              role="tab"
              class="tab"
              aria-label="File picker"
              checked={attachmentMode() === 0}
              onClick={() => setAttachmentMode(0)}
            />
            <input
              type="radio"
              name="attachment"
              role="tab"
              class="tab"
              aria-label="External link"
              checked={attachmentMode() === 1}
              onClick={() => setAttachmentMode(1)}
            />
          </div>

          <Switch>
            <Match when={attachmentMode() === 0}>
              <form class="w-full mt-2" onSubmit={uploadAttachment}>
                <label class="form-control">
                  <div class="label">
                    <span class="label-text">Pick a file</span>
                  </div>
                  <div class="join">
                    <input
                      required
                      type="file"
                      name="attachment"
                      class="join-item file-input file-input-bordered w-full"
                    />
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
                    <input
                      required
                      type="text"
                      name="mimetype"
                      class="join-item input input-bordered w-full"
                      placeholder="Mimetype"
                    />
                    <input
                      required
                      type="text"
                      name="filename"
                      class="join-item input input-bordered w-full"
                      placeholder="Name"
                    />
                  </div>
                  <div class="join">
                    <input
                      required
                      type="text"
                      name="external_url"
                      class="join-item input input-bordered w-full"
                      placeholder="External URL"
                    />
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
                {(item, idx) => (
                  <li>
                    <i class="fa-regular fa-file me-2"></i>
                    {item.filename}
                    <button class="ml-2" onClick={() => removeAttachment(idx())}>
                      <i class="fa-solid fa-delete-left"></i>
                    </button>
                  </li>
                )}
              </For>
            </ol>
          </Show>

          <div class="modal-action">
            <button type="button" class="btn" onClick={() => closeModel("#attachments")}>
              Close
            </button>
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
                <select name="category" class="join-item select select-bordered w-full">
                  <For each={availableCategories()}>{(item) => <option value={item.alias}>{item.name}</option>}</For>
                </select>
                <button type="submit" class="join-item btn btn-primary">
                  <i class="fa-solid fa-plus"></i>
                </button>
              </div>
            </label>
          </form>

          <Show when={categories().length > 0}>
            <h3 class="font-bold mt-3 mx-1">Category list</h3>
            <ol class="mt-2 mx-1 text-sm">
              <For each={categories()}>
                {(item, idx) => (
                  <li>
                    <i class="fa-solid fa-layer-group me-2"></i>
                    {item.name} <span class={styles.description}>#{item.alias}</span>
                    <button class="ml-2" onClick={() => removeCategory(idx())}>
                      <i class="fa-solid fa-delete-left"></i>
                    </button>
                  </li>
                )}
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
                  Alias is the url key of this tag. Lowercase only, required length 4-24. Leave blank for auto generate.
                </span>
              </div>
            </label>
          </form>

          <Show when={tags().length > 0}>
            <h3 class="font-bold mt-3 mx-1">Category list</h3>
            <ol class="mt-2 mx-1 text-sm">
              <For each={tags()}>
                {(item, idx) => (
                  <li>
                    <i class="fa-solid fa-tag me-2"></i>
                    {item.name} <span class={styles.description}>#{item.alias}</span>
                    <button class="ml-2" onClick={() => removeTag(idx())}>
                      <i class="fa-solid fa-delete-left"></i>
                    </button>
                  </li>
                )}
              </For>
            </ol>
          </Show>

          <div class="modal-action">
            <button type="button" class="btn" onClick={() => closeModel("#categories-and-tags")}>
              Close
            </button>
          </div>
        </div>
      </dialog>
    </>
  );
}
