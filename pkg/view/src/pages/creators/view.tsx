import styles from "./view.module.css";

export default function CreatorView(props: any) {
  return (
    <div class={`${styles.wrapper} container mx-auto`}>
      <div id="nav" class="card shadow-xl h-fit">
        <h2 class="text-xl font-bold mt-1 py-5 px-7">Creator Hub</h2>
      </div>

      <div id="content" class="card shadow-xl">
        {props.children}
      </div>

    </div>
  );
}