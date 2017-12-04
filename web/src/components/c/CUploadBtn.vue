<template>
  <q-btn @click="$refs.file.click()">
    <input type="file" ref="file" multiple hidden
           @change="fileChange" />
    <slot />
  </q-btn>
</template>
<script>
  import {
    QBtn,
  } from 'quasar';

  export default {
    components: {
      QBtn,
    },
    methods: {
      fileChange(e) {
        Array.prototype.forEach.call(e.target.files, (file) => {
          this.$emit('file', file);
        });
      },
    },
  };
</script>