import{e,p as a,u as s,M as l,q as t,g as r}from"./element-plus.2c9426ce.js";import{_ as o,u as i,D as n,S as d}from"./index.d5c032bc.js";import{r as u,i as c}from"./@element-plus.e895b9e1.js";import{J as p,r as m,o as f,c as b,a as v,O as y,R as g,u as j,W as _,V as h,av as k,aw as x}from"./@vue.2c474d7f.js";import"./lodash-es.98e49362.js";import"./async-validator.fb49d0f5.js";import"./@vueuse.d618444d.js";import"./dayjs.2790329b.js";import"./axios.765908e4.js";import"./@ctrl.b082b0c1.js";import"./@popperjs.36402333.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.8aeb3683.js";import"./vue-router.3b7aa6c5.js";import"./pinia.39330980.js";import"./vue-demi.b3a9cad9.js";import"./screenfull.7cf96174.js";import"./vue-echarts.ebbe6a42.js";import"./resize-detector.4e96b72b.js";import"./echarts.2b8dcc79.js";import"./zrender.54f191d6.js";import"./tslib.60310f1a.js";import"./nprogress.aae2b787.js";const F=e=>(k("data-v-57001553"),e=e(),x(),e),P={class:"main-container"},S={class:"main",style:{"background-color":"inherit"}},w=F((()=>v("div",{class:"card-header"},[v("span",null,"在线升级")],-1))),V={class:"item"},C=F((()=>v("br",null,null,-1))),z=F((()=>v("br",null,null,-1))),D=F((()=>v("br",null,null,-1))),U=F((()=>v("br",null,null,-1))),q={style:{"text-align":"center"}},R=F((()=>v("div",{class:"card-header"},[v("span",null,"本地升级")],-1))),I={class:"item"},J=F((()=>v("div",null,"请选择您要升级的系统文件。",-1))),M=F((()=>v("br",null,null,-1))),O=F((()=>v("br",null,null,-1))),T=F((()=>v("br",null,null,-1))),W={style:{"text-align":"center"}},A=F((()=>v("div",{class:"el-upload__tip"},"只能上传一个文件，只支持zip格式文件！",-1))),B={class:"dialog-footer"};var E=o({__name:"SysUpgrade",setup(o){const k=i(),x=p({uFlag:!1,sysParams:{}});(()=>{const e={token:k.token,data:{}};n.getSysParams(e).then((e=>{console.log("getSysParams -> res",e),"0"===e.code&&(x.sysParams=e.data,console.log("getSysParams -> ctxData.sysParams",x.sysParams))}))})();const F=()=>{G.value&&G.value.submit()},E=e=>{},G=m(null),H=e=>{G.value.clearFiles();const a=e[0];G.value.handleStart(a)},K=e=>{console.log("beforeUpload -> file",e);const a=["application/x-tar"].includes(e.type);if(a){if(a){let a=new FormData;a.append("file",e);const s={token:k.token,contentType:"multipart/form-data",data:a};d.updateSystem(s).then((e=>{"0"===e.code?(r.success(e.message),x.uFlag=!1):Q(e)}))}}else r({type:"error",message:"文件格式不正确,必须是tar文件！"})},L=()=>{x.uFlag=!1,G.value.clearFiles()},N=()=>{L()},Q=e=>{r({type:"error",message:e.message})};return(o,i)=>{const n=e,d=a,p=s,m=l,k=t;return f(),b("div",P,[v("div",S,[y(p,{class:"box-card",shadow:"hover"},{header:g((()=>[w])),default:g((()=>[v("div",V,[v("div",null,"当前系统版本："+_(j(x).sysParams.softVer)+" 。",1),C,v("div",null,"最新系统版本："+_(j(x).sysParams.softVer)+" 。",1),z,D,U,v("div",q,[y(d,{type:"primary",onClick:i[0]||(i[0]=e=>{r.info("功能升级中...")})},{default:g((()=>[y(n,{class:"el-input__icon"},{default:g((()=>[y(j(u))])),_:1}),h(" 在线升级系统到最新版本 ")])),_:1})])])])),_:1}),y(p,{class:"box-card",shadow:"hover"},{header:g((()=>[R])),default:g((()=>[v("div",I,[J,M,O,T,v("div",W,[y(d,{type:"primary",onClick:i[1]||(i[1]=e=>(console.log("importSysConfig"),void(x.uFlag=!0)))},{default:g((()=>[y(n,{class:"el-input__icon"},{default:g((()=>[y(j(c))])),_:1}),h(" 导入系统升级文件 ")])),_:1})])])])),_:1})]),y(k,{modelValue:j(x).uFlag,"onUpdate:modelValue":i[2]||(i[2]=e=>j(x).uFlag=e),title:"上传系统升级文件",width:"600px","before-close":N,"close-on-click-modal":!1},{footer:g((()=>[v("span",B,[y(d,{onClick:L},{default:g((()=>[h("取消")])),_:1}),y(d,{type:"primary",onClick:F},{default:g((()=>[h("上传")])),_:1})])])),default:g((()=>[y(m,{ref_key:"uploadRef",ref:G,action:"","auto-upload":!1,"http-request":E,limit:1,"on-exceed":H,"before-upload":K},{tip:g((()=>[A])),default:g((()=>[y(d,{type:"primary"},{default:g((()=>[h("选择文件")])),_:1})])),_:1},512)])),_:1},8,["modelValue"])])}}},[["__scopeId","data-v-57001553"]]);export{E as default};
