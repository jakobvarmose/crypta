<template>
  <q-page padding>
    <c-user-tabs :address="address"
                 :name="name"
                 :subscribed="subscribed"
                 :canPost="theyCanPost" />
    <div class="form-group" v-if="hasPermissions">
      <div><b>{{$t('Permissions:')}}</b></div>
      <div class="form-check">
        <label class="form-check-label">
          <input class="form-check-input" type="checkbox" v-model="writePermissions">
          {{$t('Write to my page')}}
        </label>
      </div>
    </div>
    <div v-if="canPost">
      <div style="display: flex; flex-direction: row; width: 100%;">
        <div
          v-for="(file, i) in postFiles"
          v-bind:key="file.id"
          style="position: relative; margin: 5px;"
        >
          <img :src="'/ipfs/'+file.hash" style="width: 100%;">
          <q-btn
            flat
            style="position: absolute; top: 0px; left: 0px; color: white;"
            @click="removeFile(i)"
          >
            <q-icon name="delete" /> {{$t('Remove')}}
          </q-btn>
        </div>
      </div>
      <q-input
        type="textarea"
        :min-rows="4"
        v-model="postText"
        :float-label="$t('Say something...')"
      />
      <q-btn
        color="primary"
        :disable="!postText&&!postFiles.length"
        @click="doPost"
      >{{$t('Post')}}</q-btn>
      <c-upload-btn @file="doUpload">{{$t('Add pictures/videos')}}</c-upload-btn>
    </div>
    <c-post v-for="post in posts" :key="`${post.hash}`" :post="post" />
  </q-page>
</template>
<script>
import uid from 'uid-safe';
import CPost from 'components/CPost.vue';
import CUser from 'components/CUser.vue';
import CTimeAgo from 'components/CTimeAgo.vue';
import CTitle from 'components/CTitle.vue';
import CUploadBtn from 'components/CUploadBtn.vue';
import CUserTabs from 'components/CUserTabs.vue';

export default {
  components: {
    CPost,
    CUser,
    CTimeAgo,
    CTitle,
    CUploadBtn,
    CUserTabs,
  },
  data() {
    return {
      name: '',
      posts: [],
      postText: '',
      postFiles: [],
      address: this.$route.params.address,
      hasPermissions: false,
      writePermissions: false,
      subscribed: false,
      canPost: false,
      theyCanPost: false,
    };
  },
  methods: {
    async doPost() {
      const obj = await this.$api('v0/page/post', {
        address: this.address,
        text: this.postText,
        attachments: this.postFiles,
      });
      if (obj) {
        this.postText = '';
        this.postFiles = [];
        obj.result.page = {
          name: this.name,
          address: this.address,
        };
        this.posts.splice(0, 0, obj.result);
      }
    },
    async refresh2() {
      const obj = await this.$api('v0/page', {
        address: this.address,
      });
      this.name = obj.info.name;
      this.posts = obj.posts;
      this.canPost = this.address in obj.writers;
      this.subscribed = obj.subscribed;
      this.theyCanPost = obj.theyCanPost;
    },
    async doUpload(file) {
      const res = await fetch('/api/v0/upload', {
        method: 'POST',
        body: file,
      });
      const obj = await res.json();
      this.postFiles.push({
        id: await uid(18),
        t: obj.t,
        hash: obj.hash,
      });
    },
    async removeFile(i) {
      this.postFiles.splice(i, 1);
    },
  },
  async created() {
    this.refresh2();
  },
};
</script>
