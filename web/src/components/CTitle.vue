<template>
  <span><slot></slot></span>
</template>
<script>
const originalTitle = document.title;
export default {
  activated() {
    document.title = this.$el.innerText;
    this.observer = new MutationObserver((/* mutations */) => {
      document.title = this.$el.innerText;
    });
    this.observer.observe(this.$el, {
      characterData: true,
      childList: true,
      subtree: true,
    });
  },
  deactivated() {
    document.title = originalTitle;
    this.observer.disconnect();
  },
};
</script>
