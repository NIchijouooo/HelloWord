import{g as e,e as a,p as l,n as t,o,m as r,A as n,C as s,D as i,B as c,q as p}from"./element-plus.3300f756.js";import{K as d,f as u,k as g,af as f,o as h,c as m,a as b,P as y,S as C,u as v,W as w,X as D,$ as _,a0 as x,n as P,aw as T,ax as z}from"./@vue.73d14ebc.js";import{_ as I,u as k,v as S}from"./index.9b2cfea2.js";import{I as V}from"./interface.4c14d44d.js";import{D as L}from"./deviceModel.16d94447.js";import{b as j,d as F,g as N}from"./@element-plus.0bef5fa5.js";const M={class:"main-container"},H={class:"main"},R={class:"search-bar"},E={class:"search-bar",style:{display:"flex"}},U={class:"title",style:{position:"relative","margin-right":"40px","justify-content":"flex-start",padding:"0px 0px",height:"40px"}},W={class:"tName"},A=(e=>(T("data-v-0fec3559"),e=e(),z(),e))((()=>b("div",null,"无数据",-1))),q={class:"pagination"},B={class:"dialog-content",style:{"min-height":"408px",overflow:"unset",padding:"0"}},K={class:"dialog-content-head"},O={class:"dialog-content-content"},X={class:"pagination dialog-pagination"},$={class:"dialog-footer"};var G=I({__name:"Device-property",props:{curDevice:{type:Object,default:{}},collInterfaceName:{type:String,default:""},pageInfo:{type:String,default:""}},emits:["changeIdFlag"],setup(T,{emit:z}){const I=T,G=k();console.log("id -> props",I);const J=d({headerCellStyle:{background:S.primaryColor,color:S.fontWhiteColor,height:"54px"},cellStyle:{height:"48px"},tableMaxHeight:0,currentPage:1,pagesize:20,propertyTableData:[],propertyInfo:"",pFlag:!1,pInfo:"",pCurrentPage:1,pPagesize:5,userHeadStyle:{color:S.headerTextColor,borderColor:"#C0C4CC",height:"48px"},userCellStyle:{borderColor:"#C0C4CC",height:"48px"},pTableData:[],pTableData1:[],selectPorperties:[],isLoading:!1}),Q=u(null),Y=a=>{const l={token:G.token,data:{collInterfaceName:I.curDevice.collInterfaceName,deviceName:I.curDevice.name}};J.isLoading=!0,V.getDeviceDataReal(l).then((async l=>{console.log("getDeviceDataReal -> res",l),l&&("0"===l.code?(J.propertyTableData=l.data,1===a&&e.success("刷新成功！")):ie(l),J.isLoading=!1,console.log("getDeviceDataReal -> ctxData.propertyTableData",J.propertyTableData),await P((()=>{let e=Q.value.clientHeight-34-36-42;I.pageInfo&&(e=Q.value.clientHeight-34-36-132),J.tableMaxHeight=e})))}))};Y();const Z=g((()=>{let e=J.propertyInfo;return J.propertyTableData.filter((a=>!e||a.name.toLowerCase().includes(e.toLowerCase())||a.label.toLowerCase().includes(e.toLowerCase()))).slice((J.currentPage-1)*J.pagesize,J.currentPage*J.pagesize)})),ee=g((()=>{let e=J.propertyInfo;return J.propertyTableData.filter((a=>!e||a.name.toLowerCase().includes(e.toLowerCase())||a.label.toLowerCase().includes(e.toLowerCase())))})),ae=e=>{J.currentPage=e},le=e=>{J.pagesize=e},te=()=>{J.pFlag=!1},oe=e=>{te()},re=e=>{J.pCurrentPage=e,J.pTableData=J.pTableData1.slice((J.pCurrentPage-1)*J.pPagesize,J.pCurrentPage*J.pPagesize)},ne=e=>{J.pPagesize=e,J.pTableData=J.pTableData1.slice((J.pCurrentPage-1)*J.pPagesize,J.pCurrentPage*J.pPagesize)},se=e=>{console.log("handleSelectionChange val = ",e),console.log("handleSelectionChange pTableData = ",J.pTableData),J.selectPorperties=e},ie=a=>{e({type:"error",message:a.message})};return(d,u)=>{const g=a,P=l,T=t,k=o,S=r,ce=f("Icon"),pe=n,de=s,ue=i,ge=c,fe=p;return h(),m("div",M,[b("div",H,[b("div",R,[y(k,{inline:!0,ref:"searchFormRef","status-icon":"","label-width":"120px"},{default:C((()=>[y(T,{style:{"margin-right":"20px","margin-left":"20px"}},{default:C((()=>[y(P,{type:"primary",plain:"",onClick:u[0]||(u[0]=e=>{z("changeIdFlag")})},{default:C((()=>[y(g,{class:"el-input__icon"},{default:C((()=>[y(v(j))])),_:1}),w(" 返回采集设备 ")])),_:1})])),_:1})])),_:1},512)]),b("div",E,[b("div",U,[b("div",W,D(I.curDevice.name)+"("+D(I.curDevice.label)+")",1)]),y(k,{inline:!0,ref:"searchFormRef","status-icon":"","label-width":"120px"},{default:C((()=>[y(T,null,{default:C((()=>[y(S,{style:{width:"200px"},placeholder:"请输入 名称/标签 过滤",clearable:"",modelValue:v(J).propertyInfo,"onUpdate:modelValue":u[1]||(u[1]=e=>v(J).propertyInfo=e)},{prefix:C((()=>[y(g,{class:"el-input__icon"},{default:C((()=>[y(v(F))])),_:1})])),_:1},8,["modelValue"])])),_:1}),y(T,null,{default:C((()=>[y(P,{type:"primary",plain:"",class:"right-btn",onClick:u[2]||(u[2]=e=>(()=>{const e={token:G.token,data:{name:I.curDevice.tsl}};L.getDeviceModelProperty(e).then((e=>{"0"===e.code?(J.pTableData=[],J.pTableData1=[],e.data.filter((e=>0!==e.accessMode)).forEach((e=>{const a={name:e.name,label:e.label,sendValue:"",result:{Code:"",Message:""}};J.pTableData.push(a),J.pTableData1.push(a)})),console.log("ctxData.pTableData -> ",J.pTableData),J.pFlag=!0):ie(e)}))})())},{default:C((()=>[y(g,{class:"el-input__icon"},{default:C((()=>[y(v(N))])),_:1}),w(" 写属性 ")])),_:1})])),_:1}),y(T,null,{default:C((()=>[y(P,{style:{color:"#fff","margin-left":"20px"},color:"#2EA554",class:"right-btn",onClick:u[3]||(u[3]=e=>{Y(1)})},{default:C((()=>[y(g,{class:"btn-icon"},{default:C((()=>[y(ce,{name:"local-refresh",size:"14px",color:"#ffffff"})])),_:1}),w(" 刷新 ")])),_:1})])),_:1})])),_:1},512)]),b("div",{class:"content",ref_key:"contentRef",ref:Q},[y(de,{data:v(Z),"cell-style":v(J).cellStyle,"header-cell-style":v(J).headerCellStyle,style:{width:"100%"},"max-height":v(J).tableMaxHeight,stripe:""},{empty:C((()=>[A])),default:C((()=>[y(pe,{prop:"index",label:"序号",width:"55"}),y(pe,{sortable:"",prop:"name",label:"变量名称",width:"auto","min-width":"180",align:"center"}),y(pe,{sortable:"",prop:"label",label:"变量标签",width:"auto","min-width":"180",align:"center"}),y(pe,{sortable:"",prop:"type",label:"变量类型",width:"auto","min-width":"120",align:"center"}),y(pe,{sortable:"",prop:"value",label:"变量值",width:"auto","min-width":"200",align:"center"}),y(pe,{sortable:"",prop:"unit",label:"单位",width:"auto","min-width":"200",align:"center"}),y(pe,{sortable:"",prop:"explain",label:"说明",width:"auto","min-width":"200",align:"center"}),y(pe,{sortable:"",prop:"timestamp",label:"实测时间",width:"auto","min-width":"220",align:"center"}),y(pe,{label:"操作",width:"auto","min-width":"200",align:"center",fixed:"right"},{default:C((a=>[y(P,{onClick:l=>(a.row,void e.info("功能完善中...")),text:"",type:"success"},{default:C((()=>[w("查看历史数据")])),_:2},1032,["onClick"])])),_:1})])),_:1},8,["data","cell-style","header-cell-style","max-height"]),b("div",q,[y(ue,{"current-page":v(J).currentPage,"page-size":v(J).pagesize,"page-sizes":[20,50,200,500],total:v(ee).length,onCurrentChange:ae,onSizeChange:le,background:"",layout:"total, sizes, prev, pager, next, jumper",style:{"margin-top":"46px"}},null,8,["current-page","page-size","total"])])],512),y(fe,{modelValue:v(J).pFlag,"onUpdate:modelValue":u[8]||(u[8]=e=>v(J).pFlag=e),title:"写属性",width:"1000px","before-close":oe,"close-on-click-modal":!1},{footer:C((()=>[b("span",$,[y(P,{onClick:u[7]||(u[7]=e=>te())},{default:C((()=>[w("关闭")])),_:1})])])),default:C((()=>[b("div",B,[b("div",K,[y(S,{placeholder:"请输入属性名称或标签",style:{width:"200px"},modelValue:v(J).pInfo,"onUpdate:modelValue":u[4]||(u[4]=e=>v(J).pInfo=e),onChange:u[5]||(u[5]=e=>{J.pTableData=J.pTableData1.filter((e=>!J.pInfo||e.name.includes(J.pInfo)||e.label.includes(J.pInfo)))})},{prefix:C((()=>[y(g,{class:"el-input__icon"},{default:C((()=>[y(v(F))])),_:1})])),_:1},8,["modelValue"]),y(P,{type:"primary",class:"right-btn",onClick:u[6]||(u[6]=a=>(a=>{if(console.log("sendCmd",J.selectPorperties),J.selectPorperties.length<1)return void e.warning("请至少选择一个属性操作");const l={};J.selectPorperties.forEach((e=>{l[e.name]=e.sendValue})),console.log(l);const t={token:G.token,data:{collInterfaceName:I.curDevice.collInterfaceName,deviceName:I.curDevice.name,serviceName:"SetVariables",serviceParam:l}};"Read"===a||V.invokeDeviceService(t).then((a=>{e.success("写属性命令下发成功！"),J.pTableData.forEach((e=>{J.selectPorperties.forEach((l=>{e.name===l.name&&(e.result=a)}))}))}))})("Write"))},{default:C((()=>[y(g,{class:"el-input__icon"},{default:C((()=>[y(v(N))])),_:1}),w(" 写属性 ")])),_:1})]),b("div",O,[y(de,{data:v(J).pTableData,onSelectionChange:se,"header-cell-style":v(J).userHeadStyle,"cell-style":v(J).userCellStyle,style:{border:"1px solid #c0c4cc"},"max-height":"290",stripe:""},{default:C((()=>[y(pe,{type:"selection",width:"55"}),y(pe,{sortable:"",prop:"name",label:"属性名称","min-width":"120",align:"center"}),y(pe,{sortable:"",prop:"label",label:"属性标签","min-width":"140",align:"center"}),y(pe,{sortable:"",label:"写入值","min-width":"100",align:"center"},{default:C((e=>[y(S,{placeholder:"请输入写入值",modelValue:e.row.sendValue,"onUpdate:modelValue":a=>e.row.sendValue=a},null,8,["modelValue","onUpdate:modelValue"])])),_:1}),y(pe,{sortable:"",label:"返回结果","min-width":"250",align:"center"},{default:C((e=>[_(y(ge,{type:"success"},{default:C((()=>[w(" 操作成功 ")])),_:2},1536),[[x,"0"===e.row.result.Code]]),_(y(ge,{type:"danger"},{default:C((()=>[w(" 操作失败 ")])),_:2},1536),[[x,"1"===e.row.result.Code]]),_(b("span",null,"-",512),[[x,""===e.row.result.Code]]),_(b("span",null,D(e.row.result.Message),513),[[x,"1"===e.row.result.Code]])])),_:1})])),_:1},8,["data","header-cell-style","cell-style"]),b("div",X,[y(ue,{"current-page":v(J).pCurrentPage,"page-size":v(J).pPagesize,"page-sizes":[5,10,20,50],total:v(J).pTableData1.length,onCurrentChange:re,onSizeChange:ne,background:"",layout:"total, sizes, prev, pager, next, jumper"},null,8,["current-page","page-size","total"])])])])])),_:1},8,["modelValue"])])])}}},[["__scopeId","data-v-0fec3559"]]);export{G as D};
