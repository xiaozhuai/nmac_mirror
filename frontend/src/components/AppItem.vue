<template>
    <div class="app-item" v-loading="loading">
        <div class="app-icon-container">
            <el-image class="app-icon"
                      :src="imageUrl"
                      lazy
                      @click="onShowDetail"/>
        </div>
        <div class="app-right">
            <h1 class="app-title" @click="onShowDetail">{{data.title}}</h1>
            <p class="app-description">{{data.description}}</p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "AppItem",
        props: {
            data: {
                type: Object,
                required: true,
            },
            useImageCache: {
                type: Boolean,
                default: true,
            },
        },
        data() {
            return {
                loading: false,
            };
        },
        computed: {
            imageUrl() {
                return this.useImageCache ? `/api/fetch_image?url=${encodeURIComponent(this.data.image_url)}` : this.data.image_url;
            }
        },
        methods: {
            async onShowDetail() {
                this.loading = true;
                try {
                    let res = await this.axios.get("/api/detail", {
                        params: {url: this.data.detail_page_url}
                    });
                    this.$emit('show-detail', res.data.data);
                } catch (e) {
                }
                this.loading = false;
            },
        },
    }
</script>

<style scoped>
    .app-item {
        border: 1px solid #dddddd;
        box-sizing: border-box;
        border-radius: 8px;
        padding: 8px;
    }

    .app-item + .app-item {
        margin-top: 16px;
    }

    .app-item:last-child {
        margin-bottom: 16px;
    }

    .app-icon-container {
        width: 100px;
        height: 100px;
        display: inline-block;
        vertical-align: top;
        overflow: hidden;
    }

    .app-icon {
        width: 100px;
        height: 100px;
        display: inline-block;
        vertical-align: top;
        cursor: pointer;
        transform: scale(1.02);
    }

    .app-icon:hover {
        opacity: 0.8;
    }

    .app-right {
        width: calc(100% - 100px);
        height: 100px;
        display: inline-block;
        vertical-align: top;
    }

    .app-title {
        padding-left: 12px;
        margin-top: 4px;
        margin-bottom: 16px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-size: 16px;
        cursor: pointer;
    }

    .app-title:hover {
        opacity: 0.7;
        text-decoration: underline;
    }

    .app-description {
        padding-left: 12px;
        margin: 0;
        font-size: 13px;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 3;
    }
</style>