<template>
    <div id="app">
        <el-container>
            <el-header>
                <img class="logo" src="./assets/nmac.png" alt="nmac.to">
                <div class="title">NMac Mirror</div>
                <el-input
                        class="search-input"
                        placeholder="Search..."
                        clearable
                        suffix-icon="el-icon-search"
                        size="small"
                        v-model="searchText"
                        :disabled="listLoading"
                        @keyup.native.enter="onSearch">
                </el-input>
            </el-header>
            <el-container>
                <el-aside>
                    <category-menu @select="onSelectCategory" :category="category"/>
                </el-aside>
                <el-container v-loading="listLoading">
                    <el-main>
                        <app-item-list ref="itemList" @loading="onLoadingChange"/>
                    </el-main>
                    <el-footer>Copyright (C) 2020 xiaozhuai.</el-footer>
                </el-container>
            </el-container>
        </el-container>
    </div>
</template>

<script>
    import CategoryMenu from "./components/CategoryMenu";
    import AppItemList from "./components/AppItemList";

    export default {
        name: 'app',
        components: {AppItemList, CategoryMenu},
        data() {
            return {
                listLoading: true,
                category: '',
                searchText: '',
            }
        },
        watch: {
            searchText(n, o) {
                if (n.trim() === '' && o.trim() !== '') {
                    this.$refs.itemList.onClearSearchText();
                }
            }
        },
        mounted() {
            let initParams = this.$router.get();
            if (initParams.isSearchMode) {
                this.searchText = initParams.params.s;
                this.$refs.itemList.search(initParams.params);
            } else {
                this.category = initParams.params.category;
                this.$refs.itemList.refresh(initParams.params);
            }
        },
        methods: {
            onLoadingChange(loading) {
                this.listLoading = loading;
            },
            onSelectCategory(category) {
                this.category = category;
                let params = {};
                if (category !== '') {
                    params.category = category;
                }
                this.$refs.itemList.refresh(params);
            },
            onSearch() {
                if (this.searchText.trim() === '') {
                    return;
                }
                let params = {
                    s: this.searchText
                };
                this.$refs.itemList.search(params);
            }
        },
    }
</script>

<style>
    html, body {
        margin: 0;
        padding: 0;
        width: 100%;
        height: 100%;
    }

    #app {
        width: 100%;
        height: 100%;
        font-family: 'Avenir', Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
    }

    .logo {
        display: inline-block;
        vertical-align: top;
        height: 48px;
        margin: 6px;
    }

    .title {
        display: inline-block;
        vertical-align: top;
        color: #eeeeee;
        line-height: 60px;
        margin-left: 24px;
    }

    .search-input {
        float: right;
        width: 280px;
        margin-top: 14px;
    }

    #app > .el-container {
        width: 100%;
        height: 100%;
    }

    #app > .el-container > .el-container {
        width: 100%;
        height: calc(100% - 60px);
    }

    #app .el-header {
        background-color: #303030;
    }

    #app .el-footer {
        border-top: 1px solid #dddddd;
        box-sizing: border-box;
        line-height: 60px;
        color: #202020;
        font-size: 13px;
    }

    .category-menu {
        height: 100%;
    }

    .app-item-list {
        width: 100%;
        height: 100%;
    }
</style>
