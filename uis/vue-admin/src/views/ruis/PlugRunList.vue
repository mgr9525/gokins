<template>
	<section>
		<!--工具条-->
		<el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
					<el-button type="warning" @click="$router.back(-1)">返回</el-button>
					<el-button type="primary" @click="getList">刷新</el-button>
					<el-button type="danger" @click="handleStop" v-if="mrstat==0||mrstat==1">停止</el-button>
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

		<div class="mains">
			<div style="width:400px;margin-right:10px">
				<el-card class="box-card " style="background:#E0EEEE;margin-bottom: 10px;">
					<div class="runrow">
					<div style="flex:1"><span style="color:blue">任务运行情况</span>
						<br/><span style="color:red">{{mrerrs}}</span>
					</div>

					<div>
					<el-tag v-if="mrstat==-1" type="danger">停止</el-tag>
					<el-tag v-if="mrstat==0" type="info">等待</el-tag>
					<el-tag v-if="mrstat==1" type="warning">运行</el-tag>
					<el-tag v-if="mrstat==2" type="danger">失败</el-tag>
					<el-tag v-if="mrstat==4" type="success">成功</el-tag>
					</div>
					</div>
				</el-card>
				<el-card class="box-card runs" :shadow="mpdata[it.Id]&&mpdata[it.Id].selected?'always':'hover'"
				:class="mpdata[it.Id]&&mpdata[it.Id].selected?'runselect':''"
				v-for="(it,idx) in listdata" :key="'run'+it.Id" @click.native="showLog(idx)">
					<div class="runrow">
					<div style="flex:1">{{idx+1}}. {{it.Title}}
						<br/><span style="color:#909399">{{it.Hstm}}s</span>
					</div>

					<div>
					<el-tag v-if="it.RunStat==0&&mrstat<2" type="info">等待</el-tag>
					<el-tag v-if="it.RunStat==0&&mrstat>=2" type="info">未运行</el-tag>
					<el-tag v-if="it.RunStat==1" type="warning">运行</el-tag>
					<el-tag v-if="it.RunStat==2" type="danger">失败</el-tag>
					<el-tag v-if="it.RunStat==4" type="success">成功</el-tag>
					</div>
					</div>
				</el-card>
			</div>
			<div style="flex:1;white-space: break-spaces;word-break: break-all;">
				<el-card class="box-card">
					<div style="color:blue">{{logs[selid]&&logs[selid].tit}}</div>
				<div style="border-top:1px dashed #aaa">
					<pre style="white-space: pre-line;">{{logs[selid]&&logs[selid].text}}</pre>
				</div>
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
				running:false,
				loading: false,
				listdata: [],

				selid:0,
				mpdata:{},
				logs:{},

				md:{},
				mrstat:0,
				mrerrs:''
			}
		},
		mounted() {
			this.tid=this.$route.query.id;
			if(this.tid==null||this.tid==''){
              	this.$router.push({ path: '/' });
				return
			}
			
			this.running=true;
			this.loading = true;
			this.getList();
		},destroyed(){
			this.running=false;
		},
		methods: {
			getInfo(tid){
				if(this.md.Id||this.md.Id>0)return;
				this.$post('/model/get',{id:tid}).then(res=>{
					this.md=res.data;
				})
			},
			//获取列表
			getList() {
				if(!this.running)return;
				this.getLog(this.selid);
				this.$post('/plug/runs',{id:this.tid,first:this.loading}).then((res) => {
              		console.log(res);
					this.loading = false;
					this.getInfo(res.data.tid);
					this.listdata = res.data.list;
					this.mrstat=res.data.state;
					this.mrerrs=res.data.errs;
					if(res.data.end==true){
						this.running=false;
					}
					this.getList();
				}).catch(err=>{
					this.loading = false;
					this.$message({
						message: err.response.data||'服务器错误',
						type: 'error'
					});
					this.getList();
				});
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
					this.mpdata[e.Id]={tit:e.Title,selected:true}
				}
				this.selid=e.Id;
				this.$forceUpdate();
				console.log('showLog:',this.mpdata[idx]);
				if(!this.running)this.getLog(this.selid);
			},getLog(selid){
				if(selid==''||selid<=0)return;
				let log=this.logs[selid];
				if(log&&!this.running)return;
				this.$post('/plug/log',{tid:this.tid,pid:selid,pos:log?log.pos:0}).then(res=>{
					res.data.tit=this.mpdata[selid].tit;
					if(log&&res.data.pos>0){
						log.pos=res.data.pos;
						log.text+=res.data.text;
					}else
						this.logs[selid]=res.data;
					this.$forceUpdate();
				})
			},handleStop(){
				this.$post('/model/stop',{id:this.tid}).then(res=>{
					this.$message('操作成功');
				}).catch(err=>{
					this.$message({
						message: err.response?err.response.data||'服务器错误':'网络错误',
						type: 'error'
					});
				});
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