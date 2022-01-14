import { isInteger, values } from "lodash";
import { TreemapState, Treemap } from "./index";

export const getters = {
  getTreemap(state: TreemapState): Treemap {
    return state.treemap;
  },
};
