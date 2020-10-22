<template>
  <section>
    <!--工具条-->
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.q" placeholder="搜索"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="getList">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="$refs.editor.show()">新增</el-button>
        </el-form-item>
      </el-form>
    </el-col>

    <!--列表-->
    <el-table :data="listdata" highlight-current-row v-loading="loading" @selection-change="selsChange"
              style="width: 100%;">
      <!-- <el-table-column type="selection" width="55">
      </el-table-column> -->
      <el-table-column type="index" width="60">
      </el-table-column>
      <el-table-column label="名称" prop="Title" width="250" sortable>
      </el-table-column>
      <el-table-column label="描述">
				<template slot-scope="{row}">
          <span>{{row.Desc}}</span>
          <div><el-tag type="danger" v-if="row.Errs!=''">{{row.Errs}}</el-tag></div>
          <div><el-tag type="info" v-if="row.Types == 'hook'">hook地址：/hook/trigger/{{row.Id}}</el-tag></div>
				</template>
      </el-table-column>
      <el-table-column prop="Types" label="触发器类型" width="150" :formatter="typesFormatter" sortable>
      </el-table-column>
      <el-table-column prop="Times" label="创建时间" width="200" :formatter="dateFormat" sortable>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template slot-scope="{row}">
          <el-button-group>
            <el-button size="small" type="warning" @click="$refs.editor.show(row)">编辑</el-button>
            <!-- <el-button size="small" @click="$router.push({path:'/Triggers/info?id='+row.Id})">插件</el-button> -->
            <el-popconfirm title="确定要删除吗？" @onConfirm="handleDel(row)">
              <el-button size="small" type="danger" slot="reference">删除</el-button>
            </el-popconfirm>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <!--工具条-->
    <el-col :span="24" class="toolbar">
      <!-- <el-button type="danger" @click="batchRemove" :disabled="this.sels.length===0">批量删除</el-button> -->
      <el-pagination layout="prev, pager, next" :current-page.sync="filters.page" :total="total" :page-size="limit"
                     @current-change="getList" style="float:right;">
      </el-pagination>
    </el-col>
    <TriggerForm ref="editor" @submitOK="getList()"/>
  </section>
</template>

<script>
import TriggerForm from './TriggerForm'
//import NProgress from 'nprogress'

export default {
  components: {TriggerForm},
  data() {
    return {
      filters: {
        page: 1,
        s: '',
        q: ''
      },
      loading: false,
      total: 0,
      limit: 0,
      listdata: [],
      sels: [],//列表选中列
    }
  },
  mounted() {
    this.getList();
  },
  methods: {
    //获取列表
    getList() {
      this.loading = true;
      //NProgress.start();
      this.$post('/trigger/list', this.filters).then((res) => {
        console.log(res);
        this.loading = false;
        this.listdata = res.data.Data;
        this.total = res.data.Total;
        this.limit = res.data.Size;
        this.filters.page = res.data.Page;
        //NProgress.done();
      }).catch(err => {
        this.loading = false;
        this.$message({
          message: err.response.data || '服务器错误',
          type: 'error'
        });
      });
    }, selsChange(sels) {
      this.sels = sels;
    }, handleAdd() {

    }, handleEdit() {

    }, handleDel(et) {
      this.$post('/trigger/del', {id: et.Id}).then(res => {
        //this.$message('操作成功');
        this.getList();
      }).catch(err => {
        this.$message({
          message: err.response ? err.response.data || '服务器错误' : '网络错误',
          type: 'error'
        });
      });
    }, batchRemove() {

    },
    typesFormatter: function (row, column) {
      let typ=row.Types == 'git' ? "git" : row.Types == 'timer' ? "定时器" : "手动";
      if(row.Enable==1){
        typ+='(已激活)';
      }else{
        typ+='(未激活)';
      }
      return typ;
    },
    dateFormat: function (row, column) {
      var t = new Date(row.Times);//row 表示一行数据, updateTime 表示要格式化的字段名称
      return t.getFullYear() + "-" + (t.getMonth() + 1) + "-" + t.getDate() + " " + t.getHours() + ":" + t.getMinutes() + ":" + t.getSeconds() + "." + t.getMilliseconds();
    },
  }
}

</script>

<style scoped>
.wxmpTit {
  line-height: 60px;
  margin-top: 5px;
  margin-bottom: 5px;
}

.wxmpTit img {
  width: 60px;
  height: 60px;
  float: left;
  margin-right: 10px;
}
</style>