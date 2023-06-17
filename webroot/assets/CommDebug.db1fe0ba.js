import{m as e,v as a,w as t,S as s,p as l,F as r,G as o,u as d,g as i}from"./element-plus.2c9426ce.js";import{_ as c,u as n}from"./index.d5c032bc.js";import{S as m}from"./sysTool.f104173e.js";import{I as u}from"./interface.52fdfea5.js";import{J as p,o as f,c as v,a as h,O as b,R as g,u as j,X as y,P as k,aj as _,V as x,W as C,L as D,S as V,av as w,aw as I}from"./@vue.2c474d7f.js";import"./lodash-es.98e49362.js";import"./async-validator.fb49d0f5.js";import"./@vueuse.d618444d.js";import"./dayjs.2790329b.js";import"./axios.765908e4.js";import"./@ctrl.b082b0c1.js";import"./@popperjs.36402333.js";import"./escape-html.e5dfadb9.js";import"./normalize-wheel-es.8aeb3683.js";import"./vue-router.3b7aa6c5.js";import"./pinia.39330980.js";import"./vue-demi.b3a9cad9.js";import"./@element-plus.e895b9e1.js";import"./screenfull.7cf96174.js";import"./vue-echarts.ebbe6a42.js";import"./resize-detector.4e96b72b.js";import"./echarts.2b8dcc79.js";import"./zrender.54f191d6.js";import"./tslib.60310f1a.js";import"./nprogress.aae2b787.js";const T={class:"main-container"},L={class:"main",style:{"background-color":"inherit"}},S=(e=>(w("data-v-3e56e29f"),e=e(),I(),e))((()=>h("div",{class:"card-header"},[h("span",null,"通讯调试助手")],-1))),U={class:"item"},A={class:"st-receive"},N={class:"str-content"},R={key:0,class:"str-noData"},z={key:0},F={class:"st-send"},O={class:"st-operate"},M={style:{"margin-bottom":"10px"}},P={class:"st-remark"};var E=c({__name:"CommDebug",setup(c){const w=n(),I=p({receiveData:[],sendData:"",curCollect:"",interfaceList:[{name:"collect01",label:"采集接口01"},{name:"modbusTCP",label:"采集接口02"}],receiveTime:!0,curCheck:0,checkOptions:[{name:0,label:"不校验"},{name:1,label:"CRC16"},{name:2,label:"SUM"}]});(()=>{const e={token:w.token,data:{}};u.getInterfaceList(e).then((e=>{console.log("getInterfaceList -> res",e),e&&("0"===e.code?I.interfaceList=e.data:E(e))}))})();const E=e=>{i({type:"error",message:e.message})};return(i,c)=>{const n=e,u=a,p=t,G=s,J=l,W=r,X=o,$=d;return f(),v("div",T,[h("div",L,[b(X,{gutter:20,style:{height:"100%"}},{default:g((()=>[b(W,{span:24},{default:g((()=>[b($,{class:"box-card",shadow:"hover"},{header:g((()=>[S])),default:g((()=>[h("div",U,[h("div",A,[h("div",N,[0==j(I).receiveData.length?(f(),v("div",R,"显示区")):y("",!0),(f(!0),v(k,null,_(j(I).receiveData,((e,a)=>(f(),v("div",{key:a},[j(I).receiveTime?(f(),v("div",z,"【"+C(e.date)+"】",1)):y("",!0),h("div",{class:D(1==e.type?"TxInfo":"RxInfo")},C(1==e.type?"Tx："+e.data:"Rx："+e.data),3)])))),128))])]),h("div",F,[b(n,{modelValue:j(I).sendData,"onUpdate:modelValue":c[0]||(c[0]=e=>j(I).sendData=e),class:"sts-sendArea",rows:"5",type:"textarea",placeholder:""==j(I).sendData?"发送区":""},null,8,["modelValue","placeholder"])]),h("div",O,[h("div",null,[b(p,{modelValue:j(I).curCollect,"onUpdate:modelValue":c[1]||(c[1]=e=>j(I).curCollect=e),style:{width:"100%","margin-bottom":"10px"},placeholder:"请选择采集接口"},{default:g((()=>[(f(!0),v(k,null,_(j(I).interfaceList,(e=>(f(),V(u,{key:e.collInterfaceName,label:e.collInterfaceName,value:e.collInterfaceName},null,8,["label","value"])))),128))])),_:1},8,["modelValue"]),h("div",M,[b(G,{modelValue:j(I).receiveTime,"onUpdate:modelValue":c[2]||(c[2]=e=>j(I).receiveTime=e),label:"接收时间",border:"",style:{width:"100%",height:"34px"}},null,8,["modelValue"])]),b(J,{type:"danger",style:{height:"34px",width:"100%"},plain:"",onClick:c[3]||(c[3]=e=>{I.receiveData=[]})},{default:g((()=>[x("清空显示区")])),_:1})]),h("div",null,[b(J,{type:"primary",style:{height:"34px",width:"100%","margin-bottom":"10px"},plain:"",onClick:c[4]||(c[4]=e=>(()=>{if(""===I.curCollect)return void E({message:"请选择采集接口！"});if(""===I.curCheck)return void E({message:"请选择校验方式！"});let e="";if(""===I.sendData)return void E({message:"发送数据不能为空！"});let a=I.sendData.replaceAll(" ","");if(!/^[0-9a-fA-F]+$/.test(a))return E({message:"报文格式错误，只能输入字符【0-9,a-f,A-F或者空格】！"}),I.sendData="",void(a="");const t=a.length;for(var s=0;s<t;s+=2)s<t&&(e+=a.substr(s,2)+(s+2>=t?"":" "));const l={token:w.token,data:{collInterfaceName:I.curCollect,directData:e,checkSum:I.curCheck}};m.sendMessage(l).then((e=>{"0"===e.code?e.data.forEach((e=>{I.receiveData.push(e)})):E(e)}))})())},{default:g((()=>[x(" 发送 ")])),_:1}),b(p,{modelValue:j(I).curCheck,"onUpdate:modelValue":c[5]||(c[5]=e=>j(I).curCheck=e),style:{width:"100%","margin-bottom":"10px"},placeholder:"请选择数据校验"},{default:g((()=>[(f(!0),v(k,null,_(j(I).checkOptions,(e=>(f(),V(u,{key:e.name,label:e.label,value:e.name},null,8,["label","value"])))),128))])),_:1},8,["modelValue"]),b(J,{type:"danger",style:{height:"34px",width:"100%"},plain:"",onClick:c[6]||(c[6]=e=>{I.sendData=""})},{default:g((()=>[x("清空发送区")])),_:1})])]),h("div",P,[b(X,{gutter:16},{default:g((()=>[b(W,{span:1},{default:g((()=>[x("操作步骤:")])),_:1}),b(W,{span:21},{default:g((()=>[x("1、选择采集接口")])),_:1})])),_:1}),b(X,{gutter:16},{default:g((()=>[b(W,{offset:1,span:21},{default:g((()=>[x("2、选择校验方式")])),_:1})])),_:1}),b(X,{gutter:16},{default:g((()=>[b(W,{offset:1,span:21},{default:g((()=>[x("3、发送区编写发送的报文")])),_:1})])),_:1}),b(X,{gutter:16},{default:g((()=>[b(W,{offset:1,span:21},{default:g((()=>[x("4、显示区显示发送的报文和返回的报文")])),_:1})])),_:1})])])])),_:1})])),_:1})])),_:1})])])}}},[["__scopeId","data-v-3e56e29f"]]);export{E as default};
