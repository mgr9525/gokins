<template>
	<el-dialog title="流水线编辑" :visible.sync="formVisible" :close-on-click-modal="false">
        <el-col :span="24" style="margin-bottom: 20px;">
            <el-form :model="formData" label-width="80px" :rules="formRules" ref="formd">
				<el-form-item label="项目名称" prop="Title">
					<el-input v-model="formData.Title" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="描述">
					<el-input type="textarea" v-model="formData.Desc" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="运行目录">
					<el-input v-model="formData.Wrkdir" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="目录问题">
					<el-switch v-model="formData.clrdir" active-text="创建或清空运行目录"></el-switch>
				</el-form-item>
				<el-form-item label="环境变量">
					<el-input type="textarea" v-model="formData.Envs" auto-complete="off" :rows="8" placeholder="每行一个变量,格式：name=xxxx"></el-input>
                    <el-tag type="info">注意：<span style="color:red">PATH</span>变量如果不填为gokins默认，如果填了请使用$PATH获取默认变量，例如：PATH=/home/bin:$PATH</el-tag>
				</el-form-item>
			</el-form>
		</el-col>
		<!--工具条-->
        <div slot="footer" class="dialog-footer">
            <el-button @click.native="formVisible = false">取消</el-button>
			<el-button type="primary" @click.native="formSubmit" :loading="formLoading">确认</el-button>
        </div>
    </el-dialog>
</template>


<script>
	export default {
		data() {
			return {
                formVisible:false,
				formLoading: false,
				formRules: {
					Title: [
						{ required: true, message: '请输入参数' }
					]
				},
				//新增界面数据
				formData: {}
			}
		},
		methods: {
            show(e){
                this.formVisible=true;
                this.formData={
                    Id:'',
                    Title: '',
                    Desc: '',
                    Envs:'',
                    Wrkdir:'',
                    clrdir:false
                }
                if(e)
                this.formData={
                    Id:e.Id,
                    Title: e.Title,
                    Desc: e.Desc,
                    Envs:e.Envs,
                    Wrkdir:e.Wrkdir,
                    clrdir:e.Clrdir==1
                }
            },/*handleSelect:function(id){
                this.tmpltCont='';
                if(id==''){
                    return
                }
                let it=this.tmplatdatas[id];
                if(it==null){
                    return
                }
                //console.log("123",this.tmplatdatas[id]);
                this.tmpltCont=it.Content;
                this.tmpltmapls=[];
                SmsTmpltKeys(id).then(res=>{
                    for(let i in res.data){
                        this.tmpltmapls.push({
                            key:res.data[i],value:''
                        });
                    }
                });
			},*/formSubmit(){
				this.$refs.formd.validate((valid) => {
					if (valid) {
						this.formLoading = true;
                        this.formData.Clrdir=this.formData.clrdir?1:2;
						this.$post('/model/edit',this.formData).then(res=>{
              				console.log(res);
                            this.$emit('submitOK');
                            this.formLoading = false;
                            this.formVisible = false;
                            //this.$message('操作成功');
						}).catch(err=>{
                            this.$emit('submitErr',err);
                            this.formLoading = false;
                            //this.formVisible = false;
                            this.$message({
                                message: err.response?err.response.data||'服务器错误':'网络错误',
                                type: 'error'
                            });
                        });
					}
				});
			}
		}
	};

</script>

<style scoped>
    .tmpdesc{
        margin-left: 10px;
        color:#d0d0d0;
    }
</style>

