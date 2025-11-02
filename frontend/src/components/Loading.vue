<template>
  <div class="flex justify-center">
    <svg class="spinner-smooth text-primary1" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
      <path fill="currentColor" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity="0.5" />
      <path fill="currentColor" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z">
        <animateTransform attributeName="transform" dur="0.8s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate" keySplines="0.4 0 0.2 1" calcMode="spline" />
      </path>
    </svg>
  </div>
</template>

<style scoped>
.spinner-smooth {
  will-change: transform;
  /* fallback for browsers that support CSS animation on SVG */
  animation: spinner-rotate 0.8s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}
@keyframes spinner-rotate {
  100% {
    transform: rotate(360deg);
  }
}
</style>
