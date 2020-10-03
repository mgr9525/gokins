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
					<el-button type="primary" @click="handleAdd">新增</el-button>
				</el-form-item>
			</el-form>
		</el-col>

		<!--列表-->
		<el-table :data="listdata" highlight-current-row v-loading="loading" @selection-change="selsChange" style="width: 100%;">
			<el-table-column type="selection" width="55">
			</el-table-column>
			<el-table-column type="index" width="60">
			</el-table-column>
			<el-table-column label="名称" sortable>
				<template slot-scope="{row}">
					<div class="wxmpTit">
					<img :src="row.Avat"/>
					<span>{{ row.Name }}</span>
					</div>
				</template>
			</el-table-column>
			<el-table-column prop="Appid" label="Appid" width="180" sortable>
			</el-table-column>
			<el-table-column label="状态" width="80">
				<template slot-scope="{row}">
					<span v-if="row.Cancel" style="color:red">取消授权</span>
					<span v-if="!row.Cancel" style="color:green">正常</span>
				</template>
			</el-table-column>
			<el-table-column label="操作" width="150">
				<template scope="scope">
					<el-button size="small" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
					<el-button type="danger" size="small" @click="handleDel(scope.$index, scope.row)">删除</el-button>
				</template>
			</el-table-column>
		</el-table>

		<!--工具条-->
		<el-col :span="24" class="toolbar">
			<el-button type="danger" @click="batchRemove" :disabled="this.sels.length===0">批量删除</el-button>
			<el-pagination layout="prev, pager, next" :current-page.sync="filters.page" :total="total" :page-size="limit" @current-change="getList" style="float:right;">
			</el-pagination>
		</el-col>
	</section>
</template>

<script>
	//import NProgress from 'nprogress'

	export default {
		data() {
			return {
				filters:{
					page: 1,
					s:'',
					q:''
				},
				loading: false,
				total:0,
				limit:0,
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
				this.$post('/api/auths',this.filters).then((res) => {
              		console.log(res);
					this.loading = false;
					this.listdata = res.data.Data;
					this.total = res.data.Total;
					this.limit = res.data.Size;
					this.filters.page=res.data.Page;
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

			},handleDel(){

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