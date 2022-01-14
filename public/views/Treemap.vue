<template>
  <div class="container-fluid d-flex flex-column h-100 join-view">
    <svg id="treemapGraph" :height="900" :width="900">
      <defs>
        <linearGradient id="vertical-fade" gradientTransform="rotate(90)">
          <stop class="fade-edge" offset="0" />
          <stop class="fade-middle" offset="50%" />
          <stop class="fade-edge" offset="100%" />
        </linearGradient>
      </defs>
      <g id="narrative-graphs-layer" />
      <g id="narrative-summaries-layer" />
    </svg>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { Treemap } from "../store/treemap/index";
import { graphTreemap } from "../util/treemap";
import { actions, getters } from "../store/treemap/module";

export default Vue.extend({
  name: "treemap",

  computed: {
    treemap(): Treemap {
      return getters.getTreemap(this.$store);
    },
  },

  async beforeMount() {
    await actions.fetchTreemap(this.$store, {
      url: "https://onedemo-telco.azurewebsites.net/",
    });
    const treemapData = getters.getTreemap(this.$store);
    graphTreemap("#treemapGraph", treemapData);
  },
});
</script>

<style>
.join-view .nav-link {
  padding: 1rem 0 0.25rem 0;
  border-bottom: 1px solid #e0e0e0;
  color: rgba(0, 0, 0, 0.87);
}
.header-label {
  padding: 1rem 0 0.5rem 0;
  font-weight: bold;
}
</style>
