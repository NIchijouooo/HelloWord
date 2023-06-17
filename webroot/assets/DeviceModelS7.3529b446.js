import{g as e,L as l,e as a,p as t,n as o,o as r,m as p,A as d,C as s,D as i,w as n,K as m,q as c,M as u,v as y}from"./element-plus.2c9426ce.js";import{_ as g,u as f,v as b}from"./index.d5c032bc.js";import{D as h}from"./deviceModel.d8833265.js";import{r as v,J as F,j as x,ae as _,o as M,c as w,a as V,O as k,R as C,u as L,V as A,W as D,P,aj as T,_ as j,$ as z,S as U,X as S,n as N,av as R,aw as q}from"./@vue.2c474d7f.js";import{b as I,h as O,i as H,d as B}from"./@element-plus.e895b9e1.js";import"./lodash-es.98e49362.js";import"./async-validator.fb49d0f5.js";import"./@vueuse.d618444d.js";import"./dayjs.2790329b.js";import"./axios.765908e4.js";import"./@ctrl.b082b0c1.js";import"./@popperjs.36402333.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.8aeb3683.js";import"./vue-router.3b7aa6c5.js";import"./pinia.39330980.js";import"./vue-demi.b3a9cad9.js";import"./screenfull.7cf96174.js";import"./vue-echarts.ebbe6a42.js";import"./resize-detector.4e96b72b.js";import"./echarts.2b8dcc79.js";import"./zrender.54f191d6.js";import"./tslib.60310f1a.js";import"./nprogress.aae2b787.js";const E=e=>(R("data-v-6bd318ce"),e=e(),q(),e),W={class:"main-container"},$={class:"main"},J={class:"search-bar"},K={class:"search-bar",style:{display:"flex"}},X={class:"title",style:{position:"relative","margin-right":"40px","justify-content":"flex-start",padding:"0px 0px",height:"40px"}},G={class:"tName"},Q={class:"param-content"},Y={class:"pc-title"},Z={class:"pct-info"},ee={class:"pc-content"},le={class:"param-value"},ae={key:0,class:"param-name"},te={key:1,class:"param-name"},oe=E((()=>V("div",null,"无数据",-1))),re={class:"pagination"},pe={class:"dialog-content"},de=E((()=>V("div",{class:"form-title"},[V("div",{class:"tName"},"配置参数")],-1))),se={class:"dialog-footer"},ie=E((()=>V("div",{class:"el-upload__tip"},"只能上传一个文件，只支持xlsx格式文件！",-1))),ne={class:"dialog-footer"};var me=g({__name:"PropertyS7",props:{curDeviceModel:{type:Object,default:{}}},emits:["changeDpFlag"],setup(g,{emit:R}){const q=g,E=f();console.log("id -> props",q);const me=/^[0-9]+(.[0-9]{1,2})?$/,ce=(e,l,a)=>{me.test(l)?a():a(new Error("只能输入大于等于0,最多两位小数的数字！"))},ue=v(null),ye=F({headerCellStyle:{background:b.primaryColor,color:b.fontWhiteColor,height:"54px"},cellStyle:{height:"48px"},tableMaxHeight:0,currentPage:1,pagesize:20,propertyList:[],deviceModelProperty:"",accessModeOptions:[{label:"只读",value:0},{label:"只写",value:1},{label:"读写",value:2}],typeOptions:[{label:"uint32",value:0},{label:"int32",value:1},{label:"double",value:2},{label:"string",value:3}],dataTypeOptions:[{label:"uint8",value:0},{label:"int8",value:1},{label:"uint16",value:2},{label:"int16",value:3},{label:"uint32",value:4},{label:"int32",value:5},{label:"float",value:6},{label:"double",value:7},{label:"bool",value:8}],typeNames:{t0:"uint32",t1:"int32",t2:"double",t3:"string"},dataTypeNames:{dt0:"uint8",dt1:"int8",dt2:"uint16",dt3:"int16",dt4:"uint32",dt5:"int32",dt6:"float",dt7:"double",dt8:"bool"},accessModeNames:{am0:"只读",am1:"只写",am2:"读写"},pFlag:!1,pTitle:"添加属性",propertyForm:{name:"",label:"",accessMode:0,type:0,dbNumber:"",startAddr:"",dataType:0,unit:"",decimals:0,min:0,max:0,minMaxAlarm:!1,step:0,stepAlarm:!1,dataLength:0,dataLengthAlarm:!1},paramName:{name:"属性名称",label:"属性标签",accessMode:"读写模式",type:"属性类型",decimals:"小数位数",dbNumber:"数据块号",dataType:"数据类型",startAddr:"PLC地址",min:"最小值",max:"最大值",minMaxAlarm:"范围报警",step:"步长",stepAlarm:"步长报警",dataLength:"字符串长度",dataLengthAlarm:"字符串长度报警"},propertyRules:{name:[{required:!0,message:"属性名称不能为空",trigger:"blur"}],accessMode:[{required:!0,message:"读写属性不能为空",trigger:"blur"}],type:[{required:!0,message:"属性类型不能为空",trigger:"blur"}],dataType:[{required:!0,message:"数据类型不能为空",trigger:"blur"}],decimals:[{type:"number",message:"小数位数只能输入数字"}],dbNumber:[{required:!0,message:"数据块号不能为空",trigger:"blur"}],startAddr:[{required:!0,message:"PLC地址不能为空",trigger:"blur"}],step:[{required:!0,message:"步长不能为空",trigger:"blur"},{trigger:"blur",validator:ce}],min:[{required:!0,message:"最小值不能为空",trigger:"blur"},{trigger:"blur",validator:ce}],max:[{required:!0,message:"最大值不能为空",trigger:"blur"},{trigger:"blur",validator:ce}],dataLength:[{required:!0,message:"字符串长度不能为空",trigger:"blur"},{trigger:"blur",validator:/^[0-9]*[1-9][0-9]*$/}]},psFlag:!1,selectedProperties:[]}),ge=l=>{const a={token:E.token,data:{name:q.curDeviceModel.name}};h.getDeviceModelProperty(a).then((async a=>{a&&("0"===a.code?(ye.propertyList=a.data,1===l&&e({type:"success",message:"刷新成功！"})):je(a),await N((()=>{ye.tableMaxHeight=ue.value.clientHeight-34-36-132})))}))};ge();const fe=x((()=>(console.log("ctxData.propertyList ->",ye.propertyList),ye.propertyList.filter((e=>{var l=!ye.deviceModelProperty,a=e.name.toLowerCase().includes(ye.deviceModelProperty.toLowerCase());return l||a})).slice((ye.currentPage-1)*ye.pagesize,ye.currentPage*ye.pagesize)))),be=x((()=>ye.propertyList.filter((e=>{var l=!ye.deviceModelProperty,a=e.name.toLowerCase().includes(ye.deviceModelProperty.toLowerCase());return l||a})))),he=e=>{ye.currentPage=e},ve=e=>{ye.pagesize=e},Fe=()=>{_e.value&&_e.value.submit()},xe=e=>{},_e=v(null),Me=e=>{_e.value.clearFiles();const l=e[0];_e.value.handleStart(l)},we=l=>{console.log("beforeUpload -> file",l);let a=["application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"].includes(l.type);if(a){if(a){let a=new FormData;a.append("name",q.curDeviceModel.name),a.append("fileName",l);const t={token:E.token,contentType:"multipart/form-data",data:a};h.importDeviceModelProptyAndService(t).then((l=>{"0"===l.code?(e.success(l.message),ye.psFlag=!1,ge()):je(l)}))}}else e({type:"error",message:"文件格式不正确,必须是xlsx文件！"})},Ve=()=>{ye.psFlag=!1,uploadRef.value.clearFiles()},ke=()=>{Ve()},Ce=e=>{ye.pFlag=!0,ye.pTitle="编辑属性",ye.propertyForm.name=e.name,ye.propertyForm.label=e.label,ye.propertyForm.accessMode=e.accessMode,ye.propertyForm.type=e.type,ye.propertyForm.decimals=e.decimals,ye.propertyForm.unit=e.unit,3!==e.type?(ye.propertyForm.min=e.params.min,ye.propertyForm.max=e.params.max,ye.propertyForm.minMaxAlarm=e.params.minMaxAlarm,ye.propertyForm.step=e.params.step,ye.propertyForm.stepAlarm=e.params.stepAlarm):(ye.propertyForm.dataLength=e.params.dataLength,ye.propertyForm.dataLengthAlarm=e.params.dataLengthAlarm),ye.propertyForm.dbNumber=e.params.dbNumber,ye.propertyForm.startAddr=e.params.startAddr,ye.propertyForm.dataType=e.params.dataType,console.log("ctxData.propertyForm",ye.propertyForm)},Le=v(null),Ae=()=>{ye.pFlag=!1,Le.value.resetFields(),Te()},De=e=>{Ae()},Pe=e=>{ye.selectedProperties=e,console.log("handleSelectionChange -> val =",e)},Te=()=>{ye.propertyForm={name:"",label:"",accessMode:0,type:0,decimals:0,unit:"",dbNumber:"",startAddr:"",dataType:0,step:0,min:0,max:0,dataLength:0}},je=l=>{e({type:"error",message:l.message})},ze=(l,a)=>{e({type:"0"===l.code?"success":"error",message:l.message}),"0"===l.code&&a&&a()};return(g,f)=>{const b=a,v=t,F=o,x=r,N=p,me=_("Icon"),ce=d,Te=s,Ue=i,Se=y,Ne=n,Re=m,qe=c,Ie=u;return M(),w("div",W,[V("div",$,[V("div",J,[k(x,{inline:!0,ref:"searchFormRef","status-icon":"","label-width":"90px"},{default:C((()=>[k(F,{style:{"margin-left":"20px"}},{default:C((()=>[k(v,{type:"primary",plain:"",onClick:f[0]||(f[0]=e=>(console.log("toDeviceModel"),void R("changeDpFlag")))},{default:C((()=>[k(b,{class:"el-input__icon"},{default:C((()=>[k(L(I))])),_:1}),A(" 返回设备模型 ")])),_:1})])),_:1}),k(F,null,{default:C((()=>[k(v,{type:"primary",plain:"",class:"right-btn",onClick:f[1]||(f[1]=e=>(console.log("importDPS"),void(ye.psFlag=!0)))},{default:C((()=>[k(b,{class:"el-input__icon"},{default:C((()=>[k(L(O))])),_:1}),A(" 导入模型属性 ")])),_:1})])),_:1}),k(F,null,{default:C((()=>[k(v,{type:"primary",plain:"",class:"right-btn",onClick:f[2]||(f[2]=e=>(()=>{console.log("exportDPS");const e={token:E.token,responseType:"blob",data:{name:q.curDeviceModel.name}};h.exportDeviceModelProptyAndService(e).then((e=>{if(console.log("exportDeviceModelProptyAndService -> res",e),e&&"1"===e.code)return void je(e);const l=new Blob([e.blob]),a=e.fileName,t=document.createElement("a");t.download=a,t.style.display="none",t.href=URL.createObjectURL(l),document.body.appendChild(t),t.click(),URL.revokeObjectURL(t.href),document.body.removeChild(t)}))})())},{default:C((()=>[k(b,{class:"el-input__icon"},{default:C((()=>[k(L(H))])),_:1}),A(" 导出模型属性 ")])),_:1})])),_:1})])),_:1},512)]),V("div",K,[V("div",X,[V("div",G,D(q.curDeviceModel.label)+"：属性列表",1)]),k(x,{inline:!0,ref:"searchFormRef2","status-icon":"","label-width":"90px"},{default:C((()=>[k(F,{label:""},{default:C((()=>[k(N,{style:{width:"200px"},placeholder:"请输入属性名称",modelValue:L(ye).deviceModelProperty,"onUpdate:modelValue":f[3]||(f[3]=e=>L(ye).deviceModelProperty=e)},{prefix:C((()=>[k(b,{class:"el-input__icon"},{default:C((()=>[k(L(B))])),_:1})])),_:1},8,["modelValue"])])),_:1}),k(F,null,{default:C((()=>[k(v,{type:"primary",bg:"",class:"right-btn",onClick:f[4]||(f[4]=e=>(console.log("addDeviceModelProperty"),ye.pFlag=!0,void(ye.pTitle="添加属性")))},{default:C((()=>[k(b,{class:"btn-icon"},{default:C((()=>[k(me,{name:"local-add",size:"14px",color:"#ffffff"})])),_:1}),A(" 添加 ")])),_:1})])),_:1}),k(F,null,{default:C((()=>[k(v,{style:{color:"#fff"},color:"#2EA554",class:"right-btn",onClick:f[5]||(f[5]=e=>{ge(1)})},{default:C((()=>[k(b,{class:"btn-icon"},{default:C((()=>[k(me,{name:"local-refresh",size:"14px",color:"#ffffff"})])),_:1}),A(" 刷新 ")])),_:1})])),_:1}),k(F,null,{default:C((()=>[k(v,{type:"danger",bg:"",class:"right-btn",onClick:f[6]||(f[6]=a=>(a=>{let t=[];0!==ye.selectedProperties.length?(ye.selectedProperties.forEach((e=>{t.push(e.name)})),console.log("pList",t),l.confirm("确认要删除这些属性吗?","警告",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then((()=>{const e={token:E.token,data:{name:q.curDeviceModel.name,properties:t}};h.deleteDeviceModelProperty(e).then((e=>{ze(e,ge)}))})).catch((()=>{e({type:"info",message:"取消删除"})}))):e.info("请至少选择一个属性！")})())},{default:C((()=>[k(b,{class:"btn-icon"},{default:C((()=>[k(me,{name:"local-delete",size:"14px",color:"#ffffff"})])),_:1}),A(" 删除 ")])),_:1})])),_:1})])),_:1},512)]),V("div",{class:"content",ref_key:"contentRef",ref:ue},[k(Te,{data:L(fe),"cell-style":L(ye).cellStyle,"header-cell-style":L(ye).headerCellStyle,"max-height":L(ye).tableMaxHeight,style:{width:"100%"},stripe:"",onSelectionChange:Pe,onRowDblclick:Ce},{empty:C((()=>[oe])),default:C((()=>[k(ce,{type:"selection",width:"55"}),k(ce,{type:"expand"},{default:C((e=>[V("div",Q,[V("div",Y,[V("div",Z,[V("b",null,D(e.row.name),1),A(" "+D("参数详情"))])]),V("div",ee,[(M(!0),w(P,null,T(e.row.params,((e,l,a)=>(M(),w("div",{class:"param-item",key:a},[V("div",le,D(L(ye).paramName[l])+"：",1),"dataType"===l?(M(),w("div",ae,D(L(ye).dataTypeNames["dt"+e]),1)):S("",!0),"dataType"!==l?(M(),w("div",te,D("boolean"==typeof e?e?"是":"否":e),1)):S("",!0)])))),128))])])])),_:1}),k(ce,{sortable:"",prop:"name",label:"属性名称",width:"auto","min-width":"150",align:"center"}),k(ce,{sortable:"",prop:"label",label:"属性标签",width:"auto","min-width":"150",align:"center"}),k(ce,{sortable:"",label:"读写属性",width:"auto","min-width":"80",align:"center"},{default:C((e=>[A(D(L(ye).accessModeNames["am"+e.row.accessMode]),1)])),_:1}),k(ce,{sortable:"",prop:"type",label:"属性类型",width:"auto","min-width":"100",align:"center"},{default:C((e=>[A(D(L(ye).typeNames["t"+e.row.type]),1)])),_:1}),k(ce,{sortable:"",label:"小数位数",width:"auto","min-width":"80",align:"center"},{default:C((e=>[A(D(""===e.row.decimals?0:e.row.decimals),1)])),_:1}),k(ce,{sortable:"",prop:"unit",label:"单位",width:"auto","min-width":"80",align:"center"}),k(ce,{label:"操作",width:"auto","min-width":"200",align:"center",fixed:"right"},{default:C((e=>[k(v,{onClick:l=>Ce(e.row),text:"",type:"primary"},{default:C((()=>[A("编辑")])),_:2},1032,["onClick"])])),_:1})])),_:1},8,["data","cell-style","header-cell-style","max-height"]),V("div",re,[k(Ue,{"current-page":L(ye).currentPage,"page-size":L(ye).pagesize,"page-sizes":[20,50,200,500],total:L(be).length,onCurrentChange:he,onSizeChange:ve,background:"",layout:"total, sizes, prev, pager, next, jumper",style:{"margin-top":"46px"}},null,8,["current-page","page-size","total"])])],512),k(qe,{modelValue:L(ye).pFlag,"onUpdate:modelValue":f[25]||(f[25]=e=>L(ye).pFlag=e),title:L(ye).pTitle,width:"800px","before-close":De,"close-on-click-modal":!1},{footer:C((()=>[V("span",se,[k(v,{onClick:f[23]||(f[23]=e=>Ae())},{default:C((()=>[A("取消")])),_:1}),k(v,{type:"primary",onClick:f[24]||(f[24]=l=>{3!==ye.propertyForm.type&&Number(ye.propertyForm.min)>Number(ye.propertyForm.max)?e.warning("最大值必须大于最小值"):Le.value.validate((e=>{if(console.log("valid",e),!e)return!1;{let e={};e.name=ye.propertyForm.name,e.label=ye.propertyForm.label,e.accessMode=ye.propertyForm.accessMode,e.type=ye.propertyForm.type,e.decimals=ye.propertyForm.decimals,e.unit=ye.propertyForm.unit;let l={dbNumber:ye.propertyForm.dbNumber,startAddr:ye.propertyForm.startAddr,dataType:ye.propertyForm.dataType};3!==ye.propertyForm.type?(l.min=ye.propertyForm.min.toString(),l.max=ye.propertyForm.max.toString(),l.minMaxAlarm=ye.propertyForm.minMaxAlarm,l.step=ye.propertyForm.step.toString(),l.stepAlarm=ye.propertyForm.stepAlarm):(l.dataLength=ye.propertyForm.dataLength.toString(),l.dataLengthAlarm=ye.propertyForm.dataLengthAlarm),e.params=l,console.log("submitPorpertyForm -> property",e);const a={token:E.token,data:{name:q.curDeviceModel.name,property:e}};ye.pTitle.includes("添加")?(console.log("添加属性"),h.addDeviceModelProperty(a).then((e=>{ze(e,ge),Ae()}))):(console.log("编辑属性"),h.editDeviceModelProperty(a).then((e=>{ze(e,ge),Ae()})))}}))})},{default:C((()=>[A("保存")])),_:1})])])),default:C((()=>[V("div",pe,[k(x,{model:L(ye).propertyForm,rules:L(ye).propertyRules,ref_key:"propertyFormRef",ref:Le,"status-icon":"","label-position":"right","label-width":"120px",inline:"true"},{default:C((()=>[k(F,{label:"属性名称",prop:"name"},{default:C((()=>[k(N,{style:{width:"220px"},disabled:L(ye).pTitle.includes("编辑"),type:"text",modelValue:L(ye).propertyForm.name,"onUpdate:modelValue":f[7]||(f[7]=e=>L(ye).propertyForm.name=e),autocomplete:"off",placeholder:"请输入属性名称"},null,8,["disabled","modelValue"])])),_:1}),k(F,{label:"属性标签",prop:"label"},{default:C((()=>[k(N,{style:{width:"220px"},type:"text",modelValue:L(ye).propertyForm.label,"onUpdate:modelValue":f[8]||(f[8]=e=>L(ye).propertyForm.label=e),autocomplete:"off",placeholder:"请输入属性标签"},null,8,["modelValue"])])),_:1}),k(F,{label:"读写属性",prop:"accessMode"},{default:C((()=>[k(Ne,{modelValue:L(ye).propertyForm.accessMode,"onUpdate:modelValue":f[9]||(f[9]=e=>L(ye).propertyForm.accessMode=e),style:{width:"220px"},placeholder:"请选择读写属性"},{default:C((()=>[(M(!0),w(P,null,T(L(ye).accessModeOptions,(e=>(M(),U(Se,{key:"accessMode_"+e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),k(F,{label:"单位",prop:"unit"},{default:C((()=>[k(N,{style:{width:"220px"},type:"text",modelValue:L(ye).propertyForm.unit,"onUpdate:modelValue":f[10]||(f[10]=e=>L(ye).propertyForm.unit=e),autocomplete:"off",placeholder:"请输入单位"},null,8,["modelValue"])])),_:1}),k(F,{label:"属性类型",prop:"type"},{default:C((()=>[k(Ne,{modelValue:L(ye).propertyForm.type,"onUpdate:modelValue":f[11]||(f[11]=e=>L(ye).propertyForm.type=e),modelModifiers:{number:!0},style:{width:"220px"},placeholder:"请选择属性类型"},{default:C((()=>[(M(!0),w(P,null,T(L(ye).typeOptions,(e=>(M(),U(Se,{key:"type_"+e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1}),j(k(F,{label:"小数位数",prop:"decimals"},{default:C((()=>[k(N,{style:{width:"220px"},type:"text",modelValue:L(ye).propertyForm.decimals,"onUpdate:modelValue":f[12]||(f[12]=e=>L(ye).propertyForm.decimals=e),modelModifiers:{number:!0},autocomplete:"off",placeholder:"请输入小数位数"},null,8,["modelValue"])])),_:1},512),[[z,2===L(ye).propertyForm.type]]),V("div",null,[3!==L(ye).propertyForm.type?(M(),U(F,{key:0,label:"范围报警",prop:"minMaxAlarm"},{default:C((()=>[k(Re,{modelValue:L(ye).propertyForm.minMaxAlarm,"onUpdate:modelValue":f[13]||(f[13]=e=>L(ye).propertyForm.minMaxAlarm=e),"inline-prompt":"","active-text":"是","inactive-text":"否"},null,8,["modelValue"])])),_:1})):S("",!0)]),V("div",null,[3!==L(ye).propertyForm.type&&L(ye).propertyForm.minMaxAlarm?(M(),U(F,{key:0,label:"最小值",prop:"min"},{default:C((()=>[k(N,{type:"text",modelValue:L(ye).propertyForm.min,"onUpdate:modelValue":f[14]||(f[14]=e=>L(ye).propertyForm.min=e),autocomplete:"off",placeholder:"请输入最小值"},null,8,["modelValue"])])),_:1})):S("",!0)]),3!==L(ye).propertyForm.type&&L(ye).propertyForm.minMaxAlarm?(M(),U(F,{key:0,label:"最大值",prop:"max"},{default:C((()=>[k(N,{type:"text",modelValue:L(ye).propertyForm.max,"onUpdate:modelValue":f[15]||(f[15]=e=>L(ye).propertyForm.max=e),autocomplete:"off",placeholder:"请输入最大值"},null,8,["modelValue"])])),_:1})):S("",!0),V("div",null,[3!==L(ye).propertyForm.type?(M(),U(F,{key:0,label:"步长报警",prop:"stepAlarm"},{default:C((()=>[k(Re,{modelValue:L(ye).propertyForm.stepAlarm,"onUpdate:modelValue":f[16]||(f[16]=e=>L(ye).propertyForm.stepAlarm=e),"inline-prompt":"","active-text":"是","inactive-text":"否"},null,8,["modelValue"])])),_:1})):S("",!0)]),3!==L(ye).propertyForm.type&&L(ye).propertyForm.stepAlarm?(M(),U(F,{key:1,label:"步长",prop:"step"},{default:C((()=>[k(N,{type:"text",modelValue:L(ye).propertyForm.step,"onUpdate:modelValue":f[17]||(f[17]=e=>L(ye).propertyForm.step=e),autocomplete:"off",placeholder:"请输入步长"},null,8,["modelValue"])])),_:1})):S("",!0),V("div",null,[3===L(ye).propertyForm.type?(M(),U(F,{key:0,label:"字符串长度报警",prop:"dataLengthAlarm"},{default:C((()=>[k(Re,{modelValue:L(ye).propertyForm.dataLengthAlarm,"onUpdate:modelValue":f[18]||(f[18]=e=>L(ye).propertyForm.dataLengthAlarm=e),"inline-prompt":"","active-text":"是","inactive-text":"否"},null,8,["modelValue"])])),_:1})):S("",!0)]),3===L(ye).propertyForm.type&&L(ye).propertyForm.dataLengthAlarm?(M(),U(F,{key:2,label:"字符串长度",prop:"dataLength"},{default:C((()=>[k(N,{type:"text",modelValue:L(ye).propertyForm.dataLength,"onUpdate:modelValue":f[19]||(f[19]=e=>L(ye).propertyForm.dataLength=e),autocomplete:"off",placeholder:"请输入字符串长度"},null,8,["modelValue"])])),_:1})):S("",!0),de,k(F,{label:"数据块号",prop:"dbNumber"},{default:C((()=>[k(N,{type:"text",style:{width:"220px"},modelValue:L(ye).propertyForm.dbNumber,"onUpdate:modelValue":f[20]||(f[20]=e=>L(ye).propertyForm.dbNumber=e),autocomplete:"off",placeholder:"请输入数据块号"},null,8,["modelValue"])])),_:1}),k(F,{label:"",prop:"",style:{width:"220px"}}),k(F,{label:"PLC地址",prop:"startAddr"},{default:C((()=>[k(N,{type:"text",style:{width:"220px"},modelValue:L(ye).propertyForm.startAddr,"onUpdate:modelValue":f[21]||(f[21]=e=>L(ye).propertyForm.startAddr=e),autocomplete:"off",placeholder:"请输入PLC地址"},null,8,["modelValue"])])),_:1}),k(F,{label:"数据类型",prop:"dataType"},{default:C((()=>[k(Ne,{modelValue:L(ye).propertyForm.dataType,"onUpdate:modelValue":f[22]||(f[22]=e=>L(ye).propertyForm.dataType=e),modelModifiers:{number:!0},style:{width:"220px"},placeholder:"请选择数据类型"},{default:C((()=>[(M(!0),w(P,null,T(L(ye).dataTypeOptions,(e=>(M(),U(Se,{key:"type_"+e.value,label:e.label,value:e.value},null,8,["label","value"])))),128))])),_:1},8,["modelValue"])])),_:1})])),_:1},8,["model","rules"])])])),_:1},8,["modelValue","title"]),k(qe,{modelValue:L(ye).psFlag,"onUpdate:modelValue":f[26]||(f[26]=e=>L(ye).psFlag=e),title:"导入模型属性",width:"600px","before-close":ke,"close-on-click-modal":!1},{footer:C((()=>[V("span",ne,[k(v,{onClick:Ve},{default:C((()=>[A("取消")])),_:1}),k(v,{type:"primary",onClick:Fe},{default:C((()=>[A("上传")])),_:1})])])),default:C((()=>[k(Ie,{ref_key:"uploadDPSRef",ref:_e,action:"","auto-upload":!1,"http-request":xe,limit:1,"on-exceed":Me,"before-upload":we},{tip:C((()=>[ie])),default:C((()=>[k(v,{type:"primary"},{default:C((()=>[A("选择文件")])),_:1})])),_:1},512)])),_:1},8,["modelValue"])])])}}},[["__scopeId","data-v-6bd318ce"]]);const ce={class:"main-container"},ue={key:0,class:"main"},ye={class:"search-bar"},ge={class:"tool-bar"},fe=(e=>(R("data-v-189462a2"),e=e(),q(),e))((()=>V("div",null,"无数据",-1))),be={class:"pagination"},he={class:"dialog-content"},ve={class:"dialog-footer"};var Fe=g({__name:"DeviceModelS7",setup(n){const m=f(),u=v(null),y=F({modelType:1,dpFlag:!0,headerCellStyle:{background:b.primaryColor,color:b.fontWhiteColor,height:"54px"},cellStyle:{height:"48px"},tableMaxHeight:0,currentPage:1,pagesize:20,deviceModelInfo:"",tableData:[],modelForm:{name:"",label:"",type:0,param:""},modelRules:{name:[{required:!0,message:"采集模型名称不能为空",trigger:"blur"}],label:[{required:!0,message:"采集模型标识不能为空",trigger:"blur"}]},udFlag:!1,dmFlag:!1,dmTitle:"添加采集模型",pluginList:[],curDeviceModel:""}),g=l=>{const a={token:m.token,data:{type:y.modelType}};h.getDeviceModelList(a).then((async a=>{console.log("getDeviceModelList -> res = ",a),a&&("0"===a.code?(y.tableData=a.data,1===l&&e({type:"success",message:"刷新成功！"})):O(a),await N((()=>{y.tableMaxHeight=u.value.clientHeight-34-36-22})))}))};g();const D=x((()=>y.tableData.filter((e=>{var l=!y.deviceModelInfo,a=e.name.toLowerCase().includes(y.deviceModelInfo.toLowerCase())||e.label.toLowerCase().includes(y.deviceModelInfo.toLowerCase());return l||a})).slice((y.currentPage-1)*y.pagesize,y.currentPage*y.pagesize))),P=x((()=>y.tableData.filter((e=>{var l=!y.deviceModelInfo,a=e.name.toLowerCase().includes(y.deviceModelInfo.toLowerCase())||e.label.toLowerCase().includes(y.deviceModelInfo.toLowerCase());return l||a})))),T=e=>{y.currentPage=e},j=e=>{y.pagesize=e},z=e=>{y.dmFlag=!0,y.dmTitle="编辑采集模型",y.modelForm={name:e.name,label:e.label,type:void 0===e.type?1:e.type,param:""}},S=v(null),R=()=>{y.dmFlag=!1,S.value.resetFields(),I()},q=()=>{R()},I=()=>{y.modelForm={name:"",label:"",type:0}},O=l=>{e({type:"error",message:l.message})},H=(l,a)=>{e({type:"0"===l.code?"success":"error",message:l.message}),"0"===l.code&&a&&a()};return(n,f)=>{const b=a,v=p,F=o,x=_("Icon"),N=t,I=r,O=d,E=s,W=i,$=c;return M(),w("div",ce,[L(y).dpFlag?(M(),w("div",ue,[V("div",ye,[k(I,{inline:!0,ref:"searchFormRef","status-icon":"","label-width":"120px"},{default:C((()=>[k(F,{label:"采集模型名称"},{default:C((()=>[k(v,{style:{width:"200px"},placeholder:"请输入采集模型名称",modelValue:L(y).deviceModelInfo,"onUpdate:modelValue":f[0]||(f[0]=e=>L(y).deviceModelInfo=e)},{prefix:C((()=>[k(b,{class:"el-input__icon"},{default:C((()=>[k(L(B))])),_:1})])),_:1},8,["modelValue"])])),_:1}),k(F,null,{default:C((()=>[k(N,{style:{color:"#fff","margin-left":"20px"},color:"#2EA554",class:"right-btn",onClick:f[1]||(f[1]=e=>{g(1)})},{default:C((()=>[k(b,{class:"btn-icon"},{default:C((()=>[k(x,{name:"local-refresh",size:"14px",color:"#ffffff"})])),_:1}),A(" 刷新 ")])),_:1})])),_:1})])),_:1},512)]),V("div",ge,[k(N,{type:"primary",bg:"",class:"right-btn",onClick:f[2]||(f[2]=e=>(y.dmFlag=!0,void(y.dmTitle="添加采集模型")))},{default:C((()=>[k(b,{class:"btn-icon"},{default:C((()=>[k(x,{name:"local-add",size:"14px",color:"#ffffff"})])),_:1}),A(" 添加 ")])),_:1})]),V("div",{class:"content",ref_key:"contentRef",ref:u},[k(E,{data:L(D),"cell-style":L(y).cellStyle,"header-cell-style":L(y).headerCellStyle,"max-height":L(y).tableMaxHeight,style:{width:"100%"},stripe:"",onRowDblclick:z},{empty:C((()=>[fe])),default:C((()=>[k(O,{sortable:"",prop:"name",label:"采集模型名称",width:"auto","min-width":"200",align:"center"}),k(O,{sortable:"",prop:"label",label:"采集模型标签",width:"auto","min-width":"200",align:"center"}),k(O,{label:"操作",width:"auto","min-width":"300",align:"center",fixed:"right"},{default:C((a=>[k(N,{onClick:e=>{return l=a.row,y.curDeviceModel=l,void(y.dpFlag=!1);var l},text:"",type:"success"},{default:C((()=>[A("变量详情")])),_:2},1032,["onClick"]),k(N,{onClick:e=>z(a.row),text:"",type:"primary"},{default:C((()=>[A("编辑")])),_:2},1032,["onClick"]),k(N,{onClick:t=>{return o=a.row,void l.confirm("确定要删除这个采集模型吗?","警告",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then((()=>{const e={token:m.token,data:{name:o.name,type:1}};h.deleteDeviceModel(e).then((e=>{H(e,g)}))})).catch((()=>{e({type:"info",message:"取消删除"})}));var o},text:"",type:"danger"},{default:C((()=>[A("删除")])),_:2},1032,["onClick"])])),_:1})])),_:1},8,["data","cell-style","header-cell-style","max-height"]),V("div",be,[k(W,{"current-page":L(y).currentPage,"page-size":L(y).pagesize,"page-sizes":[20,50,200,500],total:L(P).length,onCurrentChange:T,onSizeChange:j,background:"",layout:"total, sizes, prev, pager, next, jumper",style:{"margin-top":"46px"}},null,8,["current-page","page-size","total"])])],512)])):(M(),U(me,{key:1,curDeviceModel:L(y).curDeviceModel,onChangeDpFlag:f[3]||(f[3]=e=>{y.dpFlag=!0}),style:{width:"100%",height:"100%",overflow:"hidden"}},null,8,["curDeviceModel"])),k($,{modelValue:L(y).dmFlag,"onUpdate:modelValue":f[8]||(f[8]=e=>L(y).dmFlag=e),title:L(y).dmTitle,width:"600px","before-close":q,"close-on-click-modal":!1},{footer:C((()=>[V("span",ve,[k(N,{onClick:f[6]||(f[6]=e=>R())},{default:C((()=>[A("取消")])),_:1}),k(N,{type:"primary",onClick:f[7]||(f[7]=e=>{S.value.validate((e=>{if(!e)return!1;{y.modelForm.type=y.modelType;const e={token:m.token,data:y.modelForm};y.dmTitle.includes("添加")?h.addDeviceModel(e).then((e=>{H(e,g),R()})):h.editDeviceModel(e).then((e=>{H(e,g),R()}))}}))})},{default:C((()=>[A("保存")])),_:1})])])),default:C((()=>[V("div",he,[k(I,{model:L(y).modelForm,rules:L(y).modelRules,ref_key:"modelFormRef",ref:S,"status-icon":"","label-position":"right","label-width":"120px"},{default:C((()=>[k(F,{label:"采集模型名称",prop:"name"},{default:C((()=>[k(v,{type:"text",disabled:L(y).dmTitle.includes("编辑"),modelValue:L(y).modelForm.name,"onUpdate:modelValue":f[4]||(f[4]=e=>L(y).modelForm.name=e),autocomplete:"off",placeholder:"请输入采集模型名称"},null,8,["disabled","modelValue"])])),_:1}),k(F,{label:"采集模型标签",prop:"label"},{default:C((()=>[k(v,{type:"text",modelValue:L(y).modelForm.label,"onUpdate:modelValue":f[5]||(f[5]=e=>L(y).modelForm.label=e),autocomplete:"off",placeholder:"请输入采集模型标签"},null,8,["modelValue"])])),_:1})])),_:1},8,["model","rules"])])])),_:1},8,["modelValue","title"])])}}},[["__scopeId","data-v-189462a2"]]);export{Fe as default};
