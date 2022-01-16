<template>
  <div class="container-fluid d-flex flex-column join-view">
    <div class="">
      <form ref="urlInputForm">
        <b-form-group label-for="url-input">
          <b-form-input id="url-input" v-model="url" />
          <b-form-select v-model="selectedGraphType" size="sm">
            <b-form-select-option
              v-for="gt in graphTypes"
              :key="gt"
              :value="gt"
            >
              {{ gt }}
            </b-form-select-option>
          </b-form-select>
        </b-form-group>
      </form>
      <b-button variant="primary" @click="crawl" :disabled="isCrawling">
        <b-spinner v-if="isCrawling" small />
        <span v-else>crawl</span>
      </b-button>
    </div>
    <div class="chart">
      <svg id="treemapGraph" :height="800" :width="1200"></svg>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { Treemap } from "../store/treemap/index";
import { graphTreemap, graphTreegraph } from "../util/treemap";
import { actions, getters } from "../store/treemap/module";

export default Vue.extend({
  name: "treemap",

  data() {
    return {
      url: "https://onedemo-telco.azurewebsites.net/",
      isCrawling: false,
      selectedGraphType: "treegraph",
      graphTypes: ["treemap", "treegraph"],
    };
  },

  computed: {
    treemap(): Treemap {
      return getters.getTreemap(this.$store);
    },
  },

  methods: {
    async loadTreemap() {
      await actions.fetchTreemap(this.$store, {
        url: this.url,
        maxDepth: 15,
      });
      const treemapData = getters.getTreemap(this.$store);
      graphTreemap("#treemapGraph", treemapData);
    },
    async loadTreegraph() {
      await actions.fetchTreegraph(this.$store, {
        url: this.url,
        maxDepth: 10,
      });
      const treegraphData = getters.getTreegraph(this.$store);
      graphTreegraph("#treemapGraph", treegraphData);
    },
    async crawl() {
      this.isCrawling = true;
      switch (this.selectedGraphType) {
        case "treemap":
          await this.loadTreemap();
          break;
        case "treegraph":
          await this.loadTreegraph();
          break;
      }
      this.isCrawling = false;
    },
  },

  async beforeMount() {
    this.crawl();
  },
});
</script>

<style>
.chart {
  height: 600px;
  overflow: scroll;
}
.header-label {
  padding: 1rem 0 0.5rem 0;
  font-weight: bold;
}

.node circle {
  fill: #999;
}

.node text {
  font: 10px sans-serif;
}

.node--internal circle {
  fill: #555;
}

.node--internal text {
  text-shadow: 0 1px 0 #fff, 0 -1px 0 #fff, 1px 0 0 #fff, -1px 0 0 #fff;
}

.link {
  fill: none;
  stroke: #555;
  stroke-opacity: 0.4;
  stroke-width: 1.5px;
}
</style>
