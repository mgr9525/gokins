import Vue from 'vue'
import Vuex from 'vuex'
import * as actions from './actions'
import * as getters from './getters'
import * as mutations from './mutations'

Vue.use(Vuex)

// 应用初始状态
const state = {
    count: 10,
    userinfo:null
}
// 创建 store 实例
export default new Vuex.Store({
    actions,
    getters,
    state,
    mutations
})