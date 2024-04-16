<template>
  <a v-bind="$attrs" v-on="$listeners"><span class="text"><slot/></span>
    <v-icon x-small>mdi-open-in-new</v-icon>
  </a>
</template>

<script>
export default {
  name: "ANew"
}
</script>

<style scoped>
a {
  text-decoration: none;
}

.text {
  text-decoration: underline;
}

.v-icon {
  color: inherit;
}
</style>
