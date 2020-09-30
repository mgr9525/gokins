<template>
	<el-dialog title="插件编辑" :visible.sync="formVisible" :close-on-click-modal="false">
        <el-col :span="24" style="margin-bottom: 20px;">
            <el-form :model="formData" label-width="80px" :rules="formRules" ref="formd">
				<el-form-item label="名称" prop="Title">
					<el-input v-model="formData.Title" auto-complete="off"></el-input>
				</el-form-item>
				<el-form-item label="类型" prop="Type">
					<!-- <el-input v-model="formData.Type" auto-complete="off"></el-input> -->
                    Shell
				</el-form-item>
				<el-form-item label="内容" prop="Cont">
					<el-input type="textarea" v-model="formData.Cont" auto-complete="off" :rows="20"></el-input>
				</el-form-item>
				<el-form-item label="排序" prop="Sort">
					<el-input v-model="formData.Sort" auto-complete="off"></el-input>
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
					],Type: [
						{ required: true, message: '请输入参数' }
					],Cont: [
						{ required: true, message: '请输入参数' }
					],Sort: [
						{ required: true, message: '请输入参数' }
					]
				},
				//新增界面数据
				formData: {}
			}
		},
		methods: {
            show(tid,e){
                this.formVisible=true;
                this.formData={
                    Id:'',
                    Tid:tid,
                    Type:1,
                    Title: '',
                    Para: '',
                    Cont: '',
                    Sort: '',
                }
                if(e)
                this.formData={
                    Id:e.Id,
                    Tid:e.Tid,
                    Type:e.Type,
                    Title: e.Title,
                    Para: e.Para,
                    Cont: e.Cont,
                    Sort: e.Sort,
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
						this.$post('/plug/edit',this.formData).then(res=>{
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

