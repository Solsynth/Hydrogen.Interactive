import { Show } from "solid-js";

export default function Avatar(props: { user: any }) {
  return (
    <Show
      when={props.user?.avatar}
      fallback={
        <div class="avatar placeholder">
          <div class="w-12 h-12 bg-neutral text-neutral-content">
            <span class="text-xl uppercase">{props.user?.name?.substring(0, 1)}</span>
          </div>
        </div>
      }
    >
      <div class="avatar">
        <div class="w-12">
          <img alt="avatar" src={props.user?.avatar} />
        </div>
      </div>
    </Show>
  );
}