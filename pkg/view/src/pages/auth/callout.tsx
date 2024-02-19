import { createSignal, Show } from "solid-js";
import { request } from "../../scripts/request.ts";

export default function AuthCallout() {
  const [error, setError] = createSignal<string | null>(null);
  const [status, setStatus] = createSignal("Communicating with Goatpass...");

  async function communicate() {
    const res = await request(`/api/auth${location.search}`);
    if (res.status !== 200) {
      setError(await res.text());
    } else {
      const data = await res.json();
      setStatus("Got you! Now redirecting...");
      window.open(data["target"], "_self");
    }
  }

  communicate();

  return (
    <div class="w-full h-full flex justify-center items-center">
      <div class="card w-[480px] max-w-screen shadow-xl">
        <div class="card-body">
          <div id="header" class="text-center mb-5">
            <h1 class="text-xl font-bold">Authenticate</h1>
            <p>Via your Goatpass account</p>
          </div>

          <div class="pt-16 text-center">
            <div class="text-center">
              <div>
                <span class="loading loading-lg loading-bars"></span>
              </div>
              <span>{status()}</span>
            </div>
          </div>

          <Show when={error()} fallback={<div class="mt-16"></div>}>
            <div id="alerts" class="mt-16">
              <div role="alert" class="alert alert-error">
                <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6"
                     fill="none"
                     viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="capitalize">{error()}</span>
              </div>
            </div>
          </Show>
        </div>
      </div>
    </div>
  );
}