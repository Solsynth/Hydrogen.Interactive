import styles from "./feed.module.css";

export default function DashboardPage(props: any) {
  return (
    <div class={`${styles.wrapper} container mx-auto`}>
      <div id="trending" class="card shadow-xl h-fit"></div>

      <div id="content max-w-[100vw]" class="card shadow-xl">
        {props.children}
      </div>

      <div id="well-known" class="card shadow-xl h-fit"></div>

    </div>
  );
}