<template>
  <div class="d-flex align-center justify-center" style="height: 100vh">
    <v-sheet width="400" class="mx-auto">
      <form @submit.prevent="submit">
        <v-text-field
          v-model="email.value.value"
          :error-messages="email.errorMessage.value"
          label="E-mail"
          placeholder="johndoe@email.com"
          autofocus
        ></v-text-field>

        <v-text-field
          v-model="password.value.value"
          :counter="10"
          :error-messages="password.errorMessage.value"
          label="Password"
          type="password"
          placeholder="··········"
        ></v-text-field>

        <v-checkbox
          v-model="remember.value.value"
          :error-messages="remember.errorMessage.value"
          value="1"
          label="Remember me"
          type="checkbox"
        ></v-checkbox>

        <v-btn class="me-4" type="submit" :loading="loading"> login </v-btn>
      </form>
    </v-sheet>
  </div>
</template>

<script lang="ts" setup>
import { useField, useForm } from "vee-validate";
import { useAuthStore } from "@/store/auth";
import { toast } from "vue3-toastify";
import "vue3-toastify/dist/index.css";

const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    email(value: string) {
      if (/^[a-z.-]+@[a-z.-]+\.[a-z]+$/i.test(value)) return true;

      return "Must be a valid e-mail.";
    },
    password(value: string) {
      if (value?.length >= 6) return true;

      return "Password needs to be at least 6 characters.";
    },
  },
});
const email = useField("email");
const password = useField("password");
const remember = useField("remember");
const loading = ref(false);

const submit = handleSubmit((values) => {
  const store = useAuthStore();

  const loadingId = toast.loading("Please wait...", {
    transition: toast.TRANSITIONS.SLIDE,
    position: toast.POSITION.TOP_RIGHT,
    theme: "dark",
  });

  loading.value = true;

  store
    .login("admin", "admin")
    .then(() => {
      toast.update(loadingId, {
        render: "Successfully logged in! Redirecting...",
        autoClose: true,
        closeOnClick: true,
        closeButton: true,
        type: "success",
        isLoading: false,
      });
    })
    .catch((error) => {
      toast.update(loadingId, {
        render: error.message,
        autoClose: true,
        closeOnClick: true,
        closeButton: true,
        type: "error",
        isLoading: false,
      });
    })
    .finally(() => {
      setTimeout(() => {
        loading.value = false;
      }, 3000);
    });
});
</script>
