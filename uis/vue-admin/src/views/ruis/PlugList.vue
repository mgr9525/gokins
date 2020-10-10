<template>
	<section>
		<!--工具条-->
		<el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
			<el-form :inline="true">
				<el-form-item>
					<el-button type="warning" @click="$router.back(-1)">返回</el-button>
					<el-button type="primary" v-on:click="getList">刷新</el-button>
					<el-button type="primary" @click="$refs.editor.show(tid)">新增</el-button>
				</el-form-item>
			</el-form>
		</el-col>

		<!--列表-->
		<el-table :data="listdata" highlight-current-row v-loading="loading" @selection-change="selsChange" style="width: 100%;">
			<el-table-column type="index" width="60">
			</el-table-column>
			<el-table-column label="名称" sortable>
				<template slot-scope="{row}">
					<span>{{ row.Title }}</span>
				</template>
			</el-table-column>
			<el-table-column label="类型" width="80" sortable>
				<template slot-scope="{row}">
					<span>Shell</span>
				</template>
			</el-table-column>
			<el-table-column prop="Sort" label="排序" width="100" sortable>
			</el-table-column>
			<el-table-column prop="Times" label="创建时间" width="200" sortable>
			</el-table-column>
			<el-table-column label="操作" width="150">
				<template slot-scope="{row}">
					<el-button-group>
					<el-button size="small" type="warning" @click="$refs.editor.show(tid,row)">编辑</el-button>
              <el-popconfirm title="确定要删除吗？" @onConfirm="handleDel(row)">
					<el-button size="small" type="danger" slot="reference">删除</el-button>
              </el-popconfirm>
					</el-button-group>
				</template>
			</el-table-column>
		</el-table>

		<!--工具条-->
		<!-- <el-col :span="24" class="toolbar">
			<el-button type="danger" @click="batchRemove" :disabled="this.sels.length===0">批量删除</el-button>
		</el-col> -->
		<PlugForm ref="editor" @submitOK="getList()"/>
	</section>
</template>

<script>
import PlugForm from './PlugForm'
	//import NProgress from 'nprogress'

	export default {
		components:{PlugForm},
		data() {
			return {
				tid:'',
				loading: false,
				total:0,
				limit:0,
				listdata: [],
				sels: [],//列表选中列
			}
		},
		mounted() {
			this.tid=this.$route.query.id;
			if(this.tid==null||this.tid==''){
              this.$router.push({ path: '/' });
				return
			}
			this.getList();
		},
		methods: {
			//获取列表
			getList() {
				this.loading = true;
				//NProgress.start();
				this.$post('/plug/list',{tid:this.tid}).then((res) => {
              		console.log(res);
					this.loading = false;
					this.listdata = res.data;
					//NProgress.done();
				}).catch(err=>{
					this.loading = false;
					this.$message({
						message: err.response.data||'服务器错误',
						type: 'error'
					});
				});
			},selsChange(sels) {
				this.sels = sels;
			},handleAdd(){

			},handleEdit(){

			},handleDel(et){
				this.$post('/plug/del',{id:et.Id}).then(res=>{
					//this.$message('操作成功');
					this.getList();
				}).catch(err=>{
					this.$message({
						message: err.response?err.response.data||'服务器错误':'网络错误',
						type: 'error'
					});
				});
			},batchRemove(){

			}
		}
	}

</script>

<style scoped>
.wxmpTit{
	line-height: 60px;
	margin-top: 5px;
	margin-bottom: 5px;
}
.wxmpTit img{
	width: 60px;
	height: 60px;
	float: left;
	margin-right: 10px;
}
</style>