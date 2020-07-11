import Vue from 'vue';
import axios from 'axios';
import VueAxios from 'vue-axios';
import App from './App.vue';

import {
    Aside,
    Button,
    Container,
    Dialog,
    Dropdown,
    DropdownMenu,
    DropdownItem,
    Footer,
    Header,
    Loading,
    Main,
    Menu,
    MenuItem,
    MenuItemGroup,
    Message,
    MessageBox,
    Notification,
    Pagination,
    Submenu,
} from 'element-ui';

axios.interceptors.response.use(response => {
    if (response.data.code !== 0) {
        let err = new Error(response.data.msg);
        err.response = response;
        return Promise.reject(err);
    }
    return response;
}, error => {
    if (error.response.data.msg) {
        error.message = error.response.data.msg;
    }
    console.error(error);
    Message.error(error.message);
    return Promise.reject(error);
});

Vue.use(VueAxios, axios);

Vue.use(Aside);
Vue.use(Button);
Vue.use(Container);
Vue.use(Dialog);
Vue.use(Dropdown);
Vue.use(DropdownMenu);
Vue.use(DropdownItem);
Vue.use(Footer);
Vue.use(Header);
Vue.prototype.$loading = Loading.service;
Vue.use(Loading.directive);
Vue.use(Main);
Vue.use(Menu);
Vue.use(MenuItem);
Vue.use(MenuItemGroup);
Vue.prototype.$message = Message;
Vue.prototype.$msgbox = MessageBox;
Vue.prototype.$alert = MessageBox.alert;
Vue.prototype.$confirm = MessageBox.confirm;
Vue.prototype.$prompt = MessageBox.prompt;
Vue.prototype.$notify = Notification;
Vue.use(Pagination);
Vue.use(Submenu);


Vue.config.productionTip = false;

new Vue({
    render: h => h(App),
}).$mount('#app');
