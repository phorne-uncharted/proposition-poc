export interface TreeGraphItem {
  id: string;
  value: string;
}

export interface TreeGraph {
  items: TreeGraphItem[];
}

export interface Treemap {
  name: string;
  colName: string;
  children?: Node[];
  value?: number;
}

export interface TreemapState {
  treemap: Treemap;
  treegraph: TreeGraph;
}

export const defaultState = (): TreemapState => {
  return {
    treemap: null,
    treegraph: null,
  };
};

export const state: TreemapState = defaultState();
