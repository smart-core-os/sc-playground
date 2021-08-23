import {defineConfig} from 'vite';
import {createVuePlugin} from 'vite-plugin-vue2';
import VitePluginComponents, {VuetifyResolver} from 'vite-plugin-components';
import rollupPluginVuetify from 'rollup-plugin-vuetify';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    createVuePlugin(),
    VitePluginComponents({
      dirs: '',
      customComponentResolvers: [
        VuetifyResolver()
      ]
    }),
    rollupPluginVuetify()
  ]
})
