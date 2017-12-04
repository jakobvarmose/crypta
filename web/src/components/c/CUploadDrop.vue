<template>
  <div :class="hover ? 'bg-positive text-white' : ''">
    <slot></slot>
  </div>
</template>
<script>
  export default {
    data() {
      return {
        hover: false,
      };
    },
    created() {
      this.$nextTick(() => {
        this.$el.addEventListener('dragover', (e) => {
          e.preventDefault();
          this.hover = true;
        });
        this.$el.addEventListener('dragleave', (e) => {
          e.preventDefault();
          this.hover = false;
        });
        this.$el.addEventListener('drop', (e) => {
          e.preventDefault();
          this.hover = false;
          Array.prototype.forEach.call(e.dataTransfer.files, (file) => {
            this.$emit('file', file);
          });
        });
      });
    },
  };
</script>