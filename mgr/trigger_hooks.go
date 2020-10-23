package mgr

type Hookjs struct {
	Uis  map[string]string
	Desc string
	Defs string
	js   string
}

var HookjsMap map[string]*Hookjs

func init() {
	HookjsMap = make(map[string]*Hookjs)
	HookjsMap["web"] = &Hookjs{
		Uis:  map[string]string{"password": "string"},
		Desc: "password:触发密码",
		Defs: `{"password":"pwd"}`,
		js: `

function main(){
	console.log('start run main function!!!!');
	var ret={check:false};
	var conf=getConf();
    var body=getBody();
	if(conf.password!=body.password){
		ret.errs='触发请求密码错误';
		return ret;
    }
	ret.check=true;
	return ret
}
`,
	}
	HookjsMap["gitee"] = &Hookjs{
		Uis:  map[string]string{"password": "string", "branch": "string"},
		Desc: "password:推送密码,branch:push对象分支",
		Defs: `{"password":"pwd","branch":"master"}`,
		js: `

function main(){
	console.log('start run main function!!!!');
	var ret={check:false};
	var conf=getConf();
    var body=getBody();
	if(conf.password!=body.password){
		ret.errs='触发请求密码错误';
		return ret;
    }
    if(body.hook_name!='push_hooks'||!body.ref||body.ref==''){
        return ret;
    }
    console.log(conf.branch,body.ref);
    if(conf.branch&&conf.branch!=''&&body.ref!='refs/heads/'+conf.branch){
        return ret;
    }
	ret.check=true;
	return ret
}
`,
	}
	HookjsMap["github"] = &Hookjs{
		Uis:  map[string]string{"secretkey": "string", "branch": "string"},
		Desc: "secretkey:签名秘钥,branch:push对象分支",
		Defs: `{"secretkey":"pwd","branch":"master"}`,
		js: `

	function main(){
		console.log('start run main function!!!!');
		var ret={check:false};
		var conf=getConf();
	    var body=getBody();
	    var bodys=getBodys();
		var nm=getHeader('X-GitHub-Event');
	    var tk=getHeader('X-Hub-Signature');
		if(verifySignature(tk,conf.password,bodys)){
			ret.errs='触发请求秘钥错误';
			return ret;
	    }
	    if(nm!='push'||!body.ref||body.ref==''){
	        return ret;
	    }
	    console.log(conf.branch,body.ref);
	    if(conf.branch&&conf.branch!=''&&body.ref!='refs/heads/'+conf.branch){
	        return ret;
	    }
		ret.check=true;
		return ret
	}
	`,
	}
	HookjsMap["gitlab"] = &Hookjs{
		Uis:  map[string]string{"token": "string", "branch": "string"},
		Desc: "token:秘钥,branch:push对象分支",
		Defs: `{"token":"pwd","branch":"master"}`,
		js: `

function main(){
	console.log('start run main function!!!!');
	var ret={check:false};
	var conf=getConf();
    var body=getBody();
	var tk=getHeader('X-Gitlab-Token');
	if(conf.token!=tk){
		ret.errs='触发请求秘钥错误';
		return ret;
    }
    if(body.object_kind!='push'||!body.ref||body.ref==''){
        return ret;
    }
    console.log(conf.branch,body.ref);
    if(conf.branch&&conf.branch!=''&&body.ref!='refs/heads/'+conf.branch){
        return ret;
    }
	ret.check=true;
	return ret
}
`,
	}
}
