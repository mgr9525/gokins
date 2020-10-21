package mgr

var hookjsMap map[string]string

func init() {
	hookjsMap = make(map[string]string)
	hookjsMap["gitee"] = `

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
		ret.errs='请求密码错误';
		return ret;
	}

	console.log('start run main function!!!!');
	console.log('hook_name:',getBody().hook_name);
	console.log('head_commit.id:',getBody().head_commit.id);
	ret.check=true;
	return ret
}
`
}
