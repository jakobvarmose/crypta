<template>
  <div>
    <c-post :destinationAddress="wallAddress"
            :creatorAddress="creatorAddress"
            :postHash="postHash" />
  </div>
</template>
<script>
  import CPost from '../c/CPost.vue';
  import CTitle from '../c/CTitle.vue';
  import CUser from '../c/CUser.vue';
  import CUserTabs from '../c/CUserTabs.vue';

  export default {
    components: {
      CPost,
      CTitle,
      CUser,
      CUserTabs,
    },
    computed: {
      wallAddress() {
        return this.$route.params.wallAddress;
      },
      creatorAddress() {
        return this.$route.params.creatorAddress;
      },
      postHash() {
        return this.$route.params.postHash;
      },
    },
  };
</script>