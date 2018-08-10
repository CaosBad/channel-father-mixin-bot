import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Bot from './views/Bot.vue'
import About from './views/About.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/bot',
      name: 'bot',
      component: Bot
    },
    {
      path: '/about',
      name: 'about',
      component: About
    }
  ]
})