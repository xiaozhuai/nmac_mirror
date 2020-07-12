<template>
    <el-dialog :title="detail.title" :visible.sync="showDetailDialog" :destroy-on-close="true">
        <div style="height: 600px; overflow: auto;">
            <div>
                <div v-if="detail.version">
                    <span class="detail-field">Version: </span>
                    <span class="detail-value">{{detail.version}}</span>
                </div>
                <div v-if="detail.size">
                    <span class="detail-field">Size: </span>
                    <span class="detail-value">{{detail.size}}</span>
                </div>
                <div v-if="detail.date_published">
                    <span class="detail-field">Posted: </span>
                    <span class="detail-value">{{detail.date_published}}</span>
                </div>
                <download-button v-if="detail.urls && detail.urls.length > 0" :urls="detail.urls"/>
                <el-button v-else class="detail-download" type="primary" plain size="small">
                    No Resource!
                </el-button>
                <br>
                <el-button v-if="detail.previous_page_url && !previousVersionLoaded"
                           v-loading="previousLoading"
                           class="detail-download"
                           type="primary" plain size="small"
                           @click="onClickPreviousVersion">
                    Previous Version
                </el-button>
                <template v-if="previousVersionLoaded">
                    <span v-if="previousVersions.length === 0" style="margin-top: 12px; display: block;">No Previous Version!</span>
                    <div class="previous-version-container" v-else>
                        <download-button v-for="(version, index) of previousVersions"
                                         :key="index"
                                         :title="version.version"
                                         :urls="version.urls"/>
                    </div>
                </template>

                <div class="detail-content" v-html="detail.content"></div>
            </div>
        </div>
    </el-dialog>
</template>

<script>
    import DownloadButton from "./DownloadButton";

    export default {
        name: "DetailDialog",
        components: {DownloadButton},
        data() {
            return {
                previousLoading: false,
                showDetailDialog: false,
                detail: {},
                previousVersionLoaded: false,
                previousVersions: [],
            };
        },
        methods: {
            show(detail) {
                this.previousVersionLoaded = false;
                this.previousVersions = [];
                this.detail = detail;
                this.showDetailDialog = true;
            },
            async onClickPreviousVersion() {
                this.previousLoading = true;
                try {
                    let res = await this.axios.get("/api/previous_version", {
                        params: {url: this.detail.previous_page_url}
                    });
                    this.previousVersionLoaded = true;
                    this.previousVersions = res.data.data;
                } catch (e) {
                }
                this.previousLoading = false;
            },
        },
    }
</script>

<style>
    .el-dialog__header {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .el-dialog__body {
        padding: 12px 12px;
    }

    .detail-content img {
        max-width: 100%;
    }

    .el-dialog {
        width: 80%;
        max-width: 1000px;
    }

    .detail-icon {
        width: 120px;
        height: 120px;
        margin-top: 16px;
        margin-bottom: 4px;
    }
</style>

<style scoped>
    .detail-field,
    .detail-value {
        color: #202020;
        display: inline-block;
        vertical-align: top;
        line-height: 28px;
    }

    .detail-field {
        width: 64px;
        padding-right: 12px;
        text-align: right;
    }

    .detail-download {
        margin-top: 12px;
    }

    .previous-version-container .download-button {
        margin-right: 12px;
    }
</style>