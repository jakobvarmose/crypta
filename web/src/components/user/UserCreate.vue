<template>
  <div>
    <h4><c-title>{{$t('Create user')}}</c-title></h4>
    <q-input v-model="name" :float-label="$t('Name')"></q-input>
    <q-btn color="positive" :disabled="!name" @click="register()">{{$t('Create user')}}</q-btn>
  </div>
</template>
<script>
  import {
    QBtn,
    QInput,
  } from 'quasar';

  import CTitle from '../c/CTitle.vue';

  export default {
    components: {
      QBtn,
      QInput,
      CTitle,
    },
    data() {
      return {
        name: '',
      };
    },
    methods: {
      async register() {
        const obj = await this.$api('v0/user/create', {
          name: this.name,
        });
        this.$emit('myAddressChange', obj.address)
        this.$router.push(`/user/${obj.address}`);
      },
    },
  };
</script>
