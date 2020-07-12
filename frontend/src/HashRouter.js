export default class HashRouter {
    update(isSearchMode, params) {
        if (isSearchMode) {
            let {s, page} = params;
            let hash = '/search';
            hash += `/${encodeURIComponent(s)}`;
            if (page === undefined || page === 1) {
            } else {
                hash += `/${page}`
            }
            window.location.hash = hash;
        } else {
            let {category, page} = params;
            let hash = '';
            if (category === undefined || category === "") {
            } else {
                hash += `/category/${category}`
            }
            if (page === undefined || page === 1) {
            } else {
                hash += `/${page}`
            }
            window.location.hash = hash;
        }
    }

    get() {
        let hash = window.location.hash;
        if (hash.startsWith('#')) {
            hash = hash.substr(1);
        }
        if (hash.startsWith('/')) {
            hash = hash.substr(1);
        }

        let ret = {
            isSearchMode: false,
            params: {},
        };

        hash = hash.trim();

        if (hash !== '') {
            let tokens = hash.split('/');

            if (tokens.length > 0) {
                if (tokens[0] === 'search') {
                    ret.isSearchMode = true;
                    if (tokens.length > 1) {
                        ret.params.s = decodeURIComponent(tokens[1]);
                        if (tokens.length > 2) {
                            ret.params.page = parseInt(tokens[2]);
                        }
                    }
                } else if (tokens[0] === 'category') {
                    ret.isSearchMode = false;
                    if (tokens.length > 1) {
                        ret.params.category = tokens[1];
                        if (tokens.length > 2) {
                            ret.params.page = parseInt(tokens[2]);
                        }
                    }
                } else {
                    ret.isSearchMode = false;
                    ret.params.page = parseInt(tokens[0]);
                }
            }
        }

        return ret;
    }
}