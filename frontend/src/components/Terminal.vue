<template>
  <div>
    <div ref="terminalContainer" style="height: 300px"></div>
    <button @click="clearTerminal">Clear Terminal</button>
  </div>
</template>

<script setup lang="ts">
import "xterm/css/xterm.css";
import { Terminal } from "xterm";

const terminal = ref<Terminal | null>(null);
const terminalContainer = ref<HTMLElement | null>(null);

const handleInput = (data: string) => {
  // Echo the user input back to the terminal
  terminal.value?.write(data + "\r\n");

  // You can perform custom logic here based on user input
  // For simplicity, we're just echoing back the input.
};

const clearTerminal = () => {
  terminal.value?.clear();
};

onMounted(() => {
  // Create a new terminal instance
  terminal.value = new Terminal();

  // Attach the terminal to the DOM element
  terminal.value.open(terminalContainer.value!);

  // Write a welcome message to the terminal
  terminal.value.write("Welcome to the xterm.js example!\r\n");

  // Add an event listener to handle user input
  terminal.value.onData((data) => handleInput(data));
});
</script>

<style>
/* Add any custom styles here if needed */
</style>
