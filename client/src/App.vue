<template>
  <div>
    <NavBar></NavBar>
    <b-container fluid>
      <b-row>
        <b-col>
          <button
            v-if="!inRoot"
            class="item ffolder small gray"
            data-path="/api"
            alt="Go up/back one directory"
            v-on:click="goUpOneDir"
          >
            <span>..</span>
          </button>
          <button
            class="item ffolder small cyan"
            v-for="directory in directories"
            :key="directory.name"
            :data-path="directory.path"
            :data-name="directory.name"
            :data-size="directory.size"
            :data-modified="directory.modified"
            @click="() => changeDir(directory.path)"
          >
            <span class="truncate">{{ directory.name }}</span>
          </button>
        </b-col>
      </b-row>
      <b-row>
        <b-col class="pictures">
          <PictureIcon
            v-for="(picture, index) in pictures"
            :key="index"
            :name="picture.name"
            :path="picture.path"
            :size="picture.size"
            :type="picture.type"
            :modified="picture.modified"
            :width="picture.width"
            :height="picture.height"
            :lat="picture.exif.lat"
            :long="picture.exif.long"
          ></PictureIcon>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import NavBar from '@/components/NavBar'
import PictureIcon from '@/components/PictureIcon'

export default {
  name: "app",
  components: {
    NavBar,
    PictureIcon
  },
  data() {
    return {
      pictures: [],
      directories: [],
      picsArray: [],
      currPath: "",
      inRoot: true,
      visible: false,
      loading: true,
      errored: false
    }
  },
  mounted() {
    this.$axios
      .get(this.$BASEURL)
      .then(resp => {
        this.pictures = resp.data.Pictures
        this.directories = resp.data.Directories
        this.picsArray = resp.data.Pictures
          ? resp.data.Pictures.map(pic => pic.path)
          : []
      })
      .catch(err => {
        console.log(`ERROR: ${err}`)
        this.errored = true
      })
      .finally(() => (this.loading = false))
  },
  methods: {
    changeDir(path) {
      if (path.match(/\?path=/)) {
        this.inRoot = false
      } else {
        this.inRoot = true
      }
      this.currPath = path
      this.$axios
        .get(path)
        .then(resp => {
          this.pictures = resp.data.Pictures
          this.directories = resp.data.Directories
          this.picsArray = resp.data.Pictures
            ? resp.data.Pictures.map(pic => pic.path)
            : []
        })
        .catch(err => {
          console.log(`ERROR: ${err}`)
          this.errored = true
        })
        .finally(() => (this.loading = false))
    },
    goUpOneDir() {
      const pathStr = this.currPath.split("?")[1]
        ? this.currPath.split("?")[1]
        : ""
      let newPath = ""
      let path = pathStr.replace(/path=/, "").split("/")
      if (path.length > 1) {
        path.pop()
        newPath = `${this.$BASEURL}?path=${path.join("/")}`
      } else {
        newPath = `${this.$BASEURL}`
      }
      this.changeDir(newPath)
    }
  }
}
</script>

<style lang="css">
.truncate {
  width: 50px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.truncate:hover {
  width: 100px;
  white-space: normal;
  overflow: visible;
  text-overflow: unset;
}
</style>