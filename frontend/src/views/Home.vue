<template>
  <div style="text-align: center">
    Hello {{ Name }} <br /><br />
    CASL Permission test:
    <div v-if="$can('create', 'Post')">
      <a @click="createPost">Create Post</a>
    </div>
    <hr />
    <button @click="recaptcha">Execute recaptcha</button>
    <hr />
    <div class="locale-changer">
      <select v-model="$i18n.locale">
        <option
          v-for="locale in $i18n.availableLocales"
          :key="`locale-${locale}`"
          :value="locale"
        >
          {{ locale }}
        </option>
      </select>
    </div>
    <p>{{ $t("message.hello") }}</p>
    <hr />
    <div>
      <button @click="notify">Notify !</button>
    </div>
    <a href="#" @click="logout">Logout</a>
  </div>

  <div>
    <apexchart
      v-if="chartOptions && chartSeries"
      :type="chartOptions.type"
      :options="chartOptions"
      :series="chartSeries"
    />
  </div>

  <!-- <BarChart /> -->

  <canvas ref="chartCanvas"></canvas>

  <Terminal />
</template>

<script lang="ts" setup>
import { toast } from "vue3-toastify";
// import Msg from './Msg.vue';
import "vue3-toastify/dist/index.css";
import { inject, getCurrentInstance } from "vue";
import { VueReCaptcha, useReCaptcha } from "vue-recaptcha-v3";
import Terminal from "@/components/Terminal.vue";
import { useAuthStore } from "@/store/auth";
import router from "@/router";

const instance = getCurrentInstance();
const progress = instance.appContext.config.globalProperties.$Progress;

progress.start();

const axios: any = inject("axios"); // inject axios

const Name = ref("John Doe");

// axios.get("example.com").then((response: { data: any }) => {
//   console.log(response.data);
// });

const { executeRecaptcha, recaptchaLoaded } = useReCaptcha();

const recaptcha = async () => {
  // (optional) Wait until recaptcha has been loaded.
  await recaptchaLoaded();

  // Execute reCAPTCHA with action "login".
  const token = await executeRecaptcha("login");

  alert(token);

  // Do stuff with the received token.
};

const notify = () => {
  const id = toast.loading("Please wait...", {
    transition: toast.TRANSITIONS.SLIDE,
    position: toast.POSITION.BOTTOM_RIGHT,
    theme: "dark",
  });

  setTimeout(() => {
    toast.update(id, {
      render: "Successfully logged in! Redirecting...",
      autoClose: true,
      closeOnClick: true,
      closeButton: true,
      type: "success",
      isLoading: false,
    });
  }, 2000);
};

axios.get("/test");

const store = useAuthStore();

async function logout() {
  await store.logout();
  router.push("/auth/login");
}

import { Chart, ChartConfiguration } from "chart.js/auto";
import prometheus from "chartjs-plugin-datasource-prometheus";
import "chartjs-adapter-date-fns";

Chart.registry.plugins.register(prometheus);

const chartInstance = ref<Chart | null>(null);
const chartCanvas = ref<HTMLCanvasElement | null>(null);

onMounted(() => {
  if (chartCanvas.value) {
    const ctx = chartCanvas.value.getContext("2d");
    const chartConfig: ChartConfiguration = {
      type: "line",
      plugins: [prometheus],
      options: {
        responsive: true,
        plugins: {
          title: {
            display: true,
            text: "Chart.js Line Chart - Cubic interpolation mode",
          },
          "datasource-prometheus": {
            prometheus: {
              endpoint: "https://prometheus.demo.do.prometheus.io",
              baseURL: "/api/v1", // default value
            },
            query: "sum by (job) (go_gc_duration_seconds)",
            timeRange: {
              type: "relative",

              // from 12 hours ago to now
              start: -12 * 60 * 60 * 1000,
              end: 0,
            },
          },
        },
        interaction: {
          intersect: false,
          axis: "x",
        },
      },
    };
    chartInstance.value = new Chart(ctx, chartConfig);
  }
});

import VueApexCharts from "vue3-apexcharts";

const chartOptions = ref({
  type: "line",
  chart: {
    id: "vue-apexcharts-example",
  },
  xaxis: {
    categories: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep"],
  },
});

const chartSeries = ref([
  {
    name: "Series 1",
    data: [30, 40, 25, 50, 49, 21, 70, 51, 60],
  },
]);

import BarChart from "@/components/BarChart.vue";
</script>
