<template>
  <q-page padding>
    <h4><c-title>{{$t('Notifications')}}</c-title></h4>
    <q-list separator>
      <q-item v-for="notification in notifications" v-bind:key="notification.key">
        <router-link :to="notification.link">{{notification.text}}</router-link>
      </q-item>
    </q-list>
  </q-page>
</template>
<script>
import CTitle from 'components/CTitle.vue';
import CUser from 'components/CUser.vue';
import CUserTabs from 'components/CUserTabs.vue';

export default {
  components: {
    CTitle,
    CUser,
    CUserTabs,
  },
  data() {
    return {
      notifications: [],
    };
  },
  async created() {
    return this.refresh();
  },
  methods: {
    async refresh() {
      const obj = await this.$api('v0/notifications', {
        myAddress: localStorage.getItem('myAddress'),
      });
      this.notifications = obj.notifications;
    },
  },
};
</script>
