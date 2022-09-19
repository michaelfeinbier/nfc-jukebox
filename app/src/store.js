import { reactive } from 'vue'

export const store = reactive({
    albumList: [],
    async loadAlbumList() {
        if (this.albumList.length == 0) {
            await this.reloadAlbumList()
        }
    },
    async reloadAlbumList() {
        this.albumList = await(await fetch("/api/album")).json()
    }
})
