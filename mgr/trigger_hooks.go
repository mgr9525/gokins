package mgr

type Hookjs struct {
	Uis map[string]string
	js  string
}

var HookjsMap map[string]*Hookjs

func init() {
	HookjsMap = make(map[string]*Hookjs)
	HookjsMap["gitee"] = &Hookjs{
		Uis: map[string]string{"password": "string", "operate": "array"},
		js: `

function main(){
	console.log('start run main function!!!!');
	var ret={check:false};
	var conf=getConf();
	var body=getBody();
	if(conf.password==''){
		ret.errs='不支持非password的方式';
		return ret;
	}
	if(conf.password!=body.password){
		ret.errs='触发请求密码错误';
		return ret;
    }
	ret.check=true;
    if(conf.operate&&conf.operate.length>0){
        ret.check=false;
        for(var i in conf.operate){
            var it=conf.operate[i];
            console.log('operate['+i+']',it);
            if(it=='push'){
                if(body.hook_name=='push_hooks'){
                    ret.check=true;
                    break;
                }
            }else if(it=='merged'){
                if(body.hook_name=='merge_request_hooks'&&body.pull_request&&body.pull_request.merged==true){
                    ret.check=true;
                    break;
                }
            }
        }
    }
	return ret
}
`,
	}
}
