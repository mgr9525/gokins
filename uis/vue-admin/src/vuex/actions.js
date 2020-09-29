//test
export const increment = ({commit}) => {
    commit('INCREMENT')
}
export const decrement = ({commit}) => {
    commit('DECREMENT')
}


export const getLgInfo = ({commit}) => {
    return new Promise((resolve, reject)=>{
        api.post('/api/lginfo').then(res=>{
            commit('setLgInfo',res.data);
            resolve();
        }).catch(err=>{
            reject(err);
        });
    })
}