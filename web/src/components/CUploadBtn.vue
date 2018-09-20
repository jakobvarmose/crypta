<template>
  <q-btn @click="$refs.file.click()">
    <input type="file" ref="file" multiple hidden
           @change="fileChange" />
    <slot />
  </q-btn>
</template>
<script>
export default {
  methods: {
    fileChange(e) {
      Array.prototype.forEach.call(e.target.files, (file) => {
        this.$emit('file', file);
      });
    },
  },
};
</script>
