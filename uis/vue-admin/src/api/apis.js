import axios from 'axios';
import { getToken,setToken,removeToken } from '@/util/storage';

const apiUrl=process.env.NODE_ENV === 'production' ? '' : "http://localhost:8030";
// const apiUrl='http://open.vkstu.com:8050';

const serv=axios.create({
    baseURL: apiUrl, // api base_url
    // baseURL: 'http://localhost:8082', // api base_url
    // timeout: 5000, // 请求超时时间
    withCredentials: true
});

const post=function(path,params,headers){
    let hds={};
    if(headers)hds=headers;
    let token = getToken();
    if(token)hds['Authorization']='TOKEN '+token;
    return serv.post(path,params,{headers:hds});
}

export default{
    post,apiUrl
}