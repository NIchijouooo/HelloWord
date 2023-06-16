import{U as e,v as a,w as t,e as l,p as o,A as s,B as c,C as n,D as r,g as i}from"./element-plus.3300f756.js";import{_ as p,u as d,v as g}from"./index.f61f902c.js";import{I as m}from"./interface.868ad1a1.js";import{S as u}from"./service.f9434ad8.js";import{P as h}from"./papaparse.0ea4a40d.js";import{d as f}from"./dayjs.2790329b.js";import{j as b,k as y,_ as w,h as k}from"./@element-plus.0bef5fa5.js";import{f as v,K as _,n as j,k as S,h as T,J as N,o as C,c as z,a as L,P as x,S as F,u as I,X as P,T as U,Y as H,W as V,Q as D,ak as M,aw as R,ax as Y}from"./@vue.73d14ebc.js";import"./lodash-es.98e49362.js";import"./async-validator.fb49d0f5.js";import"./@vueuse.3d90a471.js";import"./@ctrl.b082b0c1.js";import"./@popperjs.36402333.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.8aeb3683.js";import"./axios.765908e4.js";import"./vue-router.43b9a342.js";import"./pinia.3f4aff4f.js";import"./vue-demi.b3a9cad9.js";import"./screenfull.7cf96174.js";import"./vue-echarts.7828944a.js";import"./resize-detector.4e96b72b.js";import"./echarts.2b8dcc79.js";import"./zrender.54f191d6.js";import"./tslib.60310f1a.js";import"./nprogress.aae2b787.js";const W={class:"SysLog"},E={class:"title"},O={class:"content"},B={class:"cTitle"},J={class:"tName"},A={class:"option"},G=(e=>(R("data-v-0b8f00d5"),e=e(),Y(),e))((()=>L("div",null,"无数据",-1))),K={class:"pagination"};var Q=p({__name:"SysLog",setup(p){const R=d();console.log("location -> ",location);const Y=v(null),Q=_({websocket:null,headerCellStyle:{background:g.primaryColor,color:g.fontWhiteColor,height:"54px"},cellStyle:{height:"48px"},tableMaxHeight:0,currentPage:1,pagesize:20,logType:0,logName:["采集日志","上报日志","系统日志"],collectName:"",collectList:[],uploadName:"",uploadList:[],startStopFlag:!0,dataTable:[],dataType:["接收方","发送方"]}),X=e=>{console.log("onChange -> val",e),Q.logType=e,console.log("ctxData.startStopFlag = ",Q.startStopFlag),Q.startStopFlag||(Q.startStopFlag=!0,ee())},q=e=>{Q.startStopFlag=e,!1===e?Z():ee()},Z=()=>{let e="";e=location.origin.replace("https://","").replace("http://","");let a="ws://"+e+"/api/ws"+("?type="+Q.logType+"&name="+(0===Q.logType?Q.collectName:1===Q.logType?Q.uploadName:"systemLog"));console.log("url = "+a),Q.websocket=new WebSocket(a),$(),ae(),te()},$=()=>{Q.websocket.onmessage=function(e){console.log("event.data",e.data),e.data&&""!=e.data&&Q.dataTable.push(JSON.parse(e.data))}},ee=()=>{Q.websocket&&(console.log("关闭ws命令"),Q.websocket.close(),console.log("closeWs",Q.websocket))},ae=()=>{Q.websocket.onclose=e=>{console.log("onClose - > e",e),Q.startStopFlag=!0,Q.websocket=null}},te=()=>{Q.websocket.onerror=e=>{console.log("onError - > e",e),Q.startStopFlag=!0,Q.websocket=null}};(()=>{const e={token:R.token,data:{}};m.getInterfaceList(e).then((e=>{e&&("0"===e.code?(Q.collectList=e.data,e.data.length>0&&(Q.collectName=e.data[0].collInterfaceName)):ce(e))}))})();(()=>{const e={token:R.token,data:{}};u.getGatewayList(e).then((e=>{e&&("0"===e.code?(Q.uploadList=e.data,e.data.length>0&&(Q.uploadName=e.data[0].serviceName)):ce(e))}))})(),j((()=>{Q.tableMaxHeight=Y.value.clientHeight-34-36-22}));const le=e=>{Q.currentPage=e},oe=e=>{Q.pagesize=e},se=S((()=>Q.dataTable.slice((Q.currentPage-1)*Q.pagesize,Q.currentPage*Q.pagesize)));T((()=>Q.collectName),((e,a)=>{Q.startStopFlag=!0,ee()})),T((()=>Q.uploadName),((e,a)=>{Q.startStopFlag=!0,ee()}));const ce=e=>{i({type:"error",message:e.message})};return N((()=>{console.log("页面关闭了"),ee()})),(i,p)=>{const d=e,g=a,m=t,u=l,v=o,_=s,j=c,S=n,T=r;return C(),z("div",W,[L("div",E,[x(d,{size:"large",checked:0===I(Q).logType,onChange:p[0]||(p[0]=e=>X(0))},{default:F((()=>[V("采集日志")])),_:1},8,["checked"]),x(d,{size:"large",checked:1===I(Q).logType,onChange:p[1]||(p[1]=e=>X(1))},{default:F((()=>[V("上报日志")])),_:1},8,["checked"]),x(d,{size:"large",checked:2===I(Q).logType,onChange:p[2]||(p[2]=e=>X(2))},{default:F((()=>[V("系统日志")])),_:1},8,["checked"])]),L("div",O,[L("div",B,[L("div",J,P(I(Q).logName[I(Q).logType]),1),L("div",A,[0===I(Q).logType?(C(),U(m,{key:0,modelValue:I(Q).collectName,"onUpdate:modelValue":p[3]||(p[3]=e=>I(Q).collectName=e),style:{width:"100%"},placeholder:"请选择采集接口"},{default:F((()=>[(C(!0),z(D,null,M(I(Q).collectList,(e=>(C(),U(g,{key:"collect_"+e.collInterfaceName,label:e.collInterfaceName,value:e.collInterfaceName},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])):H("",!0),1===I(Q).logType?(C(),U(m,{key:1,modelValue:I(Q).uploadName,"onUpdate:modelValue":p[4]||(p[4]=e=>I(Q).uploadName=e),style:{width:"100%"},placeholder:"请选择上报接口"},{default:F((()=>[(C(!0),z(D,null,M(I(Q).uploadList,(e=>(C(),U(g,{key:"upload_"+e.serviceName,label:e.serviceName,value:e.serviceName},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])):H("",!0),I(Q).startStopFlag?(C(),U(v,{key:2,onClick:p[5]||(p[5]=e=>q(!1)),class:"right-btn",type:"success",plain:""},{default:F((()=>[x(u,{class:"el-input__icon"},{default:F((()=>[x(I(b))])),_:1}),V(" 开始 ")])),_:1})):(C(),U(v,{key:3,class:"right-btn",onClick:p[6]||(p[6]=e=>q(!0)),type:"danger",plain:""},{default:F((()=>[x(u,{class:"el-input__icon"},{default:F((()=>[x(I(y))])),_:1}),V(" 停止 ")])),_:1})),x(v,{class:"right-btn",onClick:p[7]||(p[7]=e=>{Q.dataTable=[]}),type:"warning",plain:""},{default:F((()=>[x(u,{class:"el-input__icon"},{default:F((()=>[x(I(w))])),_:1}),V(" 清空 ")])),_:1}),x(v,{class:"right-btn",onClick:p[8]||(p[8]=e=>(()=>{let e=[];Q.dataTable.forEach((a=>{e.push({timeStamp:a.timeStamp+"\t",direction:Q.dataType[a.direction],label:a.label,content:a.content})}));var a=h.unparse(e);let t=new Blob([a]),l=window.URL||window.webkitURL||window,o=l.createObjectURL(t),s=document.createElement("a");s.href=o;let c="";c=0===Q.logType?"采集日志":1===Q.logType?"上报日志":"系统日志",c+=f(new Date).format("YYYY-MM-DD HH-mm-ss"),s.download=c+".csv",s.click(),l.revokeObjectURL(o)})()),type:"primary",plain:""},{default:F((()=>[x(u,{class:"el-input__icon"},{default:F((()=>[x(I(k))])),_:1}),V(" 导出 ")])),_:1})])]),L("div",{class:"cTable",ref_key:"contentRef",ref:Y},[x(S,{data:I(se),"cell-style":I(Q).cellStyle,"header-cell-style":I(Q).headerCellStyle,"max-height":I(Q).tableMaxHeight,style:{width:"100%"},stripe:""},{empty:F((()=>[G])),default:F((()=>[x(_,{prop:"timeStamp",label:"时间戳",width:"auto","min-width":"220",align:"center"}),x(_,{label:"数据方向",width:"auto","min-width":"150",align:"center"},{default:F((e=>[x(j,{type:0===e.row.direction?"success":"warning"},{default:F((()=>[V(P(I(Q).dataType[e.row.direction]),1)])),_:2},1032,["type"])])),_:1}),x(_,{prop:"label",label:"数据标识",width:"auto","min-width":"150",align:"center"}),x(_,{prop:"content",label:"数据内容",width:"auto","min-width":"1000",align:"left"})])),_:1},8,["data","cell-style","header-cell-style","max-height"]),L("div",K,[x(T,{"current-page":I(Q).currentPage,"page-size":I(Q).pagesize,"page-sizes":[20,50,200,500],total:I(Q).dataTable.length,onCurrentChange:le,onSizeChange:oe,background:"",layout:"total, sizes, prev, pager, next, jumper",style:{"margin-top":"46px"}},null,8,["current-page","page-size","total"])])],512)])])}}},[["__scopeId","data-v-0b8f00d5"]]);export{Q as default};
