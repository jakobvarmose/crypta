<template>
  <q-layout view="lHh Lpr lFf">
    <q-layout-header>
      <q-toolbar
        color="primary"
        :glossy="$q.theme === 'mat'"
        :inverted="$q.theme === 'ios'"
      >
        <!--<q-btn
          flat
          dense
          round
          @click="leftDrawerOpen = !leftDrawerOpen"
          aria-label="Menu"
        >
          <q-icon name="menu" />
        </q-btn>-->
        <q-toolbar-title>
          {{$t('Crypta')}}
          <!--<div slot="subtitle">Running on Quasar v{{ $q.version }}</div>-->
        </q-toolbar-title>
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
        <!--<router-link to="/user/notifications" style="color:inherit;">
          <q-btn flat><q-icon name="notifications" /></q-btn>
        </router-link>-->
        <router-link v-if="myAddress" to="/user/settings" style="color:inherit;">
          <q-btn flat><q-icon name="settings" /> Settings</q-btn>
        </router-link>
        <q-btn v-if="myAddress" flat @click="logOut">{{$t('Log out')}}</q-btn>
      </q-toolbar>
    </q-layout-header>

    <q-layout-drawer
      v-model="leftDrawerOpen"
      :content-class="$q.theme === 'mat' ? 'bg-grey-2' : null"
    >
      <c-navigation />
    </q-layout-drawer>

    <q-page-container>
      <keep-alive>
        <router-view :key="$route.path"
                     :myAddress="myAddress"
                     @myAddressChange="myAddressChange" />
      </keep-alive>
    </q-page-container>
  </q-layout>
</template>

<script>
import CNavigation from 'components/CNavigation.vue';

export default {
  name: 'MyLayout',
  components: {
    CNavigation,
  },
  data() {
    return {
      leftDrawerOpen: true /* this.$q.platform.is.desktop */,
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
};
</script>

<style>
</style>
