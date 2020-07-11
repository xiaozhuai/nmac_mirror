<template>
    <div id="app">
        <el-container>
            <el-header>
                <img class="logo" src="./assets/nmac.png" alt="nmac.to">
                <div class="title">NMac Mirror</div>
            </el-header>
            <el-container>
                <el-aside>
                    <category-menu @select="onSelectCategory"/>
                </el-aside>
                <el-container>
                    <el-main>
                        <app-item-list ref="itemList"/>
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
        methods: {
            onSelectCategory(category) {
                let params = {};
                if (category !== '') {
                    params.category = category;
                }
                this.$refs.itemList.refresh(params)
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
    }

    .category-menu {
        height: 100%;
    }

    .app-item-list {
        width: 100%;
        height: 100%;
    }
</style>
