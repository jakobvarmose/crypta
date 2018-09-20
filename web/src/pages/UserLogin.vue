<template>
  <q-page padding>
    <h4><c-title>{{$t('Log in')}}</c-title></h4>
    <q-select v-model="selectedUser" :options="users" />
    <p>{{selectedUser ? selectedUser.address : '&nbsp;'}}</p>
    <q-btn color="primary" :disable="!selectedUser" @click="login">{{$t('Log in')}}</q-btn>
  </q-page>
</template>
<script>
import CTitle from 'components/CTitle.vue';

export default {
  components: {
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
