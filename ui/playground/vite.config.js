import {defineConfig} from 'vite';
import {createVuePlugin} from 'vite-plugin-vue2';
import VitePluginComponents, {VuetifyResolver} from 'vite-plugin-components';
import rollupPluginVuetify from 'rollup-plugin-vuetify';

const traits = [
  'parent',
  'electric',
  'energy_storage',
  'metadata',
  'power_supply'
]
const includes = traits.reduce((arr, trait) => {
  arr.push(`@smart-core-os/sc-api-grpc-web/traits/${trait}_pb.js`);
  arr.push(`@smart-core-os/sc-api-grpc-web/traits/${trait}_grpc_web_pb.js`);
  return arr;
}, []);

// https://vitejs.dev/config/
export default defineConfig({
  optimizeDeps: {
    include: [
      ...includes
    ]
  },
  plugins: [
    createVuePlugin(),
    VitePluginComponents({
      dirs: '',
      customComponentResolvers: [
        VuetifyResolver()
      ]
    }),
    rollupPluginVuetify()
  ],
  server: {
    fs: {
      allow: ['..']
    }
  },
  build: {
    commonjsOptions: {include: []}
  }
})
