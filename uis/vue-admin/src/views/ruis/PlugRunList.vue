<template>
	<section>
		<!--工具条-->
		<el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
					<el-button type="warning" @click="$router.back(-1)">返回</el-button>
					<el-button type="primary" @click="getList">刷新</el-button>
		</el-col>

		<div class="mains">
			<div style="width:400px;margin-right:10px">
				<el-card class="box-card runs" :shadow="mpdata[it.Id]&&mpdata[it.Id].selected?'always':'hover'"
				:class="mpdata[it.Id]&&mpdata[it.Id].selected?'runselect':''"
				v-for="(it,idx) in listdata" :key="'run'+it.Id">
					<div class="runrow" @click="showLog(idx)">
					<div style="flex:1">{{idx+1}}. {{it.Title}}
						<br/><span style="color:#909399">{{it.Hstm}}s</span>
					</div>

					<div>
					<el-tag v-if="it.RunStat==0" type="info">等待</el-tag>
					<el-tag v-if="it.RunStat==1" type="warning">运行</el-tag>
					<el-tag v-if="it.RunStat==2" type="danger">失败</el-tag>
					<el-tag v-if="it.RunStat==4" type="success">成功</el-tag>
					</div>
					</div>
				</el-card>
			</div>
			<div style="flex:1;white-space: break-spaces;word-break: break-all;">
				<el-card class="box-card">
				<div style="" v-text="logs[selid]&&logs[selid].text"></div>
				</el-card>
			</div>
		</div>


		<!--工具条-->
		<el-col :span="24" class="toolbar">
		</el-col>
	</section>
</template>

<script>
	//import NProgress from 'nprogress'

	export default {
		data() {
			return {
				tid:'',
				loading: false,
				listdata: [],

				selid:0,
				mpdata:{},
				logs:{}
			}
		},
		mounted() {
			this.tid=this.$route.query.id;
			if(this.tid==null||this.tid==''){
              	this.$router.push({ path: '/' });
				return
			}
			this.start();
			this.getList();
		},destroyed(){
			clearInterval(window.plugTimer);
		},
		methods: {
			//获取列表
			getList() {
				this.loading = true;
				//NProgress.start();
				this.$post('/plug/runs',{id:this.tid}).then((res) => {
              		console.log(res);
					this.loading = false;
					this.listdata = res.data;
				}).catch(err=>{
					this.loading = false;
					this.$message({
						message: err.response.data||'服务器错误',
						type: 'error'
					});
				});
			},selsChange(sels) {
				this.sels = sels;
			},start(){
				let that=this;
				let tmr=window.plugTimer;
				if(tmr)clearInterval(tmr);
				tmr=setInterval(() => {
					that.getList();
					that.getLog();
					/*for(let k in that.mpdata){
						let v=that.mpdata[k];
						if(v&&v.selected)getLog(k);
					}*/
				}, 1000);
				window.plugTimer=tmr;
			},showLog(idx){
				for(let i in this.listdata){
					let e=this.listdata[i];
					if(this.mpdata[e.Id])
						this.mpdata[e.Id].selected=false;
				}
				let e=this.listdata[idx];
				if(this.mpdata[e.Id]){
					this.mpdata[e.Id].selected=true;
				}else{
					this.mpdata[e.Id]={selected:true}
				}
				this.selid=e.Id;
				this.getLog();
				this.$forceUpdate();
				console.log('showLog:',this.mpdata[idx]);
			},getLog(){
				if(this.selid==''||this.selid<=0)return;
				let v=this.logs[this.selid]
				if(v&&v.up==false)return
				this.$post('/plug/log',{tid:this.tid,pid:this.selid}).then(res=>{
					this.logs[this.selid]=res.data;
				})
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
.mains{
	display:flex;
	clear:both;
}
.mains .runs{
	margin-bottom: 10px;
	cursor: pointer;
}
.mains .runselect{
	border: 1px solid red;
}
.mains .runrow{
	display: flex;
    width: 100%;

}
</style>