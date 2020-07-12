<template>
    <div class="app-item-list">
        <div class="list-container" v-if="list.length > 0">
            <div class="list-container-inner">
                <app-item
                        v-for="(item, index) of list"
                        :key="index"
                        :data="item"
                        :use-image-cache="use_image_cache"
                        @show-detail="onShowDetail"/>
            </div>
        </div>
        <el-pagination
                v-if="list.length > 0"
                background
                layout="prev, pager, next, ->, jumper"
                :current-page="page"
                :page-count="max_page"
                @current-change="onPageChange"/>
        <div v-if="nothing" class="nothing">{{isSearchMode ? 'Nothing Found' : 'Nothing Here'}}</div>
        <detail-dialog ref="detailDialog"/>
    </div>
</template>

<script>
    import AppItem from "./AppItem";
    import DownloadButton from "./DownloadButton";
    import DetailDialog from "./DetailDialog";

    export default {
        name: "AppItemList",
        components: {DetailDialog, DownloadButton, AppItem},
        data() {
            return {
                loading: true,
                isSearchMode: false,
                nothing: false,
                normalResult: {
                    category: '',
                    page: 1,
                    use_image_cache: true,
                    max_page: 1,
                    length: 0,
                    list: [],
                },
                searchResult: {
                    searchText: '',
                    page: 1,
                    use_image_cache: true,
                    max_page: 1,
                    length: 0,
                    list: [],
                },
            };
        },
        computed: {
            page() {
                return this.isSearchMode ? this.searchResult.page : this.normalResult.page;
            },
            use_image_cache() {
                return this.isSearchMode ? this.searchResult.use_image_cache : this.normalResult.use_image_cache;
            },
            max_page() {
                return this.isSearchMode ? this.searchResult.max_page : this.normalResult.max_page;
            },
            length() {
                return this.isSearchMode ? this.searchResult.length : this.normalResult.length;
            },
            list() {
                return this.isSearchMode ? this.searchResult.list : this.normalResult.list;
            },
        },
        methods: {
            setLoading(b) {
                this.loading = true;
                this.$emit('loading', b);
            },
            async refresh(params = {}) {
                this.isSearchMode = false;
                this.$router.update(this.isSearchMode, params)
                this.setLoading(true);
                try {
                    let res = await this.axios.get("/api/list", {params});
                    this.normalResult = res.data.data;
                    this.nothing = this.normalResult.length === 0;
                } catch (e) {
                }
                this.setLoading(false);
            },
            async search(params = {}) {
                this.isSearchMode = true;
                this.$router.update(this.isSearchMode, params)
                this.setLoading(true);
                try {
                    let res = await this.axios.get("/api/search", {params});
                    this.searchResult = res.data.data;
                    this.nothing = this.searchResult.length === 0;
                } catch (e) {
                }
                this.setLoading(false);
            },
            onClearSearchText() {
                this.isSearchMode = false;
                this.refresh({
                    category: this.normalResult.category,
                    page: this.normalResult.page,
                });
            },
            async onShowDetail(detail) {
                this.$refs.detailDialog.show(detail);
            },
            async onPageChange(page) {
                if (this.isSearchMode) {
                    let params = {};
                    if (this.searchResult.searchText !== '') {
                        params.s = this.searchResult.searchText;
                    }
                    if (page !== 1) {
                        params.page = page;
                    }
                    await this.search(params);
                } else {
                    let params = {};
                    if (this.normalResult.category !== '') {
                        params.category = this.normalResult.category;
                    }
                    if (page !== 1) {
                        params.page = page;
                    }
                    await this.refresh(params);
                }
            }
        }
    }
</script>

<style scoped>
    .list-container {
        width: 100%;
        max-width: 1200px;
        margin: auto;
        height: calc(100% - 48px);
        overflow: auto;
    }

    .list-container::-webkit-scrollbar {
        display: none
    }

    .list-container-inner {
        width: 100%;
    }

    .el-pagination {
        margin: 16px auto 0;
        max-width: 1200px;
    }

    .nothing {
        text-align: center;
        color: #202020;
    }
</style>