export const INCREMENT=(state) =>{
    state.count++
}
export const DECREMENT=(state) =>{
    state.count--
}
export const setLgInfo=(state,par) =>{
    state.userinfo=par;
}