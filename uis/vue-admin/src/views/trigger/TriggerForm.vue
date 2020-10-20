<template>
  <el-dialog title="触发器编辑" :visible.sync="formVisible" :close-on-click-modal="false">
    <el-col :span="24" style="margin-bottom: 20px;">
      <el-form :model="formData" label-width="180px" :rules="formRules" ref="formd">
        <el-form-item  label="触发器名称" prop="Title">
          <el-input v-model="formData.Title" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="formData.Desc" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="触发器类型" prop="Types">
          <el-select v-model="formData.Types" placeholder="请选择">
            <!-- <el-option label="git触发" value="git"></el-option> -->
            <el-option label="定时器触发" value="timer"></el-option>
          </el-select>
					<el-switch v-model="formData.enable" active-text="激活"></el-switch>
        </el-form-item>
        <el-form-item label="流水线">
          <el-select v-model="formTriggerData.mid" placeholder="请选择">
            <el-option
                v-for="item in modelOptions"
                :key="item.Id"
                :label="item.Title"
                :value="item.Id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="重复" v-if="formData.Types == 'timer'">
          <el-select v-model="formTriggerData.repeated" placeholder="请选择">
            <el-option label="不重复" value="0"></el-option>
            <el-option label="每天" value="1"></el-option>
            <el-option label="每周" value="2"></el-option>
            <el-option label="每月" value="3"></el-option>
            <el-option label="每年" value="4"></el-option>
          </el-select>
          <el-date-picker
            v-model="formTriggerData.dates"
            type="datetime"
            placeholder="选择日期时间">
          </el-date-picker>
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
      cronPopover:false,
      cron:'',
      filters: {
        page: 1,
        size: 9999,
      },
      formVisible: false,
      formLoading: false,
      formRules: {
        Title: [
          {required: true, message: '请输入参数'}
        ],Types: [
          {required: true, message: '请输入参数'}
        ]
      },
      //新增界面数据
      formData: {},
      formTriggerData: {},
      modelOptions: [],
    }
  },
  mounted() {
    this.getList();
  },
  methods: {
    changeCron(val){
      this.cron=val
    },
    getList() {
      this.loading = true;
      //NProgress.start();
      this.$post('/model/list', this.filters).then((res) => {

        this.loading = false;
        this.modelOptions = res.data.Data;
        console.log(this.modelOptions);
        //NProgress.done();
      }).catch(err => {
        this.loading = false;
        this.$message({
          message: err.response.data || '服务器错误',
          type: 'error'
        });
      });
    },
    show(e) {
      this.formVisible = true;
      this.formData = {
        Id: '',
        Title: '',
        Desc: '',
        Types: '',
        Config: '',
        enable:false
      }
      this.formTriggerData = {
        mid:'',
        repeated:'',
        dates:''
      }
      if (e){
        this.formData = {
          Id: e.Id,
          Title: e.Title,
          Desc: e.Desc,
          Types: e.Types,
          Config: e.Config,
          enable:e.Enable==1
        }
        try{
        var res = JSON.parse(e.Config);
        this.formTriggerData = res
        }catch(e){}
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
			},*/formSubmit() {
      this.$refs.formd.validate((valid) => {
        if (valid) {
          if(this.formData.Types=='timer'){
            console.log('formTriggerData:',this.formTriggerData);
            if(this.formTriggerData.repeated==''){
              this.$message('请选择重复类型');
              return
            }
            if(this.formTriggerData.dates==''){
              this.$message('请选择日期');
              return
            }
          }
          this.formLoading = true;
          this.formData.Enable = this.formData.enable ? 1 : 2;
          this.formData.Config = JSON.stringify(this.formTriggerData)
          this.$post('/trigger/edit', this.formData).then(res => {
            console.log(res);
            this.$emit('submitOK');
            this.formLoading = false;
            this.formVisible = false;
            //this.$message('操作成功');
          }).catch(err => {
            this.$emit('submitErr', err);
            this.formLoading = false;
            //this.formVisible = false;
            this.$message({
              message: err.response ? err.response.data || '服务器错误' : '网络错误',
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
.tmpdesc {
  margin-left: 10px;
  color: #d0d0d0;
}
</style>

