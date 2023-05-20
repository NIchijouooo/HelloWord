<template>
  <div v-if="isUrl" :style="urlIconStyle" class="url-svg svg-icon icon" />
  <svg v-else class="svg-icon icon" :style="iconStyle">
    <use :href="iconName" />
  </svg>
</template>

<script setup>
import { isExternal } from 'utils/common'
const props = defineProps({
  name: {
    type: String,
    default: '',
  },
  size: {
    type: String,
    default: '18px',
  },
  color: {
    type: String,
    default: '#000000',
  },
})
const s = `${props.size.replace('px', '')}px`
const iconName = computed(() => `#${props.name}`)
const iconStyle = computed(() => {
  return {
    color: props.color,
    fontSize: s,
  }
})
const isUrl = computed(() => isExternal(props.name))
const urlIconStyle = computed(() => {
  return {
    width: s,
    height: s,
    mask: `url(${props.name}) no-repeat 50% 50%`,
    '-webkit-mask': `url(${props.name}) no-repeat 50% 50%`,
  }
})
</script>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  fill: currentColor;
  overflow: hidden;
}
</style>
