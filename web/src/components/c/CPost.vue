<template>
  <q-card>
    <q-card-title>
      <c-user :address="post.genesis.creator.address" :name="post.genesis.creator.name" />
      <span v-if="post.genesis.creator.address !== post.page.address">
        &gt; <c-user :address="post.page.address" :name="post.page.name"/>
      </span>
      <span slot="subtitle">
        <c-time-ago class="text-muted" :time="post.genesis.time"/>
      </span>
      <span slot="right">
        <!--<q-icon name="more_vert">
          <q-popover anchor="bottom right" self="top right" ref="popover">
            <q-list link>
              <q-side-link item :to="`/post/${post.genesis.creator.address}-${post.page.address}-${postHash}`">
                <q-item-main>
                  <q-icon name="link" />
                  {{$t('Permalink')}}
                </q-item-main>
              </q-side-link>
              <q-item @click="deletePost(), $refs.popover.close()" v-if="canDelete">
                <q-item-main>
                  <q-icon name="delete" />
                  {{$t('Delete')}}
                </q-item-main>
              </q-item>
            </q-list>
          </q-popover>
        </q-icon>-->
      </span>
    </q-card-title>
    <q-card-main>
      <span style="white-space: pre-wrap; word-wrap: break-word; overflow-wrap: break-word;">{{post.genesis.text}}</span>
    </q-card-main>
    <q-card-separator />
    <q-card-main>
      <p v-for="comment in post.comments">
        <c-user :address="comment.creator.address" :name="comment.creator.name"/>:
        <span style="white-space: pre-wrap; word-wrap: break-word; overflow-wrap: break-word;">{{comment.text}}</span>
      </p>
      <q-input v-model="commentText"
               :placeholder="$t('Your comment...')"
               :disable="commentDisabled"
               @keydown.enter="createComment" />
    </q-card-main>
  </q-card>
</template>
<script>
  import {
    QCard,
    QCardActions,
    QCardMain,
    QCardSeparator,
    QCardTitle,
    QIcon,
    QInput,
    QItem,
    QItemMain,
    QList,
    QPopover,
    QSideLink,
  } from 'quasar';

  import CUser from './CUser.vue';
  import CTimeAgo from './CTimeAgo.vue';

  export default {
    props: [
      'creatorAddress',
      'destinationAddress',
      'postHash',
      'post',
    ],
    components: {
      QCard,
      QCardActions,
      QCardMain,
      QCardSeparator,
      QCardTitle,
      QIcon,
      QInput,
      QItem,
      QItemMain,
      QList,
      QPopover,
      QSideLink,

      CUser,
      CTimeAgo,
    },
    data() {
      return {
        text: '',
        time: 0,
        canDelete: false,
        commentText: '',
        commentDisabled: false,
      };
    },
    async created() {
      this.refresh();
    },
    methods: {
      async refresh() {
        // FIXME only show post if creator has permission to
        // post on the wall
        // TODO implement fetching of post data
      },
      async deletePost() {
        
      },
      async createComment() {
        if (!this.commentDisabled && this.commentText) {
          this.commentDisabled = true;
          try {
            const obj = await this.$api('v0/page/comment', {
              address: this.post.page.address,
              postHash: this.post.hash,
              text: this.commentText,
            });
            if (obj) {
              this.commentText = '';
              this.post.comments.push(obj.result);
            }
          } catch (e) {
            this.$alert(e);
          }
          this.commentDisabled = false;
        }
      },
    },
  };
</script>
