import { createEffect, createMemo, createSignal, For, onMount, Show } from "solid-js";

import Cherry from "cherry-markdown";
import "cherry-markdown/dist/cherry-markdown.min.css";
import { getAtk } from "../../stores/userinfo.tsx";
import { request } from "../../scripts/request.ts";
import PostEditActions from "./PostEditActions.tsx";

export default function PostEditor(props: {
  editing?: any,
  onError: (message: string | null) => void,
  onPost: () => void
}) {
  let editorContainer: any;
  const [editor, setEditor] = createSignal<Cherry>();
  const [realmList, setRealmList] = createSignal<any[]>([]);

  const [submitting, setSubmitting] = createSignal(false);

  const [alias, setAlias] = createSignal("");
  const [publishedAt, setPublishedAt] = createSignal("");
  const [attachments, setAttachments] = createSignal<any[]>([]);
  const [categories, setCategories] = createSignal<{ alias: string, name: string }[]>([]);
  const [tags, setTags] = createSignal<{ alias: string, name: string }[]>([]);

  const theme = createMemo(() => {
    if (window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches) {
      return "dark";
    } else {
      return "light";
    }
  });

  createEffect(() => {
    editor()?.setTheme(theme());
  }, [editor(), theme()]);

  onMount(() => {
    if (editorContainer) {
      setEditor(new Cherry({
        el: editorContainer,
        value: "Welcome to the creator hub! " +
          "We provide a better editor than normal mode for you! " +
          "So you can tell us your mind clearly. " +
          "Delete this paragraph and getting start!"
      }));
    }
  });

  createEffect(() => {
    setAttachments(props.editing?.attachments ?? []);
    setCategories(props.editing?.categories ?? []);
    setTags(props.editing?.tags ?? []);
    editor()?.setValue(props.editing?.content);
  }, [props.editing]);

  async function listRealm() {
    const res = await request("/api/realms/me/available", {
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });
    if (res.status === 200) {
      setRealmList(await res.json());
    }
  }

  listRealm();

  async function doPost(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));
    if (!editor()?.getValue()) return;

    setSubmitting(true);
    const res = await request("/api/posts", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getAtk()}`
      },
      body: JSON.stringify({
        alias: alias() ? alias() : crypto.randomUUID().replace(/-/g, ""),
        title: data.title,
        content: editor()?.getValue(),
        attachments: attachments(),
        categories: categories(),
        tags: tags(),
        realm_id: parseInt(data.realm as string) !== 0 ? parseInt(data.realm as string) : undefined,
        published_at: publishedAt() ? new Date(publishedAt()) : new Date()
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
    if (!editor()?.getValue()) return;

    setSubmitting(true);
    const res = await request(`/api/posts/${props.editing?.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getAtk()}`
      },
      body: JSON.stringify({
        alias: alias() ? alias() : crypto.randomUUID().replace(/-/g, ""),
        title: data.title,
        content: editor()?.getValue(),
        attachments: attachments(),
        categories: categories(),
        tags: tags(),
        published_at: publishedAt() ? new Date(publishedAt()) : new Date()
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
  }

  return (
    <form onReset={resetForm} onSubmit={(evt) => props.editing ? doEdit(evt) : doPost(evt)}>
      <div>
        <div ref={editorContainer}></div>
      </div>

      <div class="border-y border-base-200">
        <PostEditActions
          editing={props.editing}
          onInputAlias={setAlias}
          onInputPublish={setPublishedAt}
          onInputAttachments={setAttachments}
          onInputCategories={setCategories}
          onInputTags={setTags}
          onError={props.onError}
        />
      </div>

      <div class="pt-3 pb-7 px-7">
        <Show when={!props.editing} fallback={
          <label class="form-control w-full mb-3">
            <div class="label">
              <span class="label-text">Publish region</span>
            </div>
            <input readonly type="text" class="input input-bordered"
                   value={`You published this post in realm #${props.editing?.realm_id ?? "global"}`} />
          </label>
        }>
          <label class="form-control w-full">
            <div class="label">
              <span class="label-text">Publish region</span>
            </div>
            <select name="realm" class="select select-bordered">
              <option value={0} selected>Global</option>
              <For each={realmList()}>
                {item => <option value={item.id}>{item.name}</option>}
              </For>
            </select>
            <div class="label">
              <span class="label-text-alt">Will show realms you joined or created.</span>
            </div>
          </label>
        </Show>

        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Post title</span>
          </div>
          <input value={props.editing?.title ?? ""} name="title" type="text" placeholder="Type here"
                 class="input input-bordered w-full" />
        </label>

        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Post description</span>
          </div>
          <textarea value={props.editing?.description ?? ""} disabled name="description"
                    placeholder="Not available now"
                    class="textarea textarea-bordered w-full" />
          <div class="label">
            <span class="label-text-alt">Won't display in the post list when your post is too long.</span>
          </div>
        </label>

        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Post thumbnail</span>
          </div>
          <input disabled name="thumbnail" type="file" placeholder="Not available now"
                 class="file-input file-input-bordered w-full" />
        </label>

        <button type="submit" class="btn btn-primary mt-7" disabled={submitting()}>
          <Show when={submitting()} fallback={"Submit"}>
            <span class="loading"></span>
          </Show>
        </button>
      </div>
    </form>
  );
}