<template>
  <div>
    <h4><c-title>{{filename}}</c-title></h4>
    <p v-if="image">
      <img :src="image.src">
    </p>
    <p v-if="type">{{type}}</p>

    <q-list v-if="files.length">
      <q-list-header>Contents</q-list-header>
      <router-link :to="`/file/${parts.map(encodeURIComponent).join('/')}/${encodeURIComponent(file.name)}`" v-for="file in files" :key="file.name">
        <q-item>
          <q-item-side icon="insert_drive_file"/>
          <q-item-main>{{file.name}}</q-item-main>
        </q-item>
      </router-link>
    </q-list>
    
  </div>
</template>
<script>
  import {
    QCard,
    QItem,
    QItemSide,
    QItemMain,
    QList,
    QListHeader,
  } from 'quasar';
  import CTitle from './c/CTitle.vue';

  export default {
    components: {
      QCard,
      QItem,
      QItemSide,
      QItemMain,
      QList,
      QListHeader,

      CTitle,
    },
    computed: {
      path() {
        return this.$route.params.path;
      },
      parts() {
        return this.$route.params.path.split('/');
      },
      hash() {
        return this.parts[0];
      },
      filename() {
        return this.parts[this.parts.length - 1];
      },
    },
    data() {
      return {
        image: null,
        type: null,
        files: [],
      };
    },
    methods: {
    },
    async created() {
      let block = await this.$blockStore.get(this.hash);
      for (let i = 1; i < this.parts.length; i += 1) {
        const payload = block.toPayload();
        if (payload instanceof Buffer) {
          throw new TypeError();
        }
        if (!(payload.files instanceof Array)) {
          throw new TypeError();
        }
        let block2 = null;
        for (let j = 0; j < payload.files.length; j += 1) {
          if (payload.files[j].name === this.parts[i]) {
            if (!(payload.files[j].link instanceof Object)) {
              throw new TypeError();
            }
            if (!('/' in payload.files[j].link)) {
              throw new TypeError();
            }
            // eslint-disable-next-line no-await-in-loop
            block2 = await this.$blockStore.get(payload.files[j].link['/']);
            break;
          }
        }
        if (block2 === null) {
          throw new TypeError();
        }
        block = block2;
      }
      const payload = block.toPayload();
      if (payload instanceof Buffer) {
        const blob = new Blob([block.toPayload()]);
        this.type = '';
        const url = URL.createObjectURL(blob);
        this.image = {
          src: url,
        };
        this.files = [];
      } else {
        this.type = '';
        this.image = null;
        this.files = payload.files;
      }
    },
    beforeDestroy() {
      for (let i = 0; i < this.images.length; i += 1) {
        const image = this.images[i];
        URL.revokeObjectURL(image.src);
      }
    },
  };
</script>
