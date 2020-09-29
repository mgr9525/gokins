
export const getToken=()=>{
    return sessionStorage.getItem('token');
}
export const setToken=(tks)=>{
    sessionStorage.setItem('token',tks);
}
export const removeToken=()=>{
    sessionStorage.removeItem('token');
}