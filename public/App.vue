<template>
  <div id="proposition-app">
    <router-view ref="view" class="view" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import VueRouterSync from "vuex-router-sync";
import VueObserveVisibility from "vue-observe-visibility";
import BootstrapVue from "bootstrap-vue";
import store from "./store/store";
import { actions as treemapActions } from "./store/treemap/module";
import vSelect from "vue-select";
import router from "./router/router";
import "bootstrap-vue/dist/bootstrap-vue.css";
import "./styles/uncharted-bootstrap-v4.5-custom.css";
import "./styles/main.css";
import "vue-select/dist/vue-select.css";

Vue.component("v-select", vSelect);
Vue.use(BootstrapVue);
Vue.use(VueObserveVisibility);
Vue.config.performance = true;
// sync store and router
VueRouterSync.sync(store, router, { moduleName: "routeModule" });
// main app component
export default Vue.extend({
  store: store,
  router: router,
});
</script>

<style>
/*
  This is global css.
*/
/*
pulse is used for hints
*/
.pulse {
  overflow: visible;
  position: relative;
}
.pulse:before {
  content: "";
  display: block;
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  background-color: inherit;
  border-radius: inherit;
  transition: opacity 0.3s, transform 0.3s;
  animation: pulse-animation 1s cubic-bezier(0.24, 0, 0.38, 1) infinite;
  z-index: -1;
  animation-iteration-count: 1;
}
.z-index-1 {
  z-index: 1;
}
@keyframes pulse-animation {
  0% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0;
    transform: scale(1.5);
  }
  100% {
    opacity: 0;
    transform: scale(1.5);
  }
}
</style>
