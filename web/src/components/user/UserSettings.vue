<template>
  <div>
    <h4><c-title>{{$t('My Settings')}}</c-title></h4>
    <p style="color: gray;">
      {{address}}
    </p>
    <q-input v-model="newName" :float-label="$t('Name')"/>
    
    <q-chips-input v-model="newWriters" :float-label="$t('Who can post and comment on my page')" />
    <q-btn color="positive" @click="save">{{$t('Save')}}</q-btn>
    <q-btn @click="reset">{{$t('Reset')}}</q-btn>
  </div>
</template>
<script>
  import {
    QBtn,
    QChipsInput,
    QField,
    QInput,
  } from 'quasar';

  import CTitle from '../c/CTitle.vue';
  import CUser from '../c/CUser.vue';
  import CUserTabs from '../c/CUserTabs.vue';

  export default {
    components: {
      QBtn,
      QChipsInput,
      QField,
      QInput,

      CTitle,
      CUser,
      CUserTabs,
    },
    data() {
      return {
        address: localStorage.getItem('myAddress'),
        name: '',
        newName: '',
        writers: [],
        newWriters: [],
        subscribed: false,
        theyCanPost: false,
      };
    },
    async created() {
      return this.refresh();
    },
    methods: {
      async refresh() {
        const obj = await this.$api('v0/page', {
          address: localStorage.getItem('myAddress'),
        });
        this.name = obj.info.name;
        this.newName = this.name;
        this.writers = Object.keys(obj.writers);
        this.newWriters = this.writers.slice();
        this.subscribed = obj.subscribed;
        this.theyCanPost = obj.theyCanPost;
      },
      async save() {
        await this.$api('v0/page/set', {
          key: 'name',
          val: this.newName,
        });
        await this.$api('v0/page/setwriters', {
          address: this.address,
          writers: this.newWriters,
        });
        this.refresh();
      },
      reset() {
        this.newName = this.name;
      },
    },
  };
</script>
