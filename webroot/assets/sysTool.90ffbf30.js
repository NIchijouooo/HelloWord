import{s as t}from"./index.559546be.js";var e={updateNTP:e=>t.request({url:"/ntp/param",method:"put",data:e}),getNTPInfo:e=>t.request({url:"/ntp/param",method:"get",data:e}),sendTimingCmd:e=>t.request({url:"/ntp/cmd",method:"post",data:e}),sendPingCmd:e=>t.request({url:"/system/ping",method:"post",data:e}),sendMessage:e=>t.request({url:"/system/commTool",method:"post",data:e})};export{e as S};
