import Vue from "vue";
import Vuex, { Store } from "vuex";
import { Route } from "vue-router";
import { treemapModule } from "./treemap/module";
import { TreemapState } from "./treemap/index";

Vue.use(Vuex);

export interface PropositionState {
  treemapModule: TreemapState;
}

const store = new Store<PropositionState>({
  modules: {
    treemapModule,
  },
  strict: true,
});

export default store;
