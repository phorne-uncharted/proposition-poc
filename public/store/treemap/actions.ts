import axios, { AxiosResponse } from "axios";
import _ from "lodash";
import { ActionContext } from "vuex";
import { treemapGetters } from "..";
import store, { PropositionState } from "../store";
import { getters, mutations } from "./module";
import { TreemapState } from "./index";

export type TreemapContext = ActionContext<TreemapState, PropositionState>;

export const actions = {
  async fetchTreemap(
    context: TreemapContext,
    args: { url: string }
  ): Promise<void> {
    try {
      const response = await axios.post(`/site/links`, {
        url: args.url,
      });
      mutations.setTreemap(context, response.data);
    } catch (error) {
      console.error(error);
      mutations.setTreemap(context, null);
    }
  },
};
