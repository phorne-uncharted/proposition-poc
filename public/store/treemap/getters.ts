import { isInteger, values } from "lodash";
import { TreemapState, Treemap, TreeGraph } from "./index";

export const getters = {
  getTreemap(state: TreemapState): Treemap {
    return state.treemap;
  },
  getTreegraph(state: TreemapState): TreeGraph {
    return state.treegraph;
  },
};
