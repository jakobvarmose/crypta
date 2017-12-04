<template>
  <div>
    <q-list no-border link inset-separator>
      <q-list-header>{{$t('Navigation')}}</q-list-header>
      <q-side-link item to="/" exact>
        <q-item-side icon="home"/>
        <q-item-main>{{$t('Home')}}</q-item-main>
      </q-side-link>
      <q-collapsible indent icon="person" :label="$t('User')">
        <q-side-link item :to="`/user/${address}`" v-if="loggedIn">
          <q-item-main>{{$t('My profile')}}</q-item-main>
        </q-side-link>
        <q-side-link item to="/user/settings" v-if="loggedIn">
          <q-item-main>{{$t('Settings')}}</q-item-main>
        </q-side-link>
        <q-side-link item to="/user/login" v-if="!loggedIn">
          <q-item-main>{{$t('Log in')}}</q-item-main>
        </q-side-link>
        <q-side-link item to="/user/register" v-if="!loggedIn">
          <q-item-main>{{$t('Register')}}</q-item-main>
        </q-side-link>
        <q-item link @click="logOut" v-if="loggedIn">
          <q-item-main>{{$t('Log out')}}</q-item-main>
        </q-item>
      </q-collapsible>
      <q-side-link item to="/about">
        <q-item-side icon="info"/>
        <q-item-main>{{$t('About')}}</q-item-main>
      </q-side-link>
    </q-list>
  </div>
</template>
<script>
  import {
    QCollapsible,
    QItem,
    QItemMain,
    QItemSide,
    QList,
    QListHeader,
    QSideLink,
  } from 'quasar';

  export default {
    components: {
      QCollapsible,
      QItem,
      QItemMain,
      QItemSide,
      QList,
      QListHeader,
      QSideLink,
    },
    data() {
      return {
        loggedIn: false,
        address: '',
      };
    },
    methods: {
      logOut() {
        this.$emit('myAddressChange', null)
      },
    },
  };
</script>