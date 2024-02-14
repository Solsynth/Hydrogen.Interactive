import { createMemo } from "solid-js";
import { useSearchParams } from "@solidjs/router";

import styles from "./view.module.css";

export default function CreatorView(props: any) {
  const [searchParams] = useSearchParams();

  const scrollContentStyles = createMemo(() => {
    if (!searchParams["embedded"]) {
      return "max-md:mb-[64px]";
    } else {
      return "h-[100vh]";
    }
  });

  return (
    <div class={`${styles.wrapper} container mx-auto`}>
      <div id="nav" class="card shadow-xl h-fit">
        <h2 class="text-xl font-bold mt-1 py-5 px-7">Creator Hub</h2>
      </div>

      <div id="content" class={`${scrollContentStyles()} card shadow-xl`}>
        {props.children}
      </div>
    </div>
  );
}
