import _ from "lodash";
import Vue from "vue";
import { defaultState, TreemapState, Treemap, TreeGraph } from "./index";

export const mutations = {
  setTreemap(state: TreemapState, treemap: Treemap) {
    if (!treemap) {
      return;
    }
    state.treemap = treemap;
  },
  setTreegraph(state: TreemapState, treegraph: TreeGraph) {
    if (!treegraph) {
      return;
    }
    state.treegraph = treegraph;
  },
  resetState(state: TreemapState) {
    Object.assign(state, defaultState());
  },
};
