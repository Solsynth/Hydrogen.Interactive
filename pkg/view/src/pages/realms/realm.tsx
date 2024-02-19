import { createSignal, Show } from "solid-js";
import { createStore } from "solid-js/store";
import { useNavigate, useParams } from "@solidjs/router";

import PostList from "../../components/posts/PostList.tsx";
import PostPublish from "../../components/posts/PostPublish.tsx";

import styles from "./realm.module.css";
import { getAtk, useUserinfo } from "../../stores/userinfo.tsx";
import { closeModel, openModel } from "../../scripts/modals.ts";

export default function RealmPage() {
  const userinfo = useUserinfo();

  const [error, setError] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const [realm, setRealm] = createSignal<any>(null);
  const [page, setPage] = createSignal(0);
  const [info, setInfo] = createSignal<any>(null);

  const params = useParams();
  const navigate = useNavigate();

  async function readRealm() {
    const res = await request(`/api/realms/${params["realmId"]}`);
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setRealm(await res.json());
    }
  }

  readRealm();

  async function readPosts(pn?: number) {
    if (pn) setPage(pn);
    const res = await request(`/api/posts?` + new URLSearchParams({
      take: (10).toString(),
      offset: ((page() - 1) * 10).toString(),
      realmId: params["realmId"]
    }));
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setError(null);
      setInfo(await res.json());
    }
  }

  async function editRealm(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));

    setSubmitting(true);
    const res = await request(`/api/realms/${params["realmId"]}`, {
      method: "PUT",
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
      await readRealm();
      closeModel("#edit-realm");
      form.reset();
    }
    setSubmitting(false);
  }

  async function inviteMember(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));

    setSubmitting(true);
    const res = await request(`/api/realms/${params["realmId"]}/invite`, {
      method: "POST",
      headers: { "Authorization": `Bearer ${getAtk()}`, "Content-Type": "application/json" },
      body: JSON.stringify(data)
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readRealm();
      closeModel("#invite-member");
      form.reset();
    }
    setSubmitting(false);
  }

  async function kickMember(evt: SubmitEvent) {
    evt.preventDefault();

    const form = evt.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form));

    setSubmitting(true);
    const res = await request(`/api/realms/${params["realmId"]}/kick`, {
      method: "POST",
      headers: { "Authorization": `Bearer ${getAtk()}`, "Content-Type": "application/json" },
      body: JSON.stringify(data)
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      await readRealm();
      closeModel("#kick-member");
      form.reset();
    }
    setSubmitting(false);
  }

  async function breakRealm() {
    if (!confirm("Are you sure about that? All posts in this realm will disappear forever.")) return;

    const res = await request(`/api/realms/${params["realmId"]}`, {
      method: "DELETE",
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      navigate("/realms");
    }
  }

  function setMeta(data: any, field: string, scroll = true) {
    const meta: { [id: string]: any } = {
      reposting: null,
      replying: null,
      editing: null
    };
    meta[field] = data;
    setPublishMeta(meta);

    if (scroll) window.scroll({ top: 0, behavior: "smooth" });
  }

  const [publishMeta, setPublishMeta] = createStore<any>({
    replying: null,
    reposting: null,
    editing: null
  });

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

      <div class="px-7 pt-7 pb-5">
        <h2 class="text-2xl font-bold">{realm()?.name}</h2>
        <p>{realm()?.description}</p>

        <div class={`${styles.description} text-sm mt-3`}>
          <p>Realm #{realm()?.id}</p>
          <Show when={realm()?.account_id === userinfo?.profiles?.id}>
            <div class="flex gap-2">
              <button class="link" onClick={() => openModel("#edit-realm")}>Edit</button>
              <button class="link" onClick={() => openModel("#invite-member")}>Invite</button>
              <button class="link" onClick={() => openModel("#kick-member")}>Kick</button>
              <button class="link" onClick={() => breakRealm()}>Break-up</button>
            </div>
          </Show>
        </div>
      </div>

      <PostPublish
        realmId={parseInt(params["realmId"])}
        replying={publishMeta.replying}
        reposting={publishMeta.reposting}
        editing={publishMeta.editing}
        onReset={() => setMeta(null, "none", false)}
        onPost={() => readPosts()}
        onError={setError}
      />

      <PostList
        info={info()}
        onUpdate={readPosts}
        onError={setError}
        onRepost={(item) => setMeta(item, "reposting")}
        onReply={(item) => setMeta(item, "replying")}
        onEdit={(item) => setMeta(item, "editing")}
      />

      <dialog id="edit-realm" class="modal">
        <div class="modal-box">
          <h2 class="card-title px-1">Edit your realm</h2>
          <form class="mt-2" onSubmit={editRealm}>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Realm name</span>
              </div>
              <input value={realm()?.name} name="name" type="text" placeholder="Type here"
                     class="input input-bordered w-full" />
            </label>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Realm description</span>
              </div>
              <textarea value={realm()?.description} name="description" placeholder="Type here"
                        class="textarea textarea-bordered w-full" />
            </label>
            <div class="form-control mt-2">
              <label class="label cursor-pointer">
                <span class="label-text">Make it public</span>
                <input checked={realm()?.is_public} type="checkbox" name="is_public"
                       class="checkbox checkbox-primary" />
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

      <dialog id="invite-member" class="modal">
        <div class="modal-box">
          <h2 class="card-title px-1">Invite someone as a member</h2>
          <form class="mt-2" onSubmit={inviteMember}>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Username</span>
              </div>
              <input name="account_name" type="text" placeholder="Type here" class="input input-bordered w-full" />
              <div class="label">
                <span class="label-text-alt">
                  Invite someone via their username so that they can publish content in non-public realm.
                </span>
              </div>
            </label>

            <button type="submit" class="btn btn-primary mt-2" disabled={submitting()}>
              <Show when={submitting()} fallback={"Submit"}>
                <span class="loading"></span>
              </Show>
            </button>
          </form>
        </div>
      </dialog>

      <dialog id="kick-member" class="modal">
        <div class="modal-box">
          <h2 class="card-title px-1">Kick someone out of your realm</h2>
          <form class="mt-2" onSubmit={kickMember}>
            <label class="form-control w-full">
              <div class="label">
                <span class="label-text">Username</span>
              </div>
              <input name="account_name" type="text" placeholder="Type here" class="input input-bordered w-full" />
              <div class="label">
                <span class="label-text-alt">
                  Remove someone out of your realm.
                </span>
              </div>
            </label>

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