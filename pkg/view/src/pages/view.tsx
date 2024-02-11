import styles from "./view.module.css";

export default function FeedView(props: any) {
  return (
    <div class={`${styles.wrapper} container mx-auto`}>
      <div id="trending" class="card shadow-xl h-fit"></div>

      <div id="content" class="card shadow-xl">
        {props.children}
      </div>

      <div id="well-known" class="card shadow-xl h-fit"></div>

    </div>
  );
}