import axios from "axios";
import router from "@/router";

axios.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error) {
    if (error.response?.status === 401) {
      router.push("/auth/login");
    }
    return Promise.reject(error);
  }
);
