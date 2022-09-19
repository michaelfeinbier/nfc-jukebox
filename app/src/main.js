import { createApp } from 'vue'
import App from './App.vue'
import VinylItem from './components/VinylItem.vue'
import List from './components/List.vue'
import { createRouter, createWebHistory } from 'vue-router'

//import './assets/main.css'
import './scss/style.scss'
import * as bootstrap from 'bootstrap'

const routes = [
    { path: '/', component: List},
    { path: '/view/:id', component: VinylItem }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

createApp(App).use(router).mount('#app')
