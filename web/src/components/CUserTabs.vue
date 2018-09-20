<template>
  <div>
    <h4>
      <c-title :key="name"><c-user :address="address" :name="name" /></c-title>
    </h4>
    <p style="color: gray;">
      {{address}}
    </p>
    <p>
      <q-checkbox v-model="subscribed"
                  label="subscribed"
                  @change="saveSubscribed" />
      <q-checkbox v-model="canPost"
                  label="allow posts"
                  @change="saveCanPost" />
    </p>
    <!--
    <q-tabs>
      <q-route-tab
        icon="home"
        :to="`/user/${address}`" exact
        slot="title">
        {{$t('Page')}}
      </q-route-tab>
      <!-/-
      <q-route-tab
        icon="folder"
        :to="`/user/${address}/documents`"
        slot="title">
        {{$t('Documents')}}
      </q-route-tab>
      <q-route-tab
        icon="chat"
        :to="`/user/${address}/chat`"
        slot="title">
        {{$t('Chat')}}
      </q-route-tab>
      -/->
      <q-route-tab
        icon="settings"
        :to="`/user/${address}/settings`"
        slot="title">
        {{$t('Settings')}}
      </q-route-tab>
    </q-tabs>
    -->
  </div>
</template>
<script>
import CUser from 'components/CUser.vue';
import CTitle from 'components/CTitle.vue';

export default {
  props: [
    'address',
    'name',
    'subscribed',
    'canPost',
  ],
  components: {
    CUser,
    CTitle,
  },
  methods: {
    async saveSubscribed() {
      await this.$api('v0/subscribe', {
        address: this.address,
        value: this.subscribed,
      });
    },
    async saveCanPost() {
      await this.$api('v0/canPost', {
        address: this.address,
        value: this.canPost,
      });
    },
  },
};
</script>
