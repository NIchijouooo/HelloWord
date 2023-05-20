import{i as t}from"./vue-demi.b3a9cad9.js";import{ae as e,f as n,ad as s,q as o,h as c,Y as a,B as r,aq as i,ag as u,t as p,A as f,H as l,a5 as h,p as d}from"./@vue.8518410d.js";
/*!
  * pinia v2.0.14
  * (c) 2022 Eduardo San Martin Morote
  * @license MIT
  */let y;const v=t=>y=t,b=Symbol();function _(t){return t&&"object"==typeof t&&"[object Object]"===Object.prototype.toString.call(t)&&"function"!=typeof t.toJSON}var j,$;function m(){const o=e(!0),c=o.run((()=>n({})));let a=[],r=[];const i=s({install(t){v(i),i._a=t,t.provide(b,i),t.config.globalProperties.$pinia=i,r.forEach((t=>a.push(t))),r=[]},use(e){return this._a||t?a.push(e):r.push(e),this},_p:a,_a:null,_e:o,_s:new Map,state:c});return i}($=j||(j={})).direct="direct",$.patchObject="patch object",$.patchFunction="patch function";const O=()=>{};function g(t,e,n,s=O){t.push(e);const c=()=>{const n=t.indexOf(e);n>-1&&(t.splice(n,1),s())};return!n&&o()&&f(c),c}function P(t,...e){t.slice().forEach((t=>{t(...e)}))}function w(t,e){for(const n in e){if(!e.hasOwnProperty(n))continue;const s=e[n],o=t[n];_(o)&&_(s)&&t.hasOwnProperty(n)&&!r(s)&&!i(s)?t[n]=w(o,s):t[n]=s}return t}const S=Symbol();const{assign:E}=Object;function A(t,o,p={},f,h,d){let y;const b=E({actions:{}},p),$={deep:!0};let m,A,I,q=s([]),x=s([]);const F=f.state.value[t];let k;function B(e){let n;m=A=!1,"function"==typeof e?(e(f.state.value[t]),n={type:j.patchFunction,storeId:t,events:I}):(w(f.state.value[t],e),n={type:j.patchObject,payload:e,storeId:t,events:I});const s=k=Symbol();l().then((()=>{k===s&&(m=!0)})),A=!0,P(q,n,f.state.value[t])}d||F||(f.state.value[t]={}),n({});const H=O;function J(e,n){return function(){v(f);const s=Array.from(arguments),o=[],c=[];function a(t){o.push(t)}function r(t){c.push(t)}let i;P(x,{args:s,name:e,store:N,after:a,onError:r});try{i=n.apply(this&&this.$id===t?this:N,s)}catch(u){throw P(c,u),u}return i instanceof Promise?i.then((t=>(P(o,t),t))).catch((t=>(P(c,t),Promise.reject(t)))):(P(o,i),i)}}const M={_p:f,$id:t,$onAction:g.bind(null,x),$patch:B,$reset:H,$subscribe(e,n={}){const s=g(q,e,n.detached,(()=>o())),o=y.run((()=>c((()=>f.state.value[t]),(s=>{("sync"===n.flush?A:m)&&e({storeId:t,type:j.direct,events:I},s)}),E({},$,n))));return s},$dispose:function(){y.stop(),q=[],x=[],f._s.delete(t)}},N=a(E({},M));f._s.set(t,N);const Y=f._e.run((()=>(y=e(),y.run((()=>o())))));for(const e in Y){const n=Y[e];if(r(n)&&(!r(C=n)||!C.effect)||i(n))d||(!F||_(z=n)&&z.hasOwnProperty(S)||(r(n)?n.value=F[e]:w(n,F[e])),f.state.value[t][e]=n);else if("function"==typeof n){const t=J(e,n);Y[e]=t,b.actions[e]=n}}var z,C;return E(N,Y),E(u(N),Y),Object.defineProperty(N,"$state",{get:()=>f.state.value[t],set:t=>{B((e=>{E(e,t)}))}}),f._p.forEach((t=>{E(N,y.run((()=>t({store:N,app:f._a,pinia:f,options:b}))))})),F&&d&&p.hydrate&&p.hydrate(N.$state,F),m=!0,A=!0,N}function I(t,e,n){let c,a;const r="function"==typeof e;function i(t,n){const i=o();(t=t||i&&p(b))&&v(t),(t=y)._s.has(c)||(r?A(c,e,a,t):function(t,e,n,o){const{state:c,actions:a,getters:r}=e,i=n.state.value[t];let u;u=A(t,(function(){i||(n.state.value[t]=c?c():{});const e=h(n.state.value[t]);return E(e,a,Object.keys(r||{}).reduce(((e,o)=>(e[o]=s(d((()=>{v(n);const e=n._s.get(t);return r[o].call(e,e)}))),e)),{}))}),e,n,0,!0),u.$reset=function(){const t=c?c():{};this.$patch((e=>{E(e,t)}))}}(c,a,t));return t._s.get(c)}return"string"==typeof t?(c=t,a=r?n:e):(a=t,c=t.id),i.$id=c,i}export{m as c,I as d};
