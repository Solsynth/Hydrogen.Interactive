import { createMemo } from "solid-js";
import { useSearchParams } from "@solidjs/router";

import styles from "./view.module.css";

export default function FeedView(props: any) {
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
      <div id="trending" class="card shadow-xl h-fit"></div>

      <div id="content" class={`${scrollContentStyles()} card shadow-xl`}>
        {props.children}
      </div>

      <div id="well-known" class="card shadow-xl h-fit"></div>
    </div>
  );
}
