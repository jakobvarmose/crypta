/* XXeslint-disable */
import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter);

function load (component) {
  // '@' is aliased to src/components
  return () => import(`@/${component}.vue`)
}

export default new VueRouter({
  /*
   * NOTE! VueRouter "history" mode DOESN'T works for Cordova builds,
   * it is only to be used only for websites.
   *
   * If you decide to go with "history" mode, please also open /config/index.js
   * and set "build.publicPath" to something other than an empty string.
   * Example: '/' instead of current ''
   *
   * If switching back to default "hash" mode, don't forget to set the
   * build publicPath back to '' so Cordova builds work again.
   */
  mode: 'history',

  routes: [
    { path: '/', component: require('./components/Index'),
     children: [
       {path: '', component: require('./components/Home')},

       {path: 'user/login', component: require('./components/user/UserLogin')},
       {path: 'user/create', component: require('./components/user/UserCreate')},
       {path: 'user/settings', component: require('./components/user/UserSettings')},
       {path: 'user/notifications', component: require('./components/user/UserNotifications')},
       {path: 'user/:address', component: require('./components/user/UserUser')},
       {path: 'post/:creatorAddress-:wallAddress-:postHash', component: require('./components/user/UserPost')},

       {path: 'about', component: require('./components/About')},

       {path: '*', component: require('./components/NotFound')},
     ],
    },

    // Always leave this last one
    { path: '*', component: require('./components/NotFound') } // Not found
  ]
})
