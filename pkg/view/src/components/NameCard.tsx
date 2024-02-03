import { createSignal } from "solid-js";

import styles from "./NameCard.module.css";

export default function NameCard(props: { accountId: number, onError: (messasge: string | null) => void }) {
  const [info, setInfo] = createSignal<any>(null);

  const [_, setLoading] = createSignal(true);

  async function readInfo() {
    setLoading(true);
    const res = await fetch(`/api/users/${props.accountId}`);
    if (res.status !== 200) {
      props.onError(await res.text());
    } else {
      setInfo(await res.json());
      props.onError(null);
    }
    setLoading(false);
  }

  readInfo();

  return (
    <div class="relative">
      <figure id="banner">
        <img class="object-cover w-full h-36" src="https://images.unsplash.com/photo-1464822759023-fed622ff2c3b"
             alt="banner" />
      </figure>

      <div id="avatar" class="avatar absolute border-4 border-base-200 left-[20px] top-[4.5rem]">
        <div class="w-24">
          <img src={info()?.avatar} alt="avatar" />
        </div>
      </div>

      <div id="actions" class="flex justify-end">
        <div>
          <button type="button" class="btn btn-primary">
            <i class="fa-solid fa-plus"></i>
            Follow
          </button>
        </div>
      </div>

      <div id="description" class="px-6 pb-7">
        <h2 class="text-2xl font-bold">{info()?.name}</h2>
        <p class="text-md">{info()?.description}</p>
        <div class={`mt-2 ${styles.description}`}>
          <p class="text-xs">
            <i class="fa-solid fa-calendar-days me-2"></i>
            Joined at {new Date(info()?.created_at).toLocaleString()}
          </p>
        </div>
      </div>
    </div>
  );
}