<script setup>
import { reactive, onBeforeMount, ref, watch } from 'vue';
import Vibrant from 'node-vibrant/dist/vibrant.worker'
import {onBeforeRouteUpdate, useRoute} from "vue-router"

const route = useRoute()
const record = reactive({
    artist: 'unknown',
    name: 'unknown',
    year: 1984,
    artwork: 'https://via.placeholder.com/640',
    color: '#DDD',
    palette: {
        main: '#ddd',
        title: '#000',
        text: '#333'
    }
})

watch(
    () => route.params.id,
    async newId => {
        return await fetchRecord(newId)
    }
)

onBeforeRouteUpdate(() => console.log)

onBeforeMount(async () => {
    await fetchRecord(route.params.id)
})

async function fetchRecord(id) {
    let res = await fetch("/api/album/" +  id)
    let data = await res.json()

    record.artist = data.Artist
    record.name = data.Name
    record.artwork = data.Metadata.Image

    let v = Vibrant.from(record.artwork)
    let p = await v.getPalette()

    record.color = p.Vibrant.hex
    record.palette.main = p.LightMuted.hex
    record.palette.title = p.Vibrant.hex
    record.palette.text = p.Muted.hex
}



</script>

<template>
    <div class="view vh-100 d-flex p-2">
    <div class="card border-0 m-auto shadow-lg">
        <img :src="record.artwork" alt="" class="card-img-top">
        <div class="card-body text-center">
            <h5 class="card-title">
                {{ record.name }}
            </h5>
            <p class="card-text">{{ record.artist }}</p>
        </div>
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
</style>
