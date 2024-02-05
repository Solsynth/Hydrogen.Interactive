import { createEffect, createMemo, createSignal, Match, Switch } from "solid-js";
import mediumZoom from "medium-zoom";

import styles from "./PostAttachments.module.css";

export default function PostAttachments(props: { attachments: any[] }) {
  if (props.attachments.length <= 0) return null;

  const [focus, setFocus] = createSignal(0);
  const item = createMemo(() => props.attachments[focus()]);

  function getRenderType(item: any): string {
    return item.mimetype.split("/")[0];
  }

  function getUrl(item: any): string {
    return item.external_url ?? `/api/attachments/o/${item.file_id}`;
  }

  createEffect(() => {
    mediumZoom(document.querySelectorAll(".attachment-image img"), {
      background: "var(--fallback-b1,oklch(var(--b1)/1))"
    });
  }, [focus()]);

  return (
    <>
      <p class="text-xs mt-3 mb-2">
        <i class="fa-solid fa-paperclip me-2"></i>
        Attached {props.attachments.length} file{props.attachments.length > 1 ? "s" : null}
      </p>
      <div class="border border-base-200">
        <Switch fallback={
          <div class="py-16 flex justify-center items-center">
            <div class="text-center">
              <i class="fa-solid fa-circle-question text-3xl"></i>
              <p class="mt-3">{item().filename}</p>

              <div class="flex gap-2 w-full">
                <p class="text-sm">{item().filesize <= 0 ? "Unknown" : item().filesize} Bytes</p>
                <p class="text-sm">{item().mimetype}</p>
              </div>

              <div class="mt-5">
              <a class="link" href={getUrl(item())} target="_blank">Open in browser</a>
              </div>
            </div>
          </div>
        }>
          <Match when={getRenderType(item()) === "image"}>
            <figure class="attachment-image">
              <img class="object-cover" src={getUrl(item())} alt={item().filename} />
            </figure>
          </Match>
        </Switch>

        <div id="attachments-control" class="flex justify-between border-t border-base-200">
          <div class="flex">
            <button class={`w-12 h-12 btn btn-ghost ${styles.attachmentsControl}`}
                    disabled={focus() - 1 < 0}
                    onClick={() => setFocus(focus() - 1)}>
              <i class="fa-solid fa-caret-left"></i>
            </button>
            <button class={`w-12 h-12 btn btn-ghost ${styles.attachmentsControl}`}
                    disabled={focus() + 1 >= props.attachments.length}
                    onClick={() => setFocus(focus() + 1)}>
              <i class="fa-solid fa-caret-right"></i>
            </button>
          </div>

          <div>
            <div class="h-12 px-5 py-3.5 text-sm">
              File {focus() + 1}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}