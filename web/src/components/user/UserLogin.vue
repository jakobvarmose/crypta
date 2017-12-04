<template>
  <div>
    <h4><c-title>{{$t('Log in')}}</c-title></h4>
    <q-select v-model="selectedUser" :options="users" />
    <p>{{selectedUser ? selectedUser.address : '&nbsp;'}}</p>
    <q-btn color="primary" :disable="!selectedUser" @click="login">{{$t('Log in')}}</q-btn>
  </div>
</template>
<script>
  import {
    QBtn,
    QInput,
    QSelect,
  } from 'quasar';
  import CTitle from '../c/CTitle.vue';

  export default {
    components: {
      QBtn,
      QInput,
      QSelect,
      CTitle,
    },
    data() {
      return {
        users: [],
        selectedUser: null,
      };
    },
    methods: {
      login() {
        this.$emit('myAddressChange', this.selectedUser.address);
        this.$router.push(`/user/${this.selectedUser.address}`);
      },
    },
    async created() {
      const obj = await this.$api('v0/user/list', {});
      this.users = obj.map(user => ({
        label: user.name,
        value: user,
      }));
    },
  };
</script>
