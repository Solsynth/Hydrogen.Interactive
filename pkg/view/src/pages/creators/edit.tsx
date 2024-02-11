import PostEdit from "../../components/posts/PostEditor.tsx";
import { useNavigate, useParams } from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import { getAtk } from "../../stores/userinfo.tsx";

export default function PublishPost() {
  const navigate = useNavigate();
  const params = useParams();

  const [error, setError] = createSignal<string | null>(null);
  const [post, setPost] = createSignal<any>();

  async function readPost() {
    const res = await fetch(`/api/creators/posts/${params["postId"]}`, {
      headers: { "Authorization": `Bearer ${getAtk()}` }
    });
    if (res.status === 200) {
      setPost((await res.json())["data"]);
    } else {
      setError(await res.text());
    }
  }

  readPost();

  return (
    <>
      <div class="flex pt-1 border-b border-base-200">
        <a class="btn btn-ghost ml-[20px] w-12 h-12" href="/creators">
          <i class="fa-solid fa-angle-left"></i>
        </a>
        <div class="px-5 flex items-center">
          <p>Edit「{post()?.title ? post()?.title : "Untitled"}」</p>
        </div>
      </div>

      <div id="alerts">
        <Show when={error()}>
          <div role="alert" class="alert alert-error">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
                 viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="capitalize">{error()}</span>
          </div>
        </Show>
      </div>

      <PostEdit
        editing={post()}
        onError={setError}
        onPost={() => navigate("/creators")}
      />
    </>
  );
}