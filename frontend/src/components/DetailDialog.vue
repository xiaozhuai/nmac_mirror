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
                <download-button class="detail-btn" v-if="detail.urls && detail.urls.length > 0" :urls="detail.urls"/>
                <el-button v-else class="detail-download detail-btn" type="primary" plain size="small">
                    No Resource!
                </el-button>
                <el-button v-if="detail.previous_page_url && !previousVersionLoaded"
                           v-loading="previousLoading"
                           class="detail-btn detail-download download-button"
                           type="primary" plain size="small"
                           @click="onClickPreviousVersion">
                    Previous Version
                </el-button>
                <template v-if="previousVersionLoaded">
                    <el-button v-if="previousVersions.length === 0" class="detail-download detail-btn" type="primary" plain size="small">
                        No Previous Version!
                    </el-button>
                    <download-button class="detail-btn" v-for="(version, index) of previousVersions"
                                     :key="index"
                                     :title="version.version"
                                     :urls="version.urls"/>
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

    .detail-content {
        margin-top: 24px;
    }

    .detail-content img {
        max-width: 100%;
    }

    .detail-content .label-important {
        padding: 4px 16px;
        border: 1px solid #dddddd;
        border-radius: 6px;
        color: white;
        background-color: #dd6161;
    }

    .detail-content blockquote {
        border-color: #e3e3e3;
        padding: 0 0 0 15px;
        margin: 0 0 20px;
        border-left: 5px solid #eee;
    }

    .detail-content blockquote p {
        font-size: 12.5px;
        font-weight: 200;
        line-height: 1.25;
    }

    .detail-content code {
        padding: 2px 4px;
        color: #d14;
        white-space: nowrap;
        background-color: #f7f7f9;
        border: 1px solid #e1e1e8;
    }

    code, pre {
        padding: 0 3px 2px;
        font-family: Monaco,Menlo,Consolas,courier new,monospace;
        font-size: 12px;
        color: #333;
        -webkit-border-radius: 3px;
        -moz-border-radius: 3px;
        border-radius: 3px;
    }

    .el-dialog {
        width: 80%;
        max-width: 1000px;
    }

    .detail-icon {
        width: 120px;
        height: 120px;
        float: left;
        margin-right: 16px;
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

    .detail-btn {
        margin-right: 12px;
    }
</style>