<template>
  <div>
    <h4><c-title>{{$t('Crypta')}}</c-title></h4>
    <div v-if="showInfo">
      <p>A decentralized social network. Built with <a href="https://ipfs.io/">IPFS</a>, <a href="http://quasar-framework.org/">Quasar</a> and <a href="https://vuejs.org/">Vue</a>.</p>
      <a href="https://github.com/jakobvarmose/crypta">View project on GitHub</a>
    </div>
    <c-post v-for="post in posts"
            :key="`${post.hash}`"
            :post="post" />
  </div>
</template>
<script>
  import {
    QAlert,
    QBtn,
    QCard,
    QCardActions,
    QCardMain,
    QCardTitle,
  } from 'quasar';

  import CPost from './c/CPost.vue';
  import CTitle from './c/CTitle.vue';

  export default {
    components: {
      QAlert,
      QBtn,
      QCard,
      QCardActions,
      QCardMain,
      QCardTitle,

      CPost,
      CTitle,
    },
    data() {
      return {
        posts: [],
        showInfo: false,
      }
    },
  
    async created() {
      try {
        const obj = await this.$api('v0/home', {});
        this.posts = obj.posts;
      } catch (e) {
      }
      this.showInfo = this.posts.length === 0;
    },
  };
</script>
