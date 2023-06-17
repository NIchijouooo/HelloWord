import{e,m as a,n as t,p as l,o,A as r,C as i,D as n,v as c,w as s,q as m,g as d,L as f}from"./element-plus.3300f756.js";import{_ as p,u,v as g}from"./index.559546be.js";import{I as b}from"./interface.5c9ce411.js";import{C as h}from"./commInterface.19fb3003.js";import{K as y,f as N,k as w,af as I,o as _,c as v,a as C,P as T,S as x,u as j,n as F,W as k,Q as P,ak as L,T as V,aw as U,ax as z}from"./@vue.73d14ebc.js";import"./lodash-es.98e49362.js";import"./async-validator.fb49d0f5.js";import"./@vueuse.3d90a471.js";import"./dayjs.2790329b.js";import"./axios.765908e4.js";import"./@ctrl.b082b0c1.js";import"./@popperjs.36402333.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.8aeb3683.js";import"./vue-router.43b9a342.js";import"./pinia.3f4aff4f.js";import"./vue-demi.b3a9cad9.js";import"./@element-plus.0bef5fa5.js";import"./screenfull.7cf96174.js";import"./vue-echarts.7828944a.js";import"./resize-detector.4e96b72b.js";import"./echarts.2b8dcc79.js";import"./zrender.54f191d6.js";import"./tslib.60310f1a.js";import"./nprogress.aae2b787.js";const D={class:"main-container"},S={class:"main"},M={class:"search-bar"},R={class:"tool-bar"},B=(e=>(U("data-v-15301b6a"),e=e(),z(),e))((()=>C("div",null,"无数据",-1))),q={class:"pagination"},A={class:"dialog-content"},O={class:"dialog-footer"};var H=p({__name:"Interface",setup(p){const U=u(),z=y({interfaceName:"",checkAll:!0,headerCellStyle:{background:g.primaryColor,color:g.fontWhiteColor,height:"54px"},cellStyle:{height:"48px"},tableMaxHeight:0,currentPage:1,pagesize:20,interfaceTableData:[],iFlag:!1,iTitle:"添加采集接口",interfaceForm:{collInterfaceName:"",commInterfaceName:"",protocolTypeName:"",pollPeriod:0,offlinePeriod:0},interfaceRules:{collInterfaceName:[{required:!0,message:"采集接口不能为空",trigger:"blur"}],commInterfaceName:[{required:!0,message:"通讯接口不能为空",trigger:"blur"}],pollPeriod:[{required:!0,message:"采集周期不能为空",trigger:"blur"},{type:"number",message:"采集周期只能输入数字"}],offlinePeriod:[{required:!0,message:"离线判断次数不能为空",trigger:"blur"},{type:"number",message:"离线判断次数只能输入数字"}]},commInterfaceList:[],protocolTypeNameList:[{name:"NULL",label:"NULL"},{name:"LUA",label:"LUA"},{name:"DLT645-2007",label:"DLT645-2007"},{name:"MODBUS-RTU",label:"MODBUS-RTU"},{name:"MODBUS-TCP",label:"MODBUS-TCP"}]}),H=N(null),E=e=>{const a={token:U.token,data:{}};b.getInterfaceList(a).then((async a=>{console.log("getInterfaceList -> res",a),a&&("0"===a.code?(z.interfaceTableData=a.data,z.interfaceTableData.forEach((e=>{e.protocolTypeName||(e.protocolTypeName=z.protocolTypeNameList[0].name)})),1===e&&d({type:"success",message:"刷新成功！"})):Z(a),await F((()=>{z.tableMaxHeight=H.value.clientHeight-34-36-22})))}))};E();const W=w((()=>z.interfaceTableData.filter((e=>{var a=!z.interfaceName,t=e.collInterfaceName.toLowerCase().includes(z.interfaceName.toLowerCase());return a||t})).slice((z.currentPage-1)*z.pagesize,z.currentPage*z.pagesize))),K=w((()=>z.interfaceTableData.filter((e=>{var a=!z.interfaceName,t=e.collInterfaceName.toLowerCase().includes(z.interfaceName.toLowerCase());return a||t}))));(()=>{const e={token:U.token,data:{}};h.getCommInterfaceList(e).then((e=>{e&&(console.log("getCommInterfaceList -> res",e),"0"===e.code?z.commInterfaceList=e.data:Z(e))}))})();const Q=e=>{z.iFlag=!0,z.iTitle="编辑采集接口",z.interfaceForm.collInterfaceName=e.collInterfaceName,z.interfaceForm.commInterfaceName=e.commInterfaceName,z.interfaceForm.protocolTypeName=e.protocolTypeName,z.interfaceForm.pollPeriod=e.pollPeriod,z.interfaceForm.offlinePeriod=e.offlinePeriod},G=N(null),J=()=>{z.iFlag=!1,G.value.resetFields(),Y()},X=e=>{J()},Y=()=>{z.interfaceForm={collInterfaceName:"",commInterfaceName:"",protocolTypeName:"",pollPeriod:0,offlinePeriod:0}},Z=e=>{d({type:"error",message:e.message})},$=(e,a)=>{d({type:"0"===e.code?"success":"1"===e.code?"error":"warning",message:e.message}),"0"===e.code&&a&&a()};return(p,u)=>{const g=I("search"),h=e,y=a,N=t,w=I("Icon"),F=l,Y=o,Z=r,ee=i,ae=n,te=c,le=s,oe=m;return _(),v("div",D,[C("div",S,[C("div",M,[T(Y,{inline:!0,ref:"searchFormRef","status-icon":"","label-width":"120px"},{default:x((()=>[T(N,{label:"采集接口名称"},{default:x((()=>[T(y,{style:{width:"200px"},placeholder:"请输入采集接口名称",modelValue:z.interfaceName,"onUpdate:modelValue":u[0]||(u[0]=e=>z.interfaceName=e)},{prefix:x((()=>[T(h,{class:"el-input__icon"},{default:x((()=>[T(g)])),_:1})])),_:1},8,["modelValue"])])),_:1}),T(N,null,{default:x((()=>[T(F,{style:{color:"#fff","margin-left":"20px"},color:"#2EA554",class:"right-btn",onClick:u[1]||(u[1]=e=>{E(1)})},{default:x((()=>[T(h,{class:"btn-icon"},{default:x((()=>[T(w,{name:"local-refresh",size:"14px",color:"#ffffff"})])),_:1}),k(" 刷新 ")])),_:1})])),_:1})])),_:1},512)]),C("div",R,[T(F,{type:"primary",bg:"",class:"right-btn",onClick:u[2]||(u[2]=e=>(z.iFlag=!0,void(z.iTitle="添加采集接口")))},{default:x((()=>[T(h,{class:"btn-icon"},{default:x((()=>[T(w,{name:"local-add",size:"14px",color:"#ffffff"})])),_:1}),k(" 添加 ")])),_:1})]),C("div",{class:"content",ref_key:"contentRef",ref:H},[T(ee,{data:j(W),"cell-style":z.cellStyle,"header-cell-style":z.headerCellStyle,"max-height":z.tableMaxHeight,style:{width:"100%"},stripe:"",onRowDblclick:Q},{empty:x((()=>[B])),default:x((()=>[T(Z,{sortable:"",prop:"collInterfaceName",label:"采集接口名称",width:"auto","min-width":"150",align:"center"}),T(Z,{sortable:"",prop:"commInterfaceName",label:"通讯接口",width:"auto","min-width":"150",align:"center"}),T(Z,{sortable:"",prop:"protocolTypeName",label:"通讯协议",width:"auto","min-width":"150",align:"center"}),T(Z,{sortable:"",prop:"pollPeriod",label:"采集周期(秒)",width:"auto","min-width":"150",align:"center"}),T(Z,{sortable:"",prop:"offlinePeriod",label:"离线判断次数",width:"auto","min-width":"150",align:"center"}),T(Z,{sortable:"",prop:"deviceNodeCnt",label:"设备总数",width:"auto","min-width":"150",align:"center"}),T(Z,{sortable:"",prop:"deviceNodeOnlineCnt",label:"设备在线数",width:"auto","min-width":"150",align:"center"}),T(Z,{label:"操作",width:"auto","min-width":"300",align:"center",fixed:"right"},{default:x((e=>[T(F,{onClick:a=>Q(e.row),text:"",type:"primary"},{default:x((()=>[k("编辑")])),_:2},1032,["onClick"]),T(F,{onClick:a=>{return t=e.row,void f.confirm("确定要删除这个采集接口吗?","警告",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then((()=>{const e={token:U.token,data:{name:t.collInterfaceName}};b.deleteInterface(e).then((e=>{$(e,E)}))})).catch((()=>{d({type:"info",message:"取消删除"})}));var t},text:"",type:"danger"},{default:x((()=>[k("删除")])),_:2},1032,["onClick"])])),_:1})])),_:1},8,["data","cell-style","header-cell-style","max-height"]),C("div",q,[T(ae,{"current-page":z.currentPage,"page-size":z.pagesize,"page-sizes":[20,50,200,500],total:j(K).length,onCurrentChange:p.handleCurrentChange,onSizeChange:p.handleSizeChange,background:"",layout:"total, sizes, prev, pager, next, jumper",style:{"margin-top":"46px"}},null,8,["current-page","page-size","total","onCurrentChange","onSizeChange"])])],512)]),T(oe,{modelValue:z.iFlag,"onUpdate:modelValue":u[10]||(u[10]=e=>z.iFlag=e),title:z.iTitle,width:"600px","before-close":X,"close-on-click-modal":!1},{footer:x((()=>[C("span",O,[T(F,{onClick:u[8]||(u[8]=e=>J())},{default:x((()=>[k("取消")])),_:1}),T(F,{type:"primary",onClick:u[9]||(u[9]=e=>{G.value.validate((e=>{if(!e)return!1;{const e={token:U.token,data:z.interfaceForm};z.iTitle.includes("添加")?b.addInterface(e).then((e=>{$(e,E),J()})):b.editInterface(e).then((e=>{$(e,E),J()}))}}))})},{default:x((()=>[k("保存")])),_:1})])])),default:x((()=>[C("div",A,[T(Y,{model:z.interfaceForm,rules:z.interfaceRules,ref_key:"interfaceFormRef",ref:G,"status-icon":"","label-position":"right","label-width":"120px"},{default:x((()=>[T(N,{label:"采集接口名称",prop:"collInterfaceName"},{default:x((()=>[T(y,{disabled:z.iTitle.includes("编辑"),type:"text",modelValue:z.interfaceForm.collInterfaceName,"onUpdate:modelValue":u[3]||(u[3]=e=>z.interfaceForm.collInterfaceName=e),autocomplete:"off",placeholder:"请输入采集接口名称"},null,8,["disabled","modelValue"])])),_:1}),T(N,{label:"通讯接口名称",prop:"commInterfaceName"},{default:x((()=>[T(le,{modelValue:z.interfaceForm.commInterfaceName,"onUpdate:modelValue":u[4]||(u[4]=e=>z.interfaceForm.commInterfaceName=e),style:{width:"100%"},placeholder:"请选择通讯接口名称"},{default:x((()=>[(_(!0),v(P,null,L(z.commInterfaceList,((e,a)=>(_(),V(te,{key:"comm_"+a,label:e.name,value:e.name},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),T(N,{label:"通讯协议",prop:"protocolTypeName"},{default:x((()=>[T(le,{modelValue:z.interfaceForm.protocolTypeName,"onUpdate:modelValue":u[5]||(u[5]=e=>z.interfaceForm.protocolTypeName=e),style:{width:"100%"},placeholder:"请选择通讯协议"},{default:x((()=>[(_(!0),v(P,null,L(z.protocolTypeNameList,(e=>(_(),V(te,{key:e.name,label:e.label,value:e.name},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),T(N,{label:"采集周期",prop:"pollPeriod"},{default:x((()=>[T(y,{type:"text",modelValue:z.interfaceForm.pollPeriod,"onUpdate:modelValue":u[6]||(u[6]=e=>z.interfaceForm.pollPeriod=e),modelModifiers:{number:!0},autocomplete:"off",placeholder:"请输入采集周期"},{append:x((()=>[k("单位秒")])),_:1},8,["modelValue"])])),_:1}),T(N,{label:"离线判断次数",prop:"offlinePeriod"},{default:x((()=>[T(y,{type:"text",modelValue:z.interfaceForm.offlinePeriod,"onUpdate:modelValue":u[7]||(u[7]=e=>z.interfaceForm.offlinePeriod=e),modelModifiers:{number:!0},autocomplete:"off",placeholder:"请输入离线判断次数"},{append:x((()=>[k("单位次")])),_:1},8,["modelValue"])])),_:1})])),_:1},8,["model","rules"])])])),_:1},8,["modelValue","title"])])}}},[["__scopeId","data-v-15301b6a"]]);export{H as default};
