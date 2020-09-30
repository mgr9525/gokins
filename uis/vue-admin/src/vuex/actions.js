import api from '@/api/apis'
//test
export const increment = ({commit}) => {
    commit('INCREMENT')
}
export const decrement = ({commit}) => {
    commit('DECREMENT')
}


export const getLgInfo = ({commit}) => {
    return new Promise((resolve, reject)=>{
        api.post('/lg/info').then(res=>{
            commit('setLgInfo',res.data);
            resolve();
        }).catch(err=>{
            reject(err);
        });
    })
}