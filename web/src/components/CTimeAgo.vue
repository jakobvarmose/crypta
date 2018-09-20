<template>
  <span :title="title">{{text}}</span>
</template>
<script>
import moment from 'moment';

export default {
  props: ['time'],
  data() {
    return {
      text: '',
      title: '',
    };
  },
  methods: {
    update() {
      if (!this.time) {
        this.text = '\u00a0';
        this.title = '';
        return;
      }
      const time = moment(this.time);
      this.text = time.fromNow();
      this.title = time.format('LLL');
    },
  },
  created() {
    this.update();
    this.timeout = setInterval(() => {
      this.update();
    }, 60 * 1000);
  },
  watch: {
    time() {
      this.update();
    },
  },
  beforeDestroy() {
    clearInterval(this.timeout);
  },
};
</script>
