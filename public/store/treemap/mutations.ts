import _ from "lodash";
import Vue from "vue";
import { defaultState, TreemapState, Treemap } from "./index";

export const mutations = {
  setTreemap(state: TreemapState, treemap: Treemap) {
    if (!treemap) {
      return;
    }
    state.treemap = treemap;
  },
  resetState(state: TreemapState) {
    Object.assign(state, defaultState());
  },
};
