/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import { loadFonts } from "./webfontloader";
import vuetify from "./vuetify";
import pinia from "../store";
import router from "../router";
import VueProgressBar from "@aacassandra/vue3-progressbar";
import Vue3Toasity, { type ToastContainerOptions } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';
import { createI18n } from 'vue-i18n'
import enUS from '../locales/en-US.json'
import ro from '../locales/ro.json'
import axios from 'axios'
import VueAxios from 'vue-axios'
import { VueReCaptcha } from 'vue-recaptcha-v3'
import { abilitiesPlugin } from '@casl/vue';
import ability from '../services/ability';
import '@/services/interceptors'
import VueApexCharts from "vue3-apexcharts";

// Types
import type { App } from "vue";


const i18n = createI18n({
  locale: 'en-US',
  messages: {
    "en-US": enUS,
    "ro": ro
  },
})

axios.defaults.baseURL = import.meta.env.VITE_API_URL || import.meta.env.BASE_URL
axios.defaults.withCredentials = true

export function registerPlugins(app: App) {
  loadFonts();
  app.use(vuetify).use(router).use(pinia).use(
    Vue3Toasity,
    {
      autoClose: 3000,
    } as ToastContainerOptions,
  ).use(i18n).use(VueAxios, axios).use(VueReCaptcha, {
    siteKey: '<site key>',
    loaderOptions: {
      useRecaptchaNet: true,
      autoHideBadge: true
    }
  }).use(abilitiesPlugin, ability, {
    useGlobalProperties: true
  }).use(VueApexCharts).use(VueProgressBar, {
    thickness: "4px",
  });
  app.provide('axios', app.config.globalProperties.axios)
}
