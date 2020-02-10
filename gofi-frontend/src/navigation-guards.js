import router from './router'
import store from './store'

import NProgress from 'nprogress' // progress bar
import '@/components/NProgress/nprogress.less' // progress bar custom style
import notification from 'ant-design-vue/es/notification'
import { defaultValue } from '@/utils/util'
import config from '@/config/defaultSettings'
import { message } from 'ant-design-vue'
import i18n from '@/locales'

NProgress.configure({ showSpinner: false })

router.beforeEach((to, from, next) => {
  NProgress.start() // start progress bar
  if (to.matched.some(record => record.meta.requireAuth)) {
    // this route requires auth, check if logged in
    // if not, redirect to login page.
    if (!store.getters.isLogin) {
      message.warn(i18n.t('auth.requireAuth.content'))
      next({ name: 'login' })
    } else {
      fetchConfigurationIfNotExist(next, to)
    }
  } else {
    fetchConfigurationIfNotExist(next, to)
  }
})

router.afterEach(() => {
  NProgress.done() // finish progress bar
})

function setupOrNext (next, to) {
  // if initialized, just call next()
  const initialized = store.getters.configuration.initialized
  console.log('initialized ' + initialized)
  if (to.name === 'setup') {
    if (initialized) {
      next({ replace: true, path: '/' })
    } else {
      next()
    }
  } else {
    if (initialized) {
      next()
    } else {
      next({ replace: true, name: 'setup' })
    }
  }
}

function applyConfigurationFromServer (configuration) {
  store.commit('TOGGLE_THEME',
    defaultValue(configuration.themeStyle, config.navTheme))
  store.commit('TOGGLE_NAV_MODE',
    defaultValue(configuration.navMode, config.navMode))
}

function fetchConfigurationIfNotExist (next, to) {
  if (store.getters.configurationValid) {
    console.log('configuration valid')
    setupOrNext(next, to)
  } else {
    console.log('get configuration')
    // make sure configuration is always exist before navigate to any where.
    store.dispatch('GetConfiguration')
      .then((data) => {
        setupOrNext(next, to)
        applyConfigurationFromServer(data)
      })
      .catch((e) => {
        notification.error({
          message: '错误',
          description: e
        })
      })
  }
}
