<template>
	<section>
		<!--工具条-->
		<el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
			<el-form :inline="true">
				<el-form-item>
					<el-button type="warning" @click="$router.back(-1)">返回</el-button>
					<el-button type="success" @click="handleRun">运行</el-button>
					<el-button type="primary" v-on:click="getList">刷新</el-button>
					<el-button type="info" v-on:click="$refs.editor.show(md)">编辑</el-button>
					<el-button type="primary" v-on:click="$router.push('/models/info?id='+md.Id)">插件列表</el-button>
				</el-form-item>
			</el-form>
		</el-col>
		<el-card class="box-card" style="margin-bottom:20px">
		<el-row class="text item infoItem">
			<el-col :span="10">任务名称：{{md.Title}}</el-col>
			<el-col :span="6">创建时间：{{md.Times}}</el-col>
		</el-row>
		<el-row class="text item infoItem">
			<el-col :span="12">任务描述：{{md.Desc}}</el-col>
		</el-row>
		<el-row class="text item infoItem">
			<el-col :span="10">运行目录：{{md.Wrkdir}}</el-col>
			<el-col :span="6">创建或清空运行目录：{{md.Clrdir==1?'是':'否'}}</el-col>
		</el-row>
		<el-row class="text item infoItem">
			<el-col :span="12">环境变量：<p v-text="md.Envs"></p></el-col>
		</el-row>
		</el-card>

		<!--列表-->
		<el-table :data="listdata" highlight-current-row v-loading="loading" @selection-change="selsChange" style="width: 100%;">
			<el-table-column label="编号" width="80">
				<template slot-scope="{row}">
					<el-link type="primary" @click="$router.push('/models/plug/runs?id='+row.Id)">
					#{{row.Id}}</el-link>
				</template>
			</el-table-column>
			<el-table-column label="触发" width="120">
				<template slot-scope="{row}">
					<el-tag>{{row.Tgtyps}}</el-tag>
				</template>
			</el-table-column>
			<el-table-column prop="Times1" label="运行时间" width="200" sortable>
			</el-table-column>
			<el-table-column prop="Times2" label="结束时间" width="200" sortable>
			</el-table-column>
			<el-table-column prop="Nick" label="执行人" width="200" sortable>
			</el-table-column>
			<el-table-column label="状态" width="80">
				<template slot-scope="{row}">
					<span v-if="row.State==-1" style="color:red">已停止</span>
					<span v-if="row.State==0" style="color:red">等待中</span>
					<span v-if="row.State==1" style="color:blue">运行中</span>
					<span v-if="row.State==2" style="color:red">运行失败</span>
					<span v-if="row.State==4" style="color:green">运行成功</span>
				</template>
			</el-table-column>
			<el-table-column prop="Errs" label="错误" sortable>
			</el-table-column>
			<el-table-column label="操作" width="100">
				<template scope="{row}">
					<el-button type="danger" size="small" @click="handleStop(row.Id)" v-if="row.State==0||row.State==1">停止</el-button>
				</template>
			</el-table-column>
		</el-table>

		<!--工具条-->
		<el-col :span="24" class="toolbar">
			<!-- <el-button type="danger" @click="batchRemove" :disabled="this.sels.length===0">批量删除</el-button> -->
			<el-pagination layout="prev, pager, next" :current-page.sync="page" :total="total" :page-size="limit" @current-change="getList" style="float:right;">
			</el-pagination>
		</el-col>
		<ModelForm ref="editor" @submitOK="getInfo()"/>
	</section>
</template>

<script>
import ModelForm from './ModelForm'
	//import NProgress from 'nprogress'

	export default {
		components:{ModelForm},
		data() {
			return {
				tid:'',
					page: 1,
				loading: false,
				total:0,
				limit:0,
				listdata: [],
				sels: [],//列表选中列

				md:{}
			}
		},
		mounted() {
			this.tid=this.$route.query.id;
			if(this.tid==null||this.tid==''){
              	this.$router.push({ path: '/' });
				return
			}
			this.getInfo();
			this.getList();
		},
		methods: {
			getInfo(){
				this.$post('/model/get',{id:this.tid}).then(res=>{
					this.md=res.data;
				})
			},
			//获取列表
			getList() {
				this.loading = true;
				//NProgress.start();
				this.$post('/model/runs',{tid:this.tid,page:this.page}).then((res) => {
              		console.log(res);
					this.loading = false;
					this.listdata = res.data.Data;
					this.total = res.data.Total;
					this.limit = res.data.Size;
					this.page=res.data.Page;
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
			},handleRun(){
				this.$post('/model/run',{id:this.tid}).then(res=>{
					//this.$message('操作成功');
					this.getList();
				}).catch(err=>{
					this.$message({
						message: err.response?err.response.data||'服务器错误':'网络错误',
						type: 'error'
					});
				});
			},handleEdit(){

			},handleStop(id){
				this.$post('/model/stop',{id:id}).then(res=>{
					this.$message('操作成功');
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
.infoItem{
	margin-bottom: 10px;
}
</style>