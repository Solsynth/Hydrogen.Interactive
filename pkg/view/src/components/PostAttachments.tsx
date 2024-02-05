import { createEffect, createMemo, createSignal, Match, Switch } from "solid-js";
import mediumZoom from "medium-zoom";

import styles from "./PostAttachments.module.css";

// @ts-ignore
import APlayer from "aplayer";
import Artplayer from "artplayer";
import HlsJs from "hls.js";
import FlvJs from "flv.js";

import "aplayer/dist/APlayer.min.css";

function Video({ url, ...rest }: any) {
  let container: any;

  function playM3u8(video: HTMLVideoElement, url: string, art: Artplayer) {
    if (HlsJs.isSupported()) {
      if (art.hls) art.hls.destroy();
      const hls = new HlsJs();
      hls.loadSource(url);
      hls.attachMedia(video);
      art.hls = hls;
      art.on('destroy', () => hls.destroy());
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = url;
    } else {
      art.notice.show = 'Unsupported playback format: m3u8';
    }
  }

  function playFlv(video: HTMLVideoElement, url: string, art: Artplayer) {
    if (FlvJs.isSupported()) {
      if (art.flv) art.flv.destroy();
      const flv = FlvJs.createPlayer({ type: 'flv', url });
      flv.attachMediaElement(video);
      flv.load();
      art.flv = flv;
      art.on('destroy', () => flv.destroy());
    } else {
      art.notice.show = 'Unsupported playback format: flv';
    }
  }

  createEffect(() => {
    new Artplayer({
      container: container as HTMLDivElement,
      url: url,
      setting: true,
      flip: true,
      loop: true,
      playbackRate: true,
      aspectRatio: true,
      subtitleOffset: true,
      fullscreen: true,
      fullscreenWeb: true,
      screenshot: true,
      autoPlayback: true,
      airplay: true,
      theme: "#49509e",
      customType: {
        m3u8: playM3u8,
        flv: playFlv,
      },
    });
  });

  return (
    <div ref={container} {...rest}></div>
  );
}

function Audio({ url, caption, ...rest }: any) {
  let container: any;

  createEffect(() => {
    new APlayer({
      container: container as HTMLDivElement,
      audio: [{
        name: caption,
        url: url,
        theme: "#49509e"
      }]
    });
  });

  return (
    <div ref={container} {...rest}></div>
  );
}


export default function PostAttachments(props: { attachments: any[] }) {
  if (props.attachments.length <= 0) return null;

  const [focus, setFocus] = createSignal(0);
  const item = createMemo(() => props.attachments[focus()]);

  function getRenderType(item: any): string {
    return item.mimetype.split("/")[0];
  }

  function getUrl(item: any): string {
    return item.external_url ? item.external_url : `/api/attachments/o/${item.file_id}`;
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
          <Match when={getRenderType(item()) === "audio"}>
            <Audio class="w-full" url={getUrl(item())} caption={item().filename} />
          </Match>
          <Match when={getRenderType(item()) === "video"}>
            <Video class="h-[360px] w-full" url={getUrl(item())} caption={item().filename} />
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