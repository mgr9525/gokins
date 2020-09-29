import babelpolyfill from 'babel-polyfill'
import Vue from 'vue'
import App from './App'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-default/index.css'
//import './assets/theme/theme-green/index.css'
import VueRouter from 'vue-router'
import store from './vuex/store'
import Vuex from 'vuex'
//import NProgress from 'nprogress'
//import 'nprogress/nprogress.css'
import routes from './routes'
// import Mock from './mock'
// Mock.bootstrap();
import 'font-awesome/css/font-awesome.min.css'
import api from '@/api/apis'

Vue.use(ElementUI)
Vue.use(VueRouter)
Vue.use(Vuex)

//NProgress.configure({ showSpinner: false });

const router = new VueRouter({
  routes
})

import { getToken,removeToken } from '@/util/storage';
router.beforeEach((to, from, next) => {
  //NProgress.start();
  if (to.path == '/login') {
    removeToken();
    next();
    return;
  }
  let token = getToken();
  if(!token||token==''){
    next({ path: '/login' })
    return;
  }
  next();
})

//router.afterEach(transition => {
//NProgress.done();
//});

Vue.prototype.$post=api.post;
Vue.prototype.$resUrl=api.resUrl;
new Vue({
  //el: '#app',
  //template: '<App/>',
  router,
  store,
  //components: { App }
  render: h => h(App)
}).$mount('#app')

