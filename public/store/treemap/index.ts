export interface Treemap {
  name: string;
  colName: string;
  children?: Node[];
  value?: number;
}

export interface TreemapState {
  treemap: Treemap;
}

export const defaultState = (): TreemapState => {
  return {
    treemap: null,
  };
};

export const state: TreemapState = defaultState();
