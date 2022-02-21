<template>
  <v-expansion-panels v-bind="$attrs">
    <template v-for="group in groups">
      <v-expansion-panel :key="group.name" :readonly="group.data.length === 0">
        <v-expansion-panel-header :hide-actions="group.data.length === 0">{{ group.name }}</v-expansion-panel-header>
        <v-expansion-panel-content v-if="group.data.length > 0">
          <dl>
            <template v-for="row in group.data">
              <dt :key="row.name + ':dt'">{{ row.name }}</dt>
              <dd :key="row.name + ':dd'" v-if="row.value">{{ row.value }}</dd>
            </template>
          </dl>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </template>
  </v-expansion-panels>
</template>
<script>
export default {
  name: 'metadata-panel',
  props: {
    metadata: {}
  },
  computed: {
    groups() {
      const m = this.metadata;
      if (!m) {
        return [];
      }

      const result = [];
      const ignoreTop = new Set(['name', 'traitsList', 'nicsList', 'moreMap']);
      const ignoreGroup = new Set(['moreMap']);
      for (const [k, g] of Object.entries(m)) {
        if (ignoreTop.has(k) || !g) {
          continue;
        }
        const group = {
          name: k,
          data: []
        };
        result.push(group);

        for (const [p, v] of Object.entries(g)) {
          if (ignoreGroup.has(p) || !v) {
            continue;
          }
          group.data.push({
            name: p,
            value: v
          });
        }
        if (g.moreMap && g.moreMap.length > 0) {
          group.data.push(...g.moreMap.map(([name, value]) => ({name, value})));
        }
      }

      if (m.nicsList && m.nicsList.length > 0) {
        for (const nic of m.nicsList) {

        }
      }

      if (m.moreMap && m.moreMap.length > 0) {
        const group = {name: 'more', data: m.moreMap.map(([name, value]) => ({name, value}))};
        result.push(group);
      }

      if (m.traitsList && m.traitsList.length > 0) {
        for (const t of m.traitsList) {
          result.push({
            name: t.name,
            data: t.moreMap.map(([name, value]) => ({name, value}))
          });
        }
      }
      return result;
    }
  }
};
</script>
<style scoped>
dl {
  display: grid;
  grid-template-columns: auto 1fr;
  grid-auto-flow: dense;
  grid-gap: 6px 10px;
}

dd {
  grid-column-start: 1;
}

dt {
  text-align: right;
  grid-column-start: 2;
  opacity: 0.6;
}
</style>
