import Login from './views/Login.vue'
import NotFound from './views/404.vue'
import Home from './views/Home.vue'

let routes = [
    {
        path: '/login',
        component: Login,
        name: '',
        hidden: true
    },
    {
        path: '/404',
        component: NotFound,
        name: '',
        hidden: true
    },
    //{ path: '/main', component: Main },
    {
        path: '/',
        component: Home,
        name: '工作区',
        iconCls: 'el-icon-message',//图标样式class
        children: [
            { path: '/models', component: require('@/views/ruis/ModelList'), name: '流水线' },
            { path: '/trigger', component: require('@/views/trigger/TriggerList'), name: '触发器' },
            { path: '/models/info', component: require('@/views/ruis/PlugList'), name: '流水线插件', hidden: true },
            { path: '/models/runs', component: require('@/views/ruis/RunList'), name: '流水线运行', hidden: true },
            { path: '/models/plug/runs', component: require('@/views/ruis/PlugRunList'), name: '流水日志', hidden: true },
        ]
    },
    /*{
        path: '/',
        component: Home,
        name: '导航二',
        iconCls: 'fa fa-id-card-o',
        children: [
            { path: '/page4', component: Page4, name: '页面4' },
            { path: '/page5', component: Page5, name: '页面5' }
        ]
    },
    {
        path: '/',
        component: Home,
        name: '',
        iconCls: 'fa fa-address-card',
        leaf: true,//只有一个节点
        children: [
            { path: '/page6', component: Page6, name: '导航三' }
        ]
    },
    {
        path: '/',
        component: Home,
        name: 'Charts',
        iconCls: 'fa fa-bar-chart',
        children: [
            { path: '/echarts', component: echarts, name: 'echarts' }
        ]
    },*/
    {
        path: '*',
        hidden: true,
        redirect: { path: '/404' }
    }
];

export default routes;