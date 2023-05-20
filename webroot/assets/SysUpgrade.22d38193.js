import{e as a,p as e,u as s,H as l,q as t,h as o}from"./element-plus.795748f2.js";import{_ as r,u as i,S as n}from"./index.3a6602e1.js";import{D as d}from"./dashboard.9cfd7434.js";import{N as u,M as c}from"./@element-plus.102632fb.js";import{Y as p,f as m,o as f,c as y,a as v,S as b,Q as g,U as _,u as h,as as j,at as k,a0 as x}from"./@vue.8518410d.js";import"./@vueuse.c0256496.js";import"./dayjs.b9b14681.js";import"./axios.765908e4.js";import"./lodash-es.8ca9079f.js";import"./async-validator.ed4c92a2.js";import"./@ctrl.82a509e0.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.1c4ac15a.js";import"./vue-router.ae8a7199.js";import"./pinia.4a496763.js";import"./vue-demi.b3a9cad9.js";import"./screenfull.258ac312.js";import"./nprogress.aae2b787.js";const F=a=>(j("data-v-56ad595b"),a=a(),k(),a),P={class:"main-container"},S={class:"main",style:{"background-color":"inherit"}},w=F((()=>v("div",{class:"card-header"},[v("span",null,"在线升级")],-1))),C={class:"item"},V=F((()=>v("br",null,null,-1))),D=F((()=>v("br",null,null,-1))),U=F((()=>v("br",null,null,-1))),q=F((()=>v("br",null,null,-1))),z={style:{"text-align":"center"}},H=x(" 在线升级系统到最新版本 "),I=F((()=>v("div",{class:"card-header"},[v("span",null,"本地升级")],-1))),M={class:"item"},N=F((()=>v("div",null,"请选择您要升级的系统文件。",-1))),Q=F((()=>v("br",null,null,-1))),R=F((()=>v("br",null,null,-1))),T=F((()=>v("br",null,null,-1))),Y={style:{"text-align":"center"}},A=x(" 导入系统升级文件 "),B=x("选择文件"),E=F((()=>v("div",{class:"el-upload__tip"},"只能上传一个文件，只支持zip格式文件！",-1))),G={class:"dialog-footer"},J=x("取消"),K=x("上传");var L=r({setup(r){const j=i(),k=p({uFlag:!1,sysParams:{}});(()=>{const a={token:j.token,data:{}};d.getSysParams(a).then((a=>{console.log("getSysParams -> res",a),"0"===a.code&&(k.sysParams=a.data,console.log("getSysParams -> ctxData.sysParams",k.sysParams))}))})();const x=()=>{L.value&&L.value.submit()},F=a=>{},L=m(null),O=a=>{L.value.clearFiles();const e=a[0];L.value.handleStart(e)},W=a=>{console.log("beforeUpload -> file",a);const e=["application/x-tar"].includes(a.type);if(e){if(e){let e=new FormData;e.append("file",a);const s={token:j.token,contentType:"multipart/form-data",data:e};n.updateSystem(s).then((a=>{"0"===a.code?(o.success(a.message),k.uFlag=!1):$(a)}))}}else o({type:"error",message:"文件格式不正确,必须是tar文件！"})},X=()=>{k.uFlag=!1,L.value.clearFiles()},Z=()=>{X()},$=a=>{o({type:"error",message:a.message})};return(r,i)=>{const n=a,d=e,p=s,m=l,j=t;return f(),y("div",P,[v("div",S,[b(p,{class:"box-card",shadow:"hover"},{header:g((()=>[w])),default:g((()=>[v("div",C,[v("div",null,"当前系统版本："+_(h(k).sysParams.softVer)+" 。",1),V,v("div",null,"最新系统版本："+_(h(k).sysParams.softVer)+" 。",1),D,U,q,v("div",z,[b(d,{type:"primary",onClick:i[0]||(i[0]=a=>{o.info("功能升级中...")})},{default:g((()=>[b(n,{class:"el-input__icon"},{default:g((()=>[b(h(u))])),_:1}),H])),_:1})])])])),_:1}),b(p,{class:"box-card",shadow:"hover"},{header:g((()=>[I])),default:g((()=>[v("div",M,[N,Q,R,T,v("div",Y,[b(d,{type:"primary",onClick:i[1]||(i[1]=a=>(console.log("importSysConfig"),void(k.uFlag=!0)))},{default:g((()=>[b(n,{class:"el-input__icon"},{default:g((()=>[b(h(c))])),_:1}),A])),_:1})])])])),_:1})]),b(j,{modelValue:h(k).uFlag,"onUpdate:modelValue":i[2]||(i[2]=a=>h(k).uFlag=a),title:"上传系统升级文件",width:"600px","before-close":Z,"close-on-click-modal":!1},{footer:g((()=>[v("span",G,[b(d,{onClick:X},{default:g((()=>[J])),_:1}),b(d,{type:"primary",onClick:x},{default:g((()=>[K])),_:1})])])),default:g((()=>[b(m,{ref_key:"uploadRef",ref:L,action:"","auto-upload":!1,"http-request":F,limit:1,"on-exceed":O,"before-upload":W},{tip:g((()=>[E])),default:g((()=>[b(d,{type:"primary"},{default:g((()=>[B])),_:1})])),_:1},512)])),_:1},8,["modelValue"])])}}},[["__scopeId","data-v-56ad595b"]]);export{L as default};
