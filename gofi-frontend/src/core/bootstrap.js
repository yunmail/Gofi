import Vue from 'vue'
import store from '@/store/'
import {
  DEFAULT_LANGUAGE,
  DEFAULT_NAV_MODE,
  DEFAULT_THEME,
  TOKEN
} from '@/store/mutation-types'
import config from '@/config/defaultSettings'

export default function Initializer () {
  store.commit('TOGGLE_THEME', Vue.ls.get(DEFAULT_THEME, config.navTheme))
  store.commit('TOGGLE_NAV_MODE', Vue.ls.get(DEFAULT_NAV_MODE, config.navMode))
  store.commit('SWITCH_LANGUAGE', Vue.ls.get(DEFAULT_LANGUAGE, config.language))
  const token = Vue.ls.get(TOKEN, '')
  store.commit('SET_TOKEN', token)
  // 如果token存在，请求用户信息
  if (token) {
    store.dispatch('GetUser')
  }
}
