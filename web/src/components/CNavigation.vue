<template>
  <q-list no-border link inset-separator>
    <q-list-header>{{$t('Navigation')}}</q-list-header>
    <q-item to="/" exact>
      <q-item-side icon="home"/>
      <q-item-main>{{$t('Home')}}</q-item-main>
    </q-item>
    <q-collapsible indent icon="person" :label="$t('User')">
      <q-item :to="`/user/${address}`" v-if="loggedIn">
        <q-item-main>{{$t('My profile')}}</q-item-main>
      </q-item>
      <q-item to="/user/settings" v-if="loggedIn">
        <q-item-main>{{$t('Settings')}}</q-item-main>
      </q-item>
      <q-item to="/user/login" v-if="!loggedIn">
        <q-item-main>{{$t('Log in')}}</q-item-main>
      </q-item>
      <q-item to="/user/create" v-if="!loggedIn">
        <q-item-main>{{$t('Register')}}</q-item-main>
      </q-item>
      <q-item @click="logOut" v-if="loggedIn">
        <q-item-main>{{$t('Log out')}}</q-item-main>
      </q-item>
    </q-collapsible>
    <q-item to="/about">
      <q-item-side icon="info"/>
      <q-item-main>{{$t('About')}}</q-item-main>
    </q-item>
  </q-list>
</template>
<script>
export default {
  data() {
    return {
      loggedIn: false,
      address: '',
    };
  },
  methods: {
    logOut() {
      this.$emit('myAddressChange', null);
    },
  },
};
</script>
