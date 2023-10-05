import { defineStore } from "pinia";
import router from "@/router";
import axios from "axios";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    // initialize state from local storage to enable user to stay logged in
    user: JSON.parse(localStorage.getItem("user")!),
    returnUrl: '',
  }),
  actions: {
    async login(username: string, password: string) {
      try {
        const user = await axios.post("/login", {
          username,
          password,
        });

        // update pinia state
        this.user = user;

        // store user details and jwt in local storage to keep user logged in between page refreshes
        localStorage.setItem("user", JSON.stringify(user));

        // redirect to previous url or default to home page
        router.push(this.returnUrl || "/");
      } catch (error) {
        throw error; // Re-throw the error to be caught by the caller
      }
    },
    logout() {
      this.user = null;
      localStorage.removeItem("user");
      router.push("/auth/login");
    },
  },
});
