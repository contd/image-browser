<template>
  <Moveable
    class="item moveable"
    v-bind="moveable"
    @drag="handleDrag"
    @resize="handleResize"
    @scale="handleScale"
    @rotate="handleRotate"
    @warp="handleWarp"
  >
    <span>
      <img
        :src="path"
        :alt="name"
        :data-size="size"
        :data-type="type"
        :data-modified="modified"
        :data-width="width"
        :data-height="height"
        :data-lat="lat"
        :data-long="long"
        :width="setWidth(width, height)"
        :height="setHeight(width, height)"
      />
    </span>
  </Moveable>
</template>

<script>
import Moveable from 'vue-moveable'

export default {
  name: "PictureIcon",
  components: {
    Moveable
  },
  data: () => ({
    moveable: {
      target: document.querySelector(".pictures"),
      container: document.querySelector(".pictures"),
      origin: false,
      draggable: true,
      throttleDrag: 0,
      resizable: false,
      throttleResize: 1,
      keepRatio: true,
      scalable: true,
      throttleScale: 0,
      rotatable: false,
      throttleRotate: 0
    }
  }),
  props: {
    path: String,
    name: String,
    size: Number,
    type: String,
    modified: String,
    width: Number,
    height: Number,
    lat: Number,
    long: Number,
    clickImage: Function
  },
  methods: {
    setWidth(width, height) {
      return width >= height ? "200px" : ""
    },
    setHeight(width, height) {
      return height >= width ? "200px" : ""
    },
    handleDrag({ target, transform }) {
      //console.log("onDrag left, top", left, top)
      //target.style.left = `${left}px`
      //target.style.top = `${top}px`
      target.style.transform = transform
    },
    handleResize({ target, width, height, delta }) {
      //console.log("onResize", width, height)
      delta[0] && (target.style.width = `${width}px`)
      delta[1] && (target.style.height = `${height}px`)
    },
    handleScale({ target, transform, scale }) {
      console.log("onScale scale", scale)
      target.style.transform = transform
    },
    handleRotate({ target, dist, transform }) {
      console.log("onRotate", dist)
      target.style.transform = transform
    },
    handleWarp({ target, transform }) {
      //console.log("onWarp", target)
      target.style.transform = transform
    }
  }
}
</script>
