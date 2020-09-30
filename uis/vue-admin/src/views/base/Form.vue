<template>
	<el-dialog title="上传文件" :visible.sync="formVisible" :close-on-click-modal="false">
        <el-col :span="24" style="margin-bottom: 20px;">
            <el-form :model="formData" label-width="80px" :rules="formRules" ref="formd">
				<el-form-item label="任务名称" prop="title">
					<el-input v-model="formData.title" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="短信名称" prop="name">
					<el-input v-model="formData.name" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="描述" prop="name">
					<el-input type="textarea" v-model="formData.name" auto-complete="off"></el-input>
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
					tmpltid: [
						{ required: true, message: '请选择模版' }
					],title: [
						{ required: true, message: '请输入任务名称' }
					],name: [
						{ required: true, message: '请输入短信名称' }
					],phone: [
						{ required: true, message: '请输入参数' }
					],level: [
						{ required: true, message: '请输入参数' }
					]
				},
				//新增界面数据
				formData: {}
			}
		},
		methods: {
            show(){
                this.formVisible=true;
                this.formData={
                    tmpltid:'',
                    name: '',
                    phone: '',
                    level: '1',
                    content: '',
					enable:true,
					send_time:''
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
                        var regd=/[^\d\,]/;
                        if(regd.test(this.formData.phone)){
                            this.$message({
                                message: '手机号码列表错误，只能是数字和逗号(,)',
                                type: 'error'
                            });
                            return;
                        }
                        var params=this.formData;
                        params['params']=conts;
                        params['phones']=this.formData.phone.split(',');

						this.formLoading = true;
						this.$post('/api/',params).then(res=>{
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

