import{K as e,n as a,v as l,w as t,m as o,p as r,o as s,u as d,F as u,G as p,g as n}from"./element-plus.2c9426ce.js";import{_ as m,u as i}from"./index.9c81e799.js";import{S as f}from"./sysTool.087279d2.js";import{J as c,r as b,o as v,c as T,a as U,O as C,R as g,u as _,P as h,aj as j,V as y,S as V,av as w,aw as F}from"./@vue.2c474d7f.js";import"./lodash-es.98e49362.js";import"./async-validator.fb49d0f5.js";import"./@vueuse.d618444d.js";import"./dayjs.2790329b.js";import"./axios.765908e4.js";import"./@ctrl.b082b0c1.js";import"./@popperjs.36402333.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.8aeb3683.js";import"./vue-router.3b7aa6c5.js";import"./pinia.39330980.js";import"./vue-demi.b3a9cad9.js";import"./@element-plus.e895b9e1.js";import"./screenfull.7cf96174.js";import"./vue-echarts.ebbe6a42.js";import"./resize-detector.4e96b72b.js";import"./echarts.2b8dcc79.js";import"./zrender.54f191d6.js";import"./tslib.60310f1a.js";import"./nprogress.aae2b787.js";const k={class:"main-container"},M={class:"main",style:{"background-color":"inherit"}},x=(e=>(w("data-v-382cb360"),e=e(),F(),e))((()=>U("div",{class:"card-header"},[U("span",null,"NTP校时")],-1))),S={class:"item"},Z={style:{"padding-right":"14px",height:"100%",overflow:"auto"}},z={class:"remark"};var P=m({__name:"NTPTiming",setup(m){const w=i(),F=/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/,P=/^(?=^.{3,255}$)(http(s)?:\/\/)?(www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\d+)*(\/\w+\.\w+)*$/,N=(e,a,l)=>{console.log("validateIP"),""!==a?F.test(a)||P.test(a)?l():l(new Error("格式错误！")):l()},A=c({showBtn:!1,timeZoneOptions:[{value:"UTC+8",label:"UTC+8"},{value:"UTC+7",label:"UTC+7"},{value:"UTC+6",label:"UTC+6"},{value:"UTC+5",label:"UTC+5"},{value:"UTC+4",label:"UTC+4"},{value:"UTC+3",label:"UTC+3"},{value:"UTC+2",label:"UTC+2"},{value:"UTC+1",label:"UTC+1"},{value:"UTC+0",label:"UTC+0"},{value:"UTC-1",label:"UTC-1"},{value:"UTC-2",label:"UTC-2"},{value:"UTC-3",label:"UTC-3"},{value:"UTC-4",label:"UTC-4"},{value:"UTC-5",label:"UTC-5"},{value:"UTC-6",label:"UTC-6"},{value:"UTC-7",label:"UTC-7"},{value:"UTC-8",label:"UTC-8"}],ntpForm:{enable:"",timeZone:"",urlMaster:"",portMaster:null,urlSlave:"",portSlave:null},ntpRules:{urlMaster:[{validator:N,trigger:"blur"}],portMaster:[{type:"number",message:"必须是数字",trigger:"blur"}],urlSlave:[{validator:N,trigger:"blur"}],portSlave:[{type:"number",message:"必须是数字",trigger:"blur"}]}}),B=()=>{const e={token:w.token,data:{},mock:!0};f.getNTPInfo(e).then((e=>{"0"===e.code?(A.ntpForm=e.data,n({type:"success",message:"获取NTP服务成功！"}),A.showBtn=!1):O(e)}))};B();const R=b(null),I=(e,a)=>{n({type:"0"===e.code?"success":"error",message:e.message}),"0"===e.code&&a&&a()},O=e=>{n({type:"error",message:e.message})};return(n,m)=>{const i=e,c=a,b=l,F=t,P=o,N=r,O=s,$=d,E=u,G=p;return v(),T("div",k,[U("div",M,[C(G,{gutter:20,style:{height:"100%"}},{default:g((()=>[C(E,{span:24},{default:g((()=>[C($,{class:"box-card",shadow:"never"},{header:g((()=>[x])),default:g((()=>[U("div",S,[C($,{style:{height:"100%"}},{default:g((()=>[U("div",Z,[C(O,{model:_(A).ntpForm,rules:_(A).ntpRules,ref_key:"ntpRef",ref:R,"status-icon":"","label-position":"right","label-width":"120px"},{default:g((()=>[C(c,{label:"启用状态"},{default:g((()=>[C(i,{modelValue:_(A).ntpForm.enable,"onUpdate:modelValue":m[0]||(m[0]=e=>_(A).ntpForm.enable=e),"inline-prompt":"","active-text":"是","inactive-text":"否"},null,8,["modelValue"])])),_:1}),C(c,{label:"时区"},{default:g((()=>[C(F,{modelValue:_(A).ntpForm.timeZone,"onUpdate:modelValue":m[1]||(m[1]=e=>_(A).ntpForm.timeZone=e),modelModifiers:{number:!0},style:{width:"100%"},placeholder:"请选时区"},{default:g((()=>[(v(!0),T(h,null,j(_(A).timeZoneOptions,(e=>(v(),V(b,{key:"type_"+e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),C(c,{label:"主服务器地址",prop:"urlMaster"},{default:g((()=>[C(P,{type:"text",modelValue:_(A).ntpForm.urlMaster,"onUpdate:modelValue":m[2]||(m[2]=e=>_(A).ntpForm.urlMaster=e),autocomplete:"off",placeholder:"请输入主服务器地址"},null,8,["modelValue"])])),_:1}),C(c,{label:"主服务器端口",prop:"portMaster"},{default:g((()=>[C(P,{type:"text",modelValue:_(A).ntpForm.portMaster,"onUpdate:modelValue":m[3]||(m[3]=e=>_(A).ntpForm.portMaster=e),modelModifiers:{number:!0},autocomplete:"off",placeholder:"请输入主服务器端口"},null,8,["modelValue"])])),_:1}),C(c,{label:"次服务器地址",prop:"urlSlave"},{default:g((()=>[C(P,{type:"text",modelValue:_(A).ntpForm.urlSlave,"onUpdate:modelValue":m[4]||(m[4]=e=>_(A).ntpForm.urlSlave=e),autocomplete:"off",placeholder:"请输入次服务器地址"},null,8,["modelValue"])])),_:1}),C(c,{label:"次服务器端口",prop:"portSlave"},{default:g((()=>[C(P,{type:"text",modelValue:_(A).ntpForm.portSlave,"onUpdate:modelValue":m[5]||(m[5]=e=>_(A).ntpForm.portSlave=e),modelModifiers:{number:!0},autocomplete:"off",placeholder:"请输入次服务器端口"},null,8,["modelValue"])])),_:1}),C(c,null,{default:g((()=>[C(N,{disabled:_(A).showBtn,type:"primary",onClick:m[6]||(m[6]=e=>{R.value.validate((e=>{if(!e)return!1;{A.showBtn=!0;const e={token:w.token,data:A.ntpForm,mock:!0};f.updateNTP(e).then((e=>{I(e,B)}))}}))})},{default:g((()=>[y("保存")])),_:1},8,["disabled"]),C(N,{type:"primary",onClick:m[7]||(m[7]=e=>(()=>{const e={token:w.token,data:{}};f.sendTimingCmd(e).then((e=>{I(e,B)}))})())},{default:g((()=>[y("立即校时")])),_:1})])),_:1})])),_:1},8,["model","rules"])])])),_:1}),U("div",z,[C(G,{gutter:16},{default:g((()=>[C(E,{span:1},{default:g((()=>[y("操作步骤:")])),_:1}),C(E,{span:21},{default:g((()=>[y("1、选择时区（默认是UTC+8时区）")])),_:1})])),_:1}),C(G,{gutter:16},{default:g((()=>[C(E,{offset:1,span:21},{default:g((()=>[y("2、输入主服务器地址、端口")])),_:1})])),_:1}),C(G,{gutter:16},{default:g((()=>[C(E,{offset:1,span:21},{default:g((()=>[y("3、输入次服务器地址、端口")])),_:1})])),_:1}),C(G,{gutter:16},{default:g((()=>[C(E,{offset:1,span:21},{default:g((()=>[y("4、保存配置，重启后生肖")])),_:1})])),_:1}),C(G,{gutter:16},{default:g((()=>[C(E,{offset:1,span:21},{default:g((()=>[y("5、立即校时，系统时间立即更新")])),_:1})])),_:1})])])])),_:1})])),_:1})])),_:1})])])}}},[["__scopeId","data-v-382cb360"]]);export{P as default};
