<template>
    <el-dropdown class="download-button"
                 v-loading="downloadLoading"
                 @command="onDownload">
        <el-button type="primary" plain size="small">
            {{title}}<i class="el-icon-arrow-down el-icon--right"></i>
        </el-button>
        <el-dropdown-menu slot="dropdown">
            <el-dropdown-item
                    v-for="(link, index) of urls"
                    :index="index"
                    :command="link.url">
                {{link.title}}
            </el-dropdown-item>
        </el-dropdown-menu>
    </el-dropdown>
</template>

<script>
    export default {
        name: "DownloadButton",
        props: {
            title: {
                type: String,
                default: 'Download'
            },
            urls: {
                type: Array,
                required: true,
            },
        },
        data() {
            return {
                downloadLoading: false,
            }
        },
        methods: {
            async onDownload(url) {
                this.downloadLoading = true;
                try {
                    let res = await this.axios.get("/api/direct_url", {params: {url}});
                    this.downloadLoading = false;
                    window.open(res.data.data, '_blank').location;
                } catch (e) {
                }
                this.downloadLoading = false;
            },
        }
    }
</script>

<style scoped>
    .download-button {
        margin-top: 12px;
    }
</style>