<template>
  <div id="q-app">
    <router-view />
  </div>
</template>

<script>
import { Dialog } from 'quasar';
import escapeHtml from 'escape-html';
import moment from 'moment';

window.onerror = (msg, url, line, column, err) => {
  // eslint-disable-next-line no-console
  console.log(msg, url, line, column, err);
  Dialog.create({
    title: escapeHtml(err.constructor.name),
    message: escapeHtml(err.message),
  });
};

window.onunhandledrejection = (err) => {
  // eslint-disable-next-line no-console
  console.log(err.promise, err.reason);
  Dialog.create({
    title: escapeHtml(err.reason.constructor.name),
    message: escapeHtml(err.reason.message),
  });
};

moment.locale('en-gb');

export default {
  name: 'App',
};
</script>

<style>
</style>
