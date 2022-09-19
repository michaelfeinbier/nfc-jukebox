<script setup>
import { reactive, onBeforeMount, ref, watch, computed } from 'vue';
import Vibrant from 'node-vibrant/dist/vibrant.worker'
import {onBeforeRouteUpdate, useRoute} from "vue-router"

const route = useRoute()
const record = reactive({
    artist: 'unknown',
    name: 'unknown',
    ReleaseDate: "1984-06-04",
    artwork: 'https://via.placeholder.com/640',
    color: '#DDD',
    tracks: [],
    palette: {
        main: '#ddd',
        title: '#000',
        text: '#333'
    }
})

const year = computed(() => {
    if(record.ReleaseDate.length >= 4) {
        return record.ReleaseDate.substring(0, 4)
    }

    return 1984
})

onBeforeMount(async () => {
    await fetchRecord(route.params.id)
})

async function fetchRecord(id) {
    let res = await fetch("/api/album/" +  id)
    let data = await res.json()

    record.artist = data.Artist
    record.name = data.Name
    record.artwork = data.Metadata.Image
    record.tracks = data.Metadata.Tracks
    record.ReleaseDate = data.Metadata.ReleaseDate

    let v = Vibrant.from(record.artwork)
    let p = await v.getPalette()

    record.color = p.Vibrant.hex
    record.palette.main = p.LightMuted.hex
    record.palette.title = p.Vibrant.hex
    record.palette.text = p.Muted.hex
}



</script>

<template>
    <div class="view min-vh-100 d-flex flex-column p-2">
    <div class="card border-0 shadow-lg m-3">
        <img :src="record.artwork" alt="" class="card-img-top">
        <div class="card-body text-center">
            <h5 class="card-title">
                {{ record.name }}
            </h5>
            <p class="card-text">{{ record.artist }} ({{ year }})</p>
        </div>
    </div>

    <div class="card border-0 m-3">
        <div class="card-header">
            Tracks
        </div>
        <ul class="list-group list-group-flush list-group-numbered">
            <li class="list-group-item d-flex" v-for="track in record.tracks"><span>{{track.Name}}</span></li>
        </ul>
    </div>

    </div>
</template>

<style lang="scss" scoped>
.view {
        background-color: v-bind('record.palette.main');
    }

    .card {
        background-size: cover;
    }

    .card-title {
        color: v-bind('record.palette.title')
    }

    .card-text {
        color: v-bind('record.palette.text')
    }

    .card-header {
        color: v-bind('record.palette.title');
        background-color: v-bind('record.palette.text');

    }
</style>
