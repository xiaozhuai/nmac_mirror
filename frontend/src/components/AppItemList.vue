<template>
    <div class="app-item-list">
        <div class="list-container">
            <app-item
                    v-for="(item, index) of list"
                    :key="index"
                    :data="item"
                    :use-image-cache="use_image_cache"
                    @show-detail="onShowDetail"/>
            <el-pagination
                    v-if="list.length > 0"
                    background
                    layout="prev, pager, next, ->, jumper"
                    :current-page="page"
                    :page-count="max_page"
                    @current-change="onPageChange"/>
        </div>
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
                isSearchMode: false,
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
        async mounted() {
            await this.refresh();
        },
        methods: {
            setLoading(b) {
                this.$emit('loading', b);
            },
            async refresh(params = {}) {
                this.setLoading(true);
                try {
                    let res = await this.axios.get("/api/list", {params});
                    this.normalResult = res.data.data;
                    this.isSearchMode = false;
                } catch (e) {
                }
                this.setLoading(false);
            },
            async search(params = {}) {
                this.setLoading(true);
                try {
                    let res = await this.axios.get("/api/search", {params});
                    this.searchResult = res.data.data;
                    this.isSearchMode = true;
                } catch (e) {
                }
                this.setLoading(false);
            },
            onClearSearchText() {
                this.isSearchMode = false;
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
    }

    .el-pagination {
        margin-top: 16px;
        margin-bottom: 16px;
    }
</style>