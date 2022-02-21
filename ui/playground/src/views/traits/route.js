/** @type {import('vue-router').RouteConfig} */
const route = {
  path: '/traits',
  name: 'traits',
  component: () => import('./TraitCards.vue')
}
export default route;
