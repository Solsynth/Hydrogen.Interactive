import { useUserinfo } from "../stores/userinfo.tsx";
import { createSignal, Show } from "solid-js";

export default function DashboardPage() {
  const userinfo = useUserinfo();

  const [error, setError] = createSignal<string | null>(null);

  function getGreeting() {
    const currentHour = new Date().getHours();

    if (currentHour >= 0 && currentHour < 12) {
      return "Good morning! Wishing you a day filled with joy and success. ☀️";
    } else if (currentHour >= 12 && currentHour < 18) {
      return "Afternoon! Hope you have a productive and joyful afternoon! ☀️";
    } else {
      return "Good evening! Wishing you a relaxing and pleasant evening. 🌙";
    }
  }

  return (
    <div class="max-w-[720px] mx-auto px-5 pt-12">
      <div id="greeting" class="px-5">
        <h1 class="text-2xl font-bold">{userinfo?.displayName}</h1>
        <p>{getGreeting()}</p>
      </div>

      <div id="alerts">
        <Show when={error()}>
          <div role="alert" class="alert alert-error mt-5">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                 viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="capitalize">{error()}</span>
          </div>
        </Show>
      </div>

    </div>
  );
}