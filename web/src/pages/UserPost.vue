<template>
  <q-page padding>
    <c-post :destinationAddress="wallAddress"
            :creatorAddress="creatorAddress"
            :postHash="postHash" />
  </q-page>
</template>
<script>
import CPost from 'components/CPost.vue';
import CTitle from 'components/CTitle.vue';
import CUser from 'components/CUser.vue';
import CUserTabs from 'components/CUserTabs.vue';

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
