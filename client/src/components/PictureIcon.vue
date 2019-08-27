<template>
    <div>
        <img 
            class="item item-icon"
            :src="path" 
            :alt="name" 
            :data-width="width" 
            :data-height="height"
            :width="width"
            :hight="height"
        />
    </div>
</template>

<script>
export default {
    name: 'PictureIcon',
    props: {
        path: String,
        name: String,
        width: String,
        height: String
    },
    mounted() {
        if (this.width > this.height) {
            this.width = "200px"
            this.height = null
        } else {
            this.height = "200px"
            this.width = null
        }
    }
}
</script>
