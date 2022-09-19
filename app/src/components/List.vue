<script setup>
import { onBeforeMount, ref } from 'vue';
import { store } from '../store'

onBeforeMount(async () => {
    store.loadAlbumList()
})

</script>

<template>
    <div class="home p-3">
        <div class="list-group">
            <router-link :to="`/view/${album.Id}`" class="list-group-item d-flex justify-content-start" v-for="album in store.albumList">
                <img class="rounded" :src="album.Metadata.Image" />
                <div class="title">
                    <h5>{{ album.Artist}}</h5>
                    {{album.Name}}
                </div>
            </router-link>
        </div>
    </div>
</template>

<style scoped lang="scss">
.title {
        h5 {
            margin-bottom: 0;
        }
        padding: var(--bs-list-group-item-padding-y) var(--bs-list-group-item-padding-x);
    }

    .list-group-item {
        padding: 0;
    }

    img.rounded {
        width: 100px;
        max-height: auto;
    }
</style>
