<template>
  <div>
    <h4><c-title>{{$t('Notifications')}}</c-title></h4>
    <q-list separator>
      <q-item v-for="notification in notifications">
        <router-link :to="notification.link">{{notification.text}}</router-link>
      </q-item>
    </q-list>
  </div>
</template>
<script>
  import {
    QItem,
    QList,
  } from 'quasar';

  import CTitle from '../c/CTitle.vue';
  import CUser from '../c/CUser.vue';
  import CUserTabs from '../c/CUserTabs.vue';

  export default {
    components: {
      QItem,
      QList,

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
