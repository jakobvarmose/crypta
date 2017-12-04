/* eslint-disable */
// === DEFAULT / CUSTOM STYLE ===
// WARNING! always comment out ONE of the two require() calls below.
// 1. use next line to activate CUSTOM STYLE (./src/themes)
// require(`./themes/app.${__THEME}.styl`)
// 2. or, use next line to activate DEFAULT QUASAR STYLE
require(`quasar/dist/quasar.${__THEME}.css`)
// ==============================

// Uncomment the following lines if you need IE11/Edge support
// require(`quasar/dist/quasar.ie`)
// require(`quasar/dist/quasar.ie.${__THEME}.css`)

import Vue from 'vue'
import Quasar from 'quasar'
import {
  Dialog,
} from 'quasar';
import escapeHtml from 'escape-html';
import moment from 'moment'
import router from './router'
import i18next from 'i18next';
import translations from './translations';

Vue.config.productionTip = false
Vue.use(Quasar) // Install Quasar Framework

i18next.init({
  lng: 'en',
  fallbackLng: 'en',
  resources: translations,
});

Vue.mixin({
  methods: {
    async $api(name, args) {
      args.myAddress = localStorage.getItem('myAddress');
      const res = await fetch('/api/'+name, {
        method: 'POST',
        body: JSON.stringify(args),
      });
      if (!res.ok) {
        throw new Error(await res.text());
      }
      return res.json();
    },
    $t() {
      return i18next.t.apply(i18next, arguments)
    },
    $alert(err) {
      let title, message;
      if (err instanceof Error) {
        title = err.constructor.name;
        message = err.message;
      } else {
        title = 'Alert';
        message = String(err);
      }
      Dialog.create({
        title: escapeHtml(title),
        message: escapeHtml(message),
      });
    },
  },
});


window.onerror = (msg, url, line, column, err) => {
  console.log(msg, url, line, column, err);
  Dialog.create({
    title: escapeHtml(err.constructor.name),
    message: escapeHtml(err.message),
  });
};

window.onunhandledrejection = (err) => {
  console.log(err.promise, err.reason);
  Dialog.create({
    title: escapeHtml(err.reason.constructor.name),
    message: escapeHtml(err.reason.message),
  });
};


moment.locale('en-gb')

if (__THEME === 'mat') {
  require('quasar-extras/roboto-font')
}
import 'quasar-extras/material-icons'
// import 'quasar-extras/ionicons'
// import 'quasar-extras/fontawesome'
// import 'quasar-extras/animate'

Quasar.start(() => {
  /* eslint-disable no-new */
  new Vue({
    el: '#q-app',
    router,
    render: h => h(require('./App'))
  })
})
