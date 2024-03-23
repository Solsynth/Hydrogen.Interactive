import { createRouter, createWebHistory } from "vue-router"
import MasterLayout from "@/layouts/master.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: MasterLayout,
      children: [
        {
          path: "/",
          name: "explore",
          component: () => import("@/views/explore.vue")
        },

        {
          path: "/p/moments/:alias",
          name: "posts.details.moments",
          component: () => import("@/views/posts/moments.vue")
        },
        {
          path: "/p/articles/:alias",
          name: "posts.details.articles",
          component: () => import("@/views/posts/articles.vue")
        },

        {
          path: "/realms/:realmId",
          name: "realms.page",
          component: () => import("@/views/realms/page.vue")
        }
      ]
    }
  ]
})

export default router
