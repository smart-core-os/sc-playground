import Vue from 'vue';
import VueRouter from 'vue-router';

import devices from './views/devices/route.js';
import traits from './views/traits/route.js';

Vue.use(VueRouter);
let router = new VueRouter({
  routes: [
    {path: '/', redirect: traits.path},
    devices,
    traits
  ]
});
export default router;
