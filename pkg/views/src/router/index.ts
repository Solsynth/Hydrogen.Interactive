import { createRouter, createWebHistory } from "vue-router"
import MasterLayout from "@/layouts/master.vue"
import LandingPage from "@/views/landing.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: MasterLayout,
      children: [
        {
          path: "/",
          name: "landing",
          component: LandingPage
        }
      ]
    }
  ]
})

export default router
