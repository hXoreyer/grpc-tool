<script setup lang="ts">

</script>

<template>
    <router-view />
</template>

<style>
body {
  min-height: calc(100vh - 6px);
  min-width: calc(100vw - 6px);
  height: calc(100vh - 6px);
  width: calc(100vw - 6px);
  padding: 0;
  margin: 0;
  background: rgba(0, 0, 0, 0);
}

#app {
  margin-left: 3px;
  margin-top: 3px;
  font-family: Consolas, 'Courier New', monospace;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  background: white;
  border-radius: 5px;
  box-shadow: 2px 2px 2px rgba(1, 1, 1, 0.1),-2px -2px 2px rgba(1, 1, 1, 0.1);
}

</style>
