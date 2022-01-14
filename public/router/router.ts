import Vue from "vue";
import VueRouter from "vue-router";
import Treemap from "../views/Treemap.vue";

Vue.use(VueRouter);

const router = new VueRouter({
  routes: [
    { path: "/", redirect: "/treemap" },
    { path: "/treemap", component: Treemap },
  ],
});

export default router;
