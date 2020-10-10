<template>
	<el-dialog title="修改密码" :visible.sync="formVisible" :close-on-click-modal="false">
        <el-col :span="24" style="margin-bottom: 20px;">
            <el-form :model="formData" label-width="80px" :rules="formRules" ref="formd">
				<el-form-item label="旧密码" prop="pass">
					<el-input type="password" v-model="formData.pass" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="新密码" prop="newpass">
					<el-input type="password" v-model="formData.newpass" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="重复密码" prop="repass">
					<el-input type="password" v-model="formData.repass" auto-complete="off"></el-input>
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
					pass: [
						{ required: true, message: '请输入参数' }
					],newpass: [
						{ required: true, message: '请输入参数' }
					],repass: [
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
                    pass:'',
                    newpass: '',
                    repass: '',
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
                        if(this.formData.newpass!=this.formData.repass){
                            this.$message({
                                message: '新密码不一致',
                                type: 'error'
                            });
                            return
                        }
						this.formLoading = true;
						this.$post('/lg/uppass',this.formData).then(res=>{
              				console.log(res);
                            this.$emit('submitOK');
                            this.formLoading = false;
                            this.formVisible = false;
                            //this.$message('操作成功');
                            this.$message({
                                message: '修改密码成功!!!',
                                type: 'success'
                            });
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

