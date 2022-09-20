<template>
    <div id="new">
        <Navigation title="Neuer Eintrag"></Navigation>

        <form @submit.prevent="onSubmit" class="bg-light border-1 m-2 p-2 rounded">
            <div class="mb-3">
                <div class="input-group">
                    <input v-model="query" type="text" class="form-control" placeholder="Freitext Suche" />
                    <button class="btn btn-secondary">Scan</button>
                </div>
            </div>
            <div class="">
                <button type="submit" class="btn btn-primary">Suche</button>
            </div>
        </form>

        <div class="card m-2" v-show="result.discogs">
            <div class="card-header">
                Discogs
            </div>
            <ul class="list-group list-group-flush ">
                <li class="list-group-item d-flex">
                    <div class="form-check">
                        <input type="radio" class="form-check-input" id="d0" v-model="selected.discogs" :value="null">
                        <label for="d0" class="form-check-label text-muted">(Keine Auswahl)</label>
                    </div>
                </li>
                <li class="list-group-item d-flex" v-for="disc in result.discogs">
                    <div class="form-check">
                        <input type="radio" :id="`discogs${disc.id}`" class="form-check-input me-2" v-model="selected.discogs" :value="disc.id" />
                        <label class="form-check-label" :for="`discogs${disc.id}`">{{ disc.title }}</label>
                        <div class="form-text">
                            type: <code>{{disc.type}}</code><br />
                            masterId: <code>{{ disc.master_id}}</code><br />
                            id: <code>{{disc.id}}</code><br />
                            {{disc.format.join(", ")}}
                        </div>
                    </div>
                    <img :src="disc.thumb" alt="" style="width: 75px;height: 75px" class="ms-auto rounded">
                </li>
            </ul>
        </div>

        <div class="card m-2" v-show="result.musicbrainz">
            <div class="card-header">
                Musicbrainz
            </div>
            <ul class="list-group list-group-flush ">
                <li class="list-group-item d-flex">
                    <div class="form-check">
                        <input type="radio" class="form-check-input" id="mb0" v-model="selected.musicbrainz" :value="null">
                        <label for="mb0" class="form-check-label text-muted">(Keine Auswahl)</label>
                    </div>
                </li>
                <li class="list-group-item d-flex" v-for="mb in result.musicbrainz">
                    <div class="form-check">
                        <input type="radio" :id="`mb${mb.ID}`" class="form-check-input" v-model="selected.musicbrainz" :value="mb.ID" />
                        <label class="form-check-label" :for="`mb${mb.ID}`">{{ mb.ArtistCredit.NameCredits[0].Artist.Name}} - {{ mb.Title }}</label>
                        <div class="form-text">
                            {{ mb.Mediums[0].Format }}
                        </div>
                    </div>
                </li>
            </ul>
        </div>


        <div class="card m-2" v-show="result.spotify">
            <div class="card-header">
                Spotify
            </div>
            <ul class="list-group list-group-flush ">
                <li class="list-group-item d-flex">
                    <div class="form-check">
                        <input type="radio" class="form-check-input" id="s0" v-model="selected.spotify" :value="null">
                        <label for="s0" class="form-check-label text-muted">(Keine Auswahl)</label>
                    </div>
                </li>
                <li class="list-group-item d-flex" v-for="s in result.spotify">
                    <div class="form-check">
                        <input type="radio" class="form-check-input" :id="`spotify${s.uri}`" v-model="selected.spotify" :value="s.uri">
                        <label :for="`spotify${s.uri}`" class="form-check-label">{{ s.artists[0].name }} - {{ s.name }}</label>
                    </div>
                    <img :src="s.images[2].url" alt="" style="" class="ms-auto rounded">
                </li>
            </ul>
        </div>

    <div class="m-2 p-2 text-center">
        <button type="button" class="btn btn-lg btn-primary">Speichern</button>
    </div>

    </div>
</template>

<script setup>
import { escapeHtml } from '@vue/shared';
import { reactive, ref } from 'vue';
import Navigation from './Navigation.vue';

const query = ref()
const result = ref([])

const selected = reactive({
    discogs: null,
    musicbrainz: null,
    spotify: null
})

const onSubmit = async() => {
    if (query.value != undefined) {
        result.value = await(await fetch("/api/search/" + escapeHtml(query.value))).json()

        selectMostMatching()
    }
}

const selectMostMatching = () => {
    if (result.value.discogs.length == 1) {
        selected.discogs = result.value.discogs[0].id
    }

    if (result.value.musicbrainz.length == 1) {
        selected.musicbrainz = result.value.musicbrainz[0].ID
    }

    if(result.value.spotify.length == 1) {
        selected.spotify = result.value.spotify[0].uri
    }
}
</script>
