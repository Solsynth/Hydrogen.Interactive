import { createSignal, For, Show } from "solid-js";
import { closeModel, openModel } from "../../scripts/modals.ts";
import { getAtk } from "../../stores/userinfo.tsx";
import { request } from "../../scripts/request.ts";

export default function RealmDirectoryPage() {
  const [error, setError] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const [realms, setRealms] = createSignal<any>(null);

  async function readRealms() {
    const res = await request(`/api/realms`);
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setRealms(await res.json());
    }
  }

  readRealms();

  async function createRealm(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));

    setSubmitting(true);
    const res = await request("/api/realms", {
      method: "POST",
      headers: { "Authorization": `Bearer ${getAtk()}`, "Content-Type": "application/json" },
      body: JSON.stringify({
        name: data.name,
        description: data.description,
        is_public: data.is_public != null
      })
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readRealms();
      closeModel("#create-realm");
      form.reset();
    }
    setSubmitting(false);
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

      <div class="mt-1 px-7 flex items-center justify-between">
        <h3 class="py-3 font-bold">Realms directory</h3>
        <button type="button" class="btn btn-primary" onClick={() => openModel("#create-realm")}>
          <i class="fa-solid fa-plus"></i>
        </button>
      </div>

      <For each={realms()}>
        {item => <div class="px-7 pt-7 pb-5 border-t border-base-200">
          <h2 class="text-xl font-bold">{item.name}</h2>
          <p>{item.description}</p>

          <div class="mt-2">
            <a href={`/realms/${item.id}`} class="link">Jump in</a>
          </div>
        </div>}
      </For>

      <dialog id="create-realm" class="modal">
        <div class="modal-box">
          <h2 class="card-title px-1">Create a realm</h2>
          <form class="mt-2" onSubmit={createRealm}>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Realm name</span>
              </div>
              <input name="name" type="text" placeholder="Type here" class="input input-bordered w-full" />
            </label>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Realm description</span>
              </div>
              <textarea name="description" placeholder="Type here" class="textarea textarea-bordered w-full" />
            </label>
            <div class="form-control mt-2">
              <label class="label cursor-pointer">
                <span class="label-text">Make it public</span>
                <input type="checkbox" name="is_public" class="checkbox checkbox-primary" />
              </label>
            </div>

            <button type="submit" class="btn btn-primary mt-2" disabled={submitting()}>
              <Show when={submitting()} fallback={"Submit"}>
                <span class="loading"></span>
              </Show>
            </button>
          </form>
        </div>
      </dialog>
    </>
  );
}