<template>
  <el-dialog title="流水线编辑" :visible.sync="formVisible" :close-on-click-modal="false">
    <el-col :span="24" style="margin-bottom: 20px;">
      <el-form :model="formData" label-width="180px" :rules="formRules" ref="formd">
        <el-form-item  label="触发器名称" prop="Name">
          <el-input v-model="formData.Name" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="formData.Desc" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="触发器类型">
          <el-select v-model="formData.Types" placeholder="请选择">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="流水线">
          <el-select v-model="formTriggerData.ModelId" placeholder="请选择">
            <el-option
                v-for="item in modelOptions"
                :key="item.Id"
                :label="item.Title"
                :value="item.Id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="表达式" v-if="formData.Types == '2' ">
          <el-input type="textarea"   v-model="formTriggerData.Expression" auto-complete="off" ></el-input>
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
        ]
      },
      //新增界面数据
      formData: {},
      formTriggerData: {},
      options: [{
        value: 0,
        label: '手动触发'
      }, {
        value: 1,
        label: 'git触发'
      }, {
        value: 2,
        label: '定时器触发'
      }],
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
        Name: '',
        Desc: '',
        Types: '',
        Config: '',
      }
      this.formTriggerData = {
        ModelId: '',
        Expression: '',
      }
      if (e){
        this.formData = {
          Id: e.Id,
          Name: e.Name,
          Desc: e.Desc,
          Types: e.Types,
          Config: e.Config,
        }
        var res = JSON.parse(e.Config);
        this.formTriggerData = {
          ModelId: res.ModelId,
          Expression: res.Expression,
        }
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
          this.formLoading = true;
          this.formData.Clrdir = this.formData.clrdir ? 1 : 2;
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

