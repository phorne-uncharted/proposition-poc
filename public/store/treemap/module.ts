import { Module } from "vuex";
import { getStoreAccessors } from "vuex-typescript";
import { PropositionState } from "../store";
import { actions as moduleActions } from "./actions";
import { getters as moduleGetters } from "./getters";
import { TreemapState, state } from "./index";
import { mutations as moduleMutations } from "./mutations";

export const treemapModule: Module<TreemapState, PropositionState> = {
  getters: moduleGetters,
  actions: moduleActions,
  mutations: moduleMutations,
  state: state,
};

const { commit, read, dispatch } = getStoreAccessors<
  TreemapState,
  PropositionState
>(null);

// Typed getters
export const getters = {
  getTreemap: read(moduleGetters.getTreemap),
};

// Typed actions
export const actions = {
  fetchTreemap: dispatch(moduleActions.fetchTreemap),
};

// Typed mutations
export const mutations = {
  setTreemap: commit(moduleMutations.setTreemap),
};
