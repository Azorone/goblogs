import axios from "axios";

const Server = axios.create({
    baseURL: "/api",
    timeout: 1000,
});

Server.interceptors.request.use(
    (config) => {

         let token = localStorage.getItem("token");
        if (token!=null) {
            config.headers["Authorization"] =  `${localStorage.getItem('token')}`;


        }
        return config;

    },
    (error) => {

        return Promise.reject(error);
    }
);
Server.interceptors.response.use(
    (response) => {
        if (response.status == 200) {
            if (response.data.code == 500) {
             //错误处理
            }

        }
        return response.data;
    },
    (error) => {

        return Promise.reject(error);
    }
);

const API = {
    register: (data) => {
        return Server.post('/register', data);
    },
    login: (data) => {
        return Server.post('/login', data);
    },
    addBlog: (data) => {
        return Server.post('/blogs', data);
    },
    updateBlog: (id, data) => {
        return Server.put(`/blogs/${id}`, data);
    },
    deleteBlog: (id) => {
        return Server.delete(`/blogs/${id}`);
    },
    addCategory: (data) => {
        return Server.post('/categories', data);
    },
    updateCategory: (id, data) => {
        return Server.put(`/categories/${id}`, data);
    },
    getBlogsManger: (page, pageSize) => {
        return Server.get("/blogs/manger", {
            params: {
                page,
                pageSize,
            },
        });
    },
    getBlogs: (page, pageSize) => {
        return Server.get("/blogs", {
            params: {
                page,
                pageSize,
            },
        });
    },
    getCategories: (page, pageSize) => {
        return Server.get("/categories", {
            params: {
                page,
                pageSize,
            },
        });
    },
    getBlogsByCategory: (categoryID, page, pageSize) => {
        return Server.get(`/blogs/category/${categoryID}`, {
            params: {
                page,
                pageSize,
            },
        });
    },
    getBlogById: (id) => {
        return Server.get(`/blogs/${id}`);
    },
};

export default API;
