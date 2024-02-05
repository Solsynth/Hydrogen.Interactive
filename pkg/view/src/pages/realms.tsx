import { createSignal, For, Show } from "solid-js";

export default function RealmDirectoryPage() {
  const [error, setError] = createSignal<string | null>(null);

  const [realms, setRealms] = createSignal<any>(null);

  async function readRealms() {
    const res = await fetch(`/api/realms`);
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      setRealms(await res.json());
    }
  }

  readRealms();

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

      <For each={realms()}>
        {item => <div class="px-7 pt-7 pb-5 border-t border-base-200">
          <h2 class="text-xl font-bold">{item.name}</h2>
          <p>{item.description}</p>

          <div class="mt-2">
            <a href={`/realms/${item.id}`} class="link">Jump in</a>
          </div>
        </div>}
      </For>
    </>
  );
}