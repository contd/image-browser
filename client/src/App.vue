<template>
  <div id="app" class="icon-grid">
    <div v-if="loading">Loading...</div>
    <div v-else>
      <PictureIcon v-for="file in files" :key="file.name"
        path="file.path"
        name="file.name"
        width="file.width"
        height="file.height"
      ></PictureIcon>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import PictureIcon from './components/PictureIcon'

export default {
  name: 'app',
  components: {
    PictureIcon
  },
  data() {
    return {
      files: null,
      loading: true,
      errored: false
    }
  },
  mounted() {
    axios
      .get('http://localhost:6969/')
      .then(resp => {
        this.files = resp.data
      })
      .catch(err => {
        console.log(`ERROR: ${err}`)
        this.errored = true
      })
      .finally(() => this.loading = false)
  }
}
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

ul {
	padding: 1em;
	margin: 1em;
}

li {
	padding: 1em;
	margin: 1em;
	justify-content: center;
	text-align: center;
	display: inline-flex;
	font-size: 2em;
	font-weight: 500;
}

.icon-grid {
    display: flexbox;
    margin: 0;
    padding: 0;
}

.item {
    margin: 1px;
    padding: 0;
}

.item-icon {
    justify-content: center;
	text-align: center;
	display: inline-flex;
}

.zoomwall {
	font-size: 0;
	overflow: hidden;
	margin: 0;
	padding: 0;
}

.zoomwall img {
	height: 15vw;
	opacity: 1;
	vertical-align: top;
	
	transform-origin: 0% 0%;
	transition-property: transform, opacity;
	transition-duration: 0.3s;
	transition-timing-function: ease-out;

	-webkit-transform-origin: 0% 0%;
	-webkit-transition-property: transform, opacity;
	-webkit-transition-duration: 0.3s;
	-webkit-transition-timing-function: ease-out;
}

.zoomwall.lightbox img {
	transition-timing-function: ease-in;
	-webkit-transition-timing-function: ease-in;
}

.zoomwall.lightbox img {
	opacity: 0.3;
}

.zoomwall.lightbox img.active {
	opacity: 1;
}
</style>
