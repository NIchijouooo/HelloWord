import{e,m as l,n as a,p as o,o as t,A as r,B as m,l as d,C as i,K as u,v as s,w as n,q as p,g as c,L as f}from"./element-plus.2c9426ce.js";import{s as b,_ as g,u as h,v}from"./index.9c81e799.js";import{t as w}from"./vue3-barcode.6ca0e488.js";import{d as y}from"./@element-plus.e895b9e1.js";import{r as F,J as V,j as _,ae as x,o as P,c as k,a as A,O as B,R as C,u as j,n as N,V as R,S as q,W as U,X as M,P as S,aj as z,av as I,aw as T}from"./@vue.2c474d7f.js";import"./lodash-es.98e49362.js";import"./async-validator.fb49d0f5.js";import"./@vueuse.d618444d.js";import"./dayjs.2790329b.js";import"./axios.765908e4.js";import"./@ctrl.b082b0c1.js";import"./@popperjs.36402333.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.8aeb3683.js";import"./vue-router.3b7aa6c5.js";import"./pinia.39330980.js";import"./vue-demi.b3a9cad9.js";import"./screenfull.7cf96174.js";import"./vue-echarts.ebbe6a42.js";import"./resize-detector.4e96b72b.js";import"./echarts.2b8dcc79.js";import"./zrender.54f191d6.js";import"./tslib.60310f1a.js";import"./nprogress.aae2b787.js";var O=e=>b.request({url:"/network/mobile",method:"post",data:e}),L=e=>b.request({url:"/network/mobile",method:"put",data:e}),D=e=>b.request({url:"/network/mobile",method:"delete",data:e}),E=e=>b.request({url:"/network/mobiles",method:"get",data:e});const H=e=>(I("data-v-63570d32"),e=e(),T(),e),Z={class:"main-container"},$={class:"main"},K={class:"search-bar"},W={class:"tool-bar"},J=H((()=>A("div",null,"无数据",-1))),X={class:"mobile-tips"},G={class:"dialog-content"},Q=H((()=>A("div",{class:"form-title"},[A("div",{class:"tName"},"配置参数")],-1))),Y=H((()=>A("div",{class:"form-title"},[A("div",{class:"tName"},"通信参数")],-1))),ee=H((()=>A("div",{class:"form-title"},[A("div",{class:"tName"},"保活参数")],-1))),le={class:"dialog-footer"},ae={style:{display:"flex",width:"100%","justify-content":"center","text-align":"center"}};var oe=g({__name:"Mobile",setup(b){const g=h(),I=/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/,T=/^(?=^.{3,255}$)(http(s)?:\/\/)?(www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\d+)*(\/\w+\.\w+)*$/,H=(e,l,a)=>{console.log("validateIP"),I.test(l)||T.test(l)?a():a(new Error("IP格式错误！"))},oe=(e,l,a)=>{isNaN(Number(l))?a(new Error("只能输入数字！")):a()},te=F(null),re=V({moduleName:"",headerCellStyle:{background:v.primaryColor,color:v.fontWhiteColor,height:"54px"},cellStyle:{height:"48px"},tableMaxHeight:0,currentPage:1,pagesize:20,moduleTableData:[],showAddBtn:!0,mFlag:!1,mTitle:"添加移动网络",moduleForm:{name:"",model:"",flowAlarm:!1,flowAlarmValue:"",serialNumber:"",baudRate:"",dataBits:"",stopBits:"",parity:"",timeout:"",interval:"",polling:"",ipMaster:"",ipSlave:"",pollPeriod:"",offlineCnt:""},moduleRules:{name:[{required:!0,message:"模块名称不能为空",trigger:"blur"}],model:[{required:!0,message:"模块型号不能为空",trigger:"blur"}],serialNumber:[{required:!0,message:"串口号不能为空",trigger:"blur"}],baudRate:[{required:!0,message:"波特率不能为空",trigger:"blur"}],dataBits:[{required:!0,message:"数据位不能为空",trigger:"blur"}],stopBits:[{required:!0,message:"停止位不能为空",trigger:"blur"}],parity:[{required:!0,message:"校验位不能为空",trigger:"blur"}],flowAlarmValue:[{required:!0,message:"流量报警值不能为空",trigger:"blur"},{trigger:"blur",validator:oe}],timeout:[{required:!0,message:"超时时间不能为空",trigger:"blur"},{trigger:"blur",validator:oe}],interval:[{required:!0,message:"间隔时间不能为空",trigger:"blur"},{trigger:"blur",validator:oe}],polling:[{required:!0,message:"轮询时间不能为空",trigger:"blur"},{trigger:"blur",validator:oe}],ipMaster:[{required:!0,validator:H,trigger:"blur"}],ipSlave:[{validator:H,trigger:"blur"}],pollPeriod:[{required:!0,message:"保活检测周期不能为空",trigger:"blur"},{type:"number",message:"保活检测周期只能输入数字"}],offlineCnt:[{required:!0,message:"离线判断次数不能为空",trigger:"blur"},{type:"number",message:"离线判断次数只能输入数字"}]},baudRateOptions:[{value:"1200",label:"1200"},{value:"2400",label:"2400"},{value:"4800",label:"4800"},{value:"9600",label:"9600"},{value:"14400",label:"14400"},{value:"19200",label:"19200"},{value:"115200",label:"115200"}],dataBitsOptions:[{value:"5",label:"5"},{value:"6",label:"6"},{value:"7",label:"7"},{value:"8",label:"8"}],stopBitsOptions:[{value:"1",label:"1"},{value:"1.5",label:"1.5"},{value:"2",label:"2"}],parityOptions:[{value:"N",label:"无校验"},{value:"O",label:"奇校验"},{value:"E",label:"偶校验"}],sFlag:!1,barcodeVale:0}),me=e=>{const l={token:g.token,data:{}};E(l).then((async l=>{console.log("getModuleList -> res",l),l&&("0"===l.code?(re.moduleTableData=l.data,re.showAddBtn=0===l.data.length,1===e&&c({type:"success",message:"刷新成功！"})):ce(l),await N((()=>{re.tableMaxHeight=te.value.clientHeight-34-36-22})))}))};me();const de=_((()=>re.moduleTableData.filter((e=>!re.moduleName||e.name.toLowerCase().includes(re.moduleName.toLowerCase()))).slice((re.currentPage-1)*re.pagesize,re.currentPage*re.pagesize)));_((()=>re.moduleTableData.filter((e=>!re.moduleName||e.name.toLowerCase().includes(re.moduleName.toLowerCase())))));const ie=e=>{re.mFlag=!0,re.mTitle="编辑移动网络",re.moduleForm.name=e.name,re.moduleForm.model=e.model,re.moduleForm.flowAlarm=e.configParam.flowAlarm,re.moduleForm.flowAlarmValue=e.configParam.flowAlarmValue,re.moduleForm.serialNumber=e.commParam.name,re.moduleForm.baudRate=e.commParam.baudRate,re.moduleForm.dataBits=e.commParam.dataBits,re.moduleForm.stopBits=e.commParam.stopBits,re.moduleForm.parity=e.commParam.parity,re.moduleForm.timeout=e.commParam.timeout,re.moduleForm.interval=e.commParam.interval,re.moduleForm.polling=e.commParam.pollPeriod,re.moduleForm.ipMaster=e.keepAliveParam.ipMaster,re.moduleForm.ipSlave=e.keepAliveParam.ipSlave,re.moduleForm.pollPeriod=e.keepAliveParam.pollPeriod,re.moduleForm.offlineCnt=e.keepAliveParam.offlineCnt},ue=F(null),se=e=>{ne()},ne=()=>{re.mFlag=!1,ue.value.resetFields(),pe()},pe=()=>{re.moduleForm={name:"",model:"",flowAlarm:!1,flowAlarmValue:"",serialNumber:"",baudRate:"",dataBits:"",stopBits:"",parity:"",timeout:"",interval:"",polling:"",ipMaster:"",ipSlave:"",pollPeriod:"",offlineCnt:""}},ce=e=>{c({type:"error",message:e.message})},fe=(e,l)=>{c({type:"0"===e.code?"success":"error",message:e.message}),"0"===e.code&&l&&l()};return(b,h)=>{const v=e,F=l,V=a,_=x("Icon"),N=o,I=t,T=r,E=m,H=d,oe=i,pe=u,ce=s,be=n,ge=p;return P(),k("div",Z,[A("div",$,[A("div",K,[B(I,{inline:!0,ref:"searchFormRef","status-icon":"","label-width":"90px"},{default:C((()=>[B(V,{label:"模块名称"},{default:C((()=>[B(F,{style:{width:"200px"},placeholder:"请输入模块名称",clearable:"",modelValue:j(re).moduleName,"onUpdate:modelValue":h[0]||(h[0]=e=>j(re).moduleName=e)},{prefix:C((()=>[B(v,{class:"el-input__icon"},{default:C((()=>[B(j(y))])),_:1})])),_:1},8,["modelValue"])])),_:1}),B(V,null,{default:C((()=>[B(N,{style:{color:"#fff","margin-left":"20px"},color:"#2EA554",class:"right-btn",onClick:h[1]||(h[1]=e=>{me(1)})},{default:C((()=>[B(v,{class:"btn-icon"},{default:C((()=>[B(_,{name:"local-refresh",size:"14px",color:"#ffffff"})])),_:1}),R(" 刷新 ")])),_:1})])),_:1})])),_:1},512)]),A("div",W,[B(N,{type:"primary",bg:"",class:"right-btn",onClick:h[2]||(h[2]=e=>{re.mFlag=!0})},{default:C((()=>[B(v,{class:"btn-icon"},{default:C((()=>[B(_,{name:"local-add",size:"14px",color:"#ffffff"})])),_:1}),R(" 添加 ")])),_:1})]),A("div",{class:"content",ref_key:"contentRef",ref:te},[B(oe,{data:j(de),"cell-style":j(re).cellStyle,"header-cell-style":j(re).headerCellStyle,style:{width:"100%"},"max-height":j(re).tableMaxHeight,stripe:"",onRowDblclick:ie},{empty:C((()=>[J])),default:C((()=>[B(T,{sortable:"",prop:"name",label:"模块名称",width:"auto","min-width":"120",align:"center"}),B(T,{sortable:"",prop:"model",label:"模块型号",width:"auto","min-width":"120",align:"center"}),B(T,{sortable:"",label:"sim卡号",width:"auto","min-width":"180",align:"center"},{default:C((e=>[""!==e.row.runParam.iccid?(P(),q(H,{key:0,class:"box-item",effect:"dark",content:"单击显示SIM卡条形码",placement:"top-start"},{default:C((()=>[B(E,{style:{cursor:"pointer"},onClick:l=>{return a=e.row.runParam.iccid,re.sFlag=!0,void(re.barcodeVale=a);var a}},{default:C((()=>[R(U(e.row.runParam.iccid),1)])),_:2},1032,["onClick"])])),_:2},1024)):M("",!0)])),_:1}),B(T,{sortable:"",label:"imei",width:"auto","min-width":"150",align:"center"},{default:C((e=>[R(U(e.row.runParam.imei),1)])),_:1}),B(T,{sortable:"",label:"信号强度",width:"auto","min-width":"80",align:"center"},{default:C((e=>[R(U(e.row.runParam.csq),1)])),_:1}),B(T,{sortable:"",label:"流量（KByte）",width:"auto","min-width":"150",align:"center"},{default:C((e=>[R(U(e.row.runParam.flow),1)])),_:1}),B(T,{sortable:"",label:"基站定位Lac",width:"auto","min-width":"110",align:"center"},{default:C((e=>[R(U(e.row.runParam.lac),1)])),_:1}),B(T,{sortable:"",label:"基站定位Ci",width:"auto","min-width":"100",align:"center"},{default:C((e=>[R(U(e.row.runParam.ci),1)])),_:1}),B(T,{sortable:"",label:"sim卡插入状态",width:"auto","min-width":"120",align:"center"},{default:C((e=>[B(E,{type:e.row.runParam.simInsert?"success":"danger"},{default:C((()=>[R(U(e.row.runParam.simInsert?"是":"否"),1)])),_:2},1032,["type"])])),_:1}),B(T,{sortable:"",label:"网络注册状态",width:"auto","min-width":"120",align:"center"},{default:C((e=>[B(E,{type:e.row.runParam.netRegister?"success":"danger"},{default:C((()=>[R(U(e.row.runParam.netRegister?"是":"否"),1)])),_:2},1032,["type"])])),_:1}),B(T,{label:"操作",width:"auto","min-width":"150",align:"center",fixed:"right"},{default:C((e=>[B(N,{onClick:l=>ie(e.row),text:"",type:"primary"},{default:C((()=>[R("编辑")])),_:2},1032,["onClick"]),B(N,{onClick:l=>{return a=e.row,void f.confirm("确定要删除这个移动网络吗?","警告",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then((()=>{const e={token:g.token,data:{name:a.name}};D(e).then((e=>{fe(e,me)}))})).catch((()=>{c({type:"info",message:"取消删除"})}));var a},text:"",type:"danger"},{default:C((()=>[R("删除")])),_:2},1032,["onClick"])])),_:1})])),_:1},8,["data","cell-style","header-cell-style","max-height"]),A("div",X,[B(E,{class:"ml-2",type:"danger"},{default:C((()=>[R("注：移动网络配置完成后，必须重启网关，才能生效！")])),_:1})])],512)]),B(ge,{modelValue:j(re).mFlag,"onUpdate:modelValue":h[21]||(h[21]=e=>j(re).mFlag=e),title:j(re).mTitle,width:"900px","before-close":se,"close-on-click-modal":!1},{footer:C((()=>[A("span",le,[B(N,{onClick:h[19]||(h[19]=e=>ne())},{default:C((()=>[R("取消")])),_:1}),B(N,{type:"primary",onClick:h[20]||(h[20]=e=>{ue.value.validate((e=>{const l={name:re.moduleForm.name,model:re.moduleForm.model,commParam:{},configParam:{}},a={name:re.moduleForm.serialNumber,baudRate:re.moduleForm.baudRate,dataBits:re.moduleForm.dataBits,stopBits:re.moduleForm.stopBits,parity:re.moduleForm.parity,timeout:re.moduleForm.timeout,interval:re.moduleForm.interval,pollPeriod:re.moduleForm.polling},o={flowAlarm:re.moduleForm.flowAlarm,flowAlarmValue:re.moduleForm.flowAlarmValue},t={ipMaster:re.moduleForm.ipMaster,ipSlave:re.moduleForm.ipSlave,pollPeriod:re.moduleForm.pollPeriod,offlineCnt:re.moduleForm.offlineCnt};if(l.commParam=a,l.configParam=o,l.keepAliveParam=t,!e)return!1;{const e={token:g.token,data:l};re.mTitle.includes("添加")?O(e).then((e=>{fe(e,me),ne()})):L(e).then((e=>{fe(e,me),ne()}))}}))})},{default:C((()=>[R("保存")])),_:1})])])),default:C((()=>[A("div",G,[B(I,{model:j(re).moduleForm,rules:j(re).moduleRules,ref_key:"moduleFormRef",ref:ue,"status-icon":"","label-position":"right","label-width":"120px",inline:"true"},{default:C((()=>[B(V,{label:"模块名称",prop:"name"},{default:C((()=>[B(F,{disabled:j(re).mTitle.includes("编辑"),type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.name,"onUpdate:modelValue":h[3]||(h[3]=e=>j(re).moduleForm.name=e),autocomplete:"off",placeholder:"请输入模块名称"},null,8,["disabled","modelValue"])])),_:1}),B(V,{label:"模块型号",prop:"model"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.model,"onUpdate:modelValue":h[4]||(h[4]=e=>j(re).moduleForm.model=e),autocomplete:"off",placeholder:"请输入模块型号"},null,8,["modelValue"])])),_:1}),Q,B(V,{label:"启用流量报警",prop:"flowAlarm"},{default:C((()=>[B(pe,{style:{width:"270px"},modelValue:j(re).moduleForm.flowAlarm,"onUpdate:modelValue":h[5]||(h[5]=e=>j(re).moduleForm.flowAlarm=e),"inline-prompt":"","active-text":"是","inactive-text":"否"},null,8,["modelValue"])])),_:1}),B(V,{label:"流量报警值",prop:j(re).moduleForm.flowAlarm?"flowAlarmValue":""},{default:C((()=>[B(F,{disabled:!j(re).moduleForm.flowAlarm,type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.flowAlarmValue,"onUpdate:modelValue":h[6]||(h[6]=e=>j(re).moduleForm.flowAlarmValue=e),autocomplete:"off",placeholder:"请输入流量报警值"},null,8,["disabled","modelValue"])])),_:1},8,["prop"]),Y,B(V,{label:"串口号",prop:"serialNumber"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.serialNumber,"onUpdate:modelValue":h[7]||(h[7]=e=>j(re).moduleForm.serialNumber=e),autocomplete:"off",placeholder:"请输入串口名称"},null,8,["modelValue"])])),_:1}),B(V,{label:"波特率",prop:"baudRate"},{default:C((()=>[B(be,{modelValue:j(re).moduleForm.baudRate,"onUpdate:modelValue":h[8]||(h[8]=e=>j(re).moduleForm.baudRate=e),style:{width:"270px"},placeholder:"请选择波特率"},{default:C((()=>[(P(!0),k(S,null,z(j(re).baudRateOptions,(e=>(P(),q(ce,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),B(V,{label:"数据位",prop:"dataBits"},{default:C((()=>[B(be,{modelValue:j(re).moduleForm.dataBits,"onUpdate:modelValue":h[9]||(h[9]=e=>j(re).moduleForm.dataBits=e),style:{width:"270px"},placeholder:"请选择数据位"},{default:C((()=>[(P(!0),k(S,null,z(j(re).dataBitsOptions,(e=>(P(),q(ce,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),B(V,{label:"停止位",prop:"stopBits"},{default:C((()=>[B(be,{modelValue:j(re).moduleForm.stopBits,"onUpdate:modelValue":h[10]||(h[10]=e=>j(re).moduleForm.stopBits=e),style:{width:"270px"},placeholder:"请选择停止位"},{default:C((()=>[(P(!0),k(S,null,z(j(re).stopBitsOptions,(e=>(P(),q(ce,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),B(V,{label:"校验位",prop:"parity"},{default:C((()=>[B(be,{modelValue:j(re).moduleForm.parity,"onUpdate:modelValue":h[11]||(h[11]=e=>j(re).moduleForm.parity=e),style:{width:"270px"},placeholder:"请选择校验位"},{default:C((()=>[(P(!0),k(S,null,z(j(re).parityOptions,(e=>(P(),q(ce,{key:e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),B(V,{label:"超时时间",prop:"timeout"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.timeout,"onUpdate:modelValue":h[12]||(h[12]=e=>j(re).moduleForm.timeout=e),autocomplete:"off",placeholder:"请输入超时时间"},{append:C((()=>[R("单位毫秒")])),_:1},8,["modelValue"])])),_:1}),B(V,{label:"间隔时间",prop:"interval"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.interval,"onUpdate:modelValue":h[13]||(h[13]=e=>j(re).moduleForm.interval=e),autocomplete:"off",placeholder:"请输入间隔时间"},{append:C((()=>[R("单位毫秒")])),_:1},8,["modelValue"])])),_:1}),B(V,{label:"轮询时间",prop:"polling"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.polling,"onUpdate:modelValue":h[14]||(h[14]=e=>j(re).moduleForm.polling=e),autocomplete:"off",placeholder:"请输入轮询时间"},{append:C((()=>[R("单位秒")])),_:1},8,["modelValue"])])),_:1}),ee,B(V,{label:"保活检测主IP",prop:"ipMaster"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.ipMaster,"onUpdate:modelValue":h[15]||(h[15]=e=>j(re).moduleForm.ipMaster=e),autocomplete:"off",placeholder:"请输入保活检测主IP"},null,8,["modelValue"])])),_:1}),B(V,{label:"保活检测备用IP",prop:"ipSlave"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.ipSlave,"onUpdate:modelValue":h[16]||(h[16]=e=>j(re).moduleForm.ipSlave=e),autocomplete:"off",placeholder:"请输入保活检测备用IP"},null,8,["modelValue"])])),_:1}),B(V,{label:"保活检测周期",prop:"pollPeriod"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.pollPeriod,"onUpdate:modelValue":h[17]||(h[17]=e=>j(re).moduleForm.pollPeriod=e),modelModifiers:{number:!0},autocomplete:"off",placeholder:"请输入保活检测周期"},{append:C((()=>[R("单位秒")])),_:1},8,["modelValue"])])),_:1}),B(V,{label:"离线判断次数",prop:"offlineCnt"},{default:C((()=>[B(F,{type:"text",style:{width:"270px"},modelValue:j(re).moduleForm.offlineCnt,"onUpdate:modelValue":h[18]||(h[18]=e=>j(re).moduleForm.offlineCnt=e),modelModifiers:{number:!0},autocomplete:"off",placeholder:"请输入离线判断次数"},{append:C((()=>[R("单位次")])),_:1},8,["modelValue"])])),_:1})])),_:1},8,["model","rules"])])])),_:1},8,["modelValue","title"]),B(ge,{modelValue:j(re).sFlag,"onUpdate:modelValue":h[22]||(h[22]=e=>j(re).sFlag=e),title:"sim卡号条形码",width:"30%"},{default:C((()=>[A("div",ae,[B(j(w),{value:j(re).barcodeVale,height:60},null,8,["value"])])])),_:1},8,["modelValue"])])}}},[["__scopeId","data-v-63570d32"]]);export{oe as default};
