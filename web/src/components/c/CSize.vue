<template>
  <span :title="title">{{text}}</span>
</template>
<script>
  import human from '../../human';

  export default {
    props: [
      'size',
    ],
    computed: {
      text() {
        if (typeof this.size !== 'number') {
          return this.size;
        }
        return human.fromSize(this.size);
      },
      title() {
        if (typeof this.size !== 'number') {
          return this.size;
        }
        return `${this.size} bytes`;
      },
    },
  };
</script>