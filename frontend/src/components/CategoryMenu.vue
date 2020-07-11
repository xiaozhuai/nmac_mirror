<template>
    <el-menu class="category-menu"
             v-loading="loading"
             background-color="#303030"
             text-color="#fff"
             active-text-color="#ffd04b"
             @select="onSelect">
        <el-menu-item
                v-for="(category, index) of categories"
                :key="index"
                :index="category.category">
            {{category.title}}
        </el-menu-item>
    </el-menu>
</template>


<script>
    export default {
        name: "CategoryMenu",
        data() {
            return {
                loading: true,
                categories: []
            }
        },
        async mounted() {
            await this.loadCategories();
        },
        methods: {
            async loadCategories() {
                this.loading = true;
                try {
                    let res = await this.axios.get("/api/categories");
                    this.categories = res.data.data;
                } catch (e) {
                }
                this.loading = false;
            },
            onSelect(key, keyPath) {
                this.$emit("select", key);
            },
        },
    }
</script>

<style scoped>

</style>