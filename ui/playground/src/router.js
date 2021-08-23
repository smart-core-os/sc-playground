import Vue from 'vue';
import VueRouter from 'vue-router';

import devices from './views/devices/route.js';
import home from './views/home/route.js';

Vue.use(VueRouter);
let router = new VueRouter({
  routes: [
    home,
    devices
  ]
});
export default router;
