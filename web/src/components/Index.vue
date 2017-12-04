<template>
  <q-layout ref="layout"
            view="lHh Lpr fff"
            :left-class="{'bg-grey-2': true}">
    <q-toolbar slot="header" class="glossy">
      <!--<q-btn flat
             @click="$refs.layout.toggleLeft()">
        <q-icon name="menu" />
      </q-btn>-->
      <q-toolbar-title>{{$t('Crypta')}}</q-toolbar-title>
      <router-link to="/" style="color:inherit;">
        <q-btn flat><q-icon name="home" /> {{$t('Home')}}</q-btn>
      </router-link>
      <router-link v-if="!myAddress" to="/user/login" style="color:inherit;">
        <q-btn flat>{{$t('Log in')}}</q-btn>
      </router-link>
      <router-link v-if="!myAddress" to="/user/create" style="color:inherit;">
        <q-btn flat>{{$t('Create user')}}</q-btn>
      </router-link>
      <router-link v-if="myAddress" :to="'/user/'+myAddress" style="color:inherit;">
        <q-btn flat><q-icon name="person" /> My Page</q-btn>
      </router-link>
      <router-link v-if="myAddress" to="/user/settings" style="color:inherit;">
        <q-btn flat><q-icon name="settings" /> Settings</q-btn>
      </router-link>
      <q-btn v-if="myAddress" flat @click="logOut">{{$t('Log out')}}</q-btn>
    </q-toolbar>
    <!--<c-navigation slot="left" />-->
    <div class="layout-padding">
      <keep-alive>
        <router-view :key="$route.path"
                     :myAddress="myAddress"
                     @myAddressChange="myAddressChange" />
      </keep-alive>
    </div>
  </q-layout>
</template>
<script>
  import {
    QBtn,
    QIcon,
    QLayout,
    QToolbar,
    QToolbarTitle,
  } from 'quasar'

  import CNavigation from './c/CNavigation.vue';

  export default {
    components: {
      QBtn,
      QIcon,
      QLayout,
      QToolbar,
      QToolbarTitle,

      CNavigation,
    },
    
    data() {
      return {
        myAddress: localStorage.getItem('myAddress'),
      };
    },
    methods: {
      logOut() {
        this.myAddressChange(null);
      },
      myAddressChange(arg) {
        if (arg) {
          localStorage.setItem('myAddress', arg);
          this.myAddress = arg;
        } else {
          localStorage.removeItem('myAddress');
          this.myAddress = null;
        }
      },
    },
  }
</script>
