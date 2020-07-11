<template>
    <div class="app-item-list">
        <div class="list-container" v-loading="loading">
            <app-item
                    v-for="(item, index) of list"
                    :key="index"
                    :data="item"
                    :use-image-cache="use_image_cache"
                    @show-detail="onShowDetail"/>
            <el-pagination
                    v-if="this.list.length > 0"
                    background
                    layout="prev, pager, next, ->, jumper"
                    :page-count="max_page"
                    @current-change="onPageChange"/>
        </div>
        <el-dialog :title="detail.title" :visible.sync="showDetailDialog">
            <div style="height: 600px; overflow: auto;">
                <div>
                    <div>
                        <span class="detail-field">Version: </span>
                        <span class="detail-value">{{detail.version}}</span>
                    </div>
                    <div>
                        <span class="detail-field">Size: </span>
                        <span class="detail-value">{{detail.size}}</span>
                    </div>
                    <div>
                        <span class="detail-field">Posted: </span>
                        <span class="detail-value">{{detail.date_published}}</span>
                    </div>
                    <download-button v-if="detail.urls" :urls="detail.urls"/>
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
                                             :title="version.version"
                                             :urls="version.urls"/>
                        </div>
                    </template>

                    <div class="detail-content" v-html="detail.content"></div>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script>
    import AppItem from "./AppItem";
    import DownloadButton from "./DownloadButton";

    export default {
        name: "AppItemList",
        components: {DownloadButton, AppItem},
        data() {
            return {
                loading: true,
                previousLoading: false,
                use_image_cache: true,
                category: '',
                page: 1,
                max_page: 1,
                length: 0,
                list: [],
                showDetailDialog: false,
                detail: {},
                previousVersionLoaded: false,
                previousVersions: [],
            };
        },
        async mounted() {
            await this.refresh();
        },
        methods: {
            async refresh(params = {}) {
                this.loading = true;
                try {
                    let res = await this.axios.get("/api/list", {params});
                    this.use_image_cache = res.data.data.use_image_cache;
                    this.category = res.data.data.category;
                    this.page = res.data.data.page;
                    this.max_page = res.data.data.max_page;
                    this.length = res.data.data.length;
                    this.list = res.data.data.list;
                } catch (e) {
                }
                this.loading = false;
            },
            async onShowDetail(detail) {
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
            async onPageChange(page) {
                let params = {};
                if (this.category !== '') {
                    params.category = this.category;
                }
                if (page !== 1) {
                    params.page = page;
                }
                await this.refresh(params);
            }
        }
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
    .list-container {
        width: 100%;
    }

    .el-pagination {
        margin-top: 16px;
        margin-bottom: 16px;
    }

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