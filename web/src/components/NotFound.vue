<template>
  <div>
    <h4><c-title>{{$t('Not found')}}</c-title></h4>
  </div>
</template>
<script>
  import CTitle from './c/CTitle.vue';

  export default {
    components: {
      CTitle,
    },
  };
</script>
