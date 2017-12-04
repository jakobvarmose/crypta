<template>
  <router-link style="white-space: pre-wrap; word-wrap: break-word; overflow-wrap: break-word;" :to="'/user/'+address" :title="address">{{name || '&nbsp;'}}</router-link>
</template>
<script>
  export default {
    props: [
      'address',
      'name',
    ],
  };
</script>
