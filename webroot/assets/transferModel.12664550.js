import{s as e}from"./index.3a6602e1.js";var r={getModelList:r=>e.request({url:"/service/report/models",method:"get",data:r}),addModel:r=>e.request({url:"/service/report/model",method:"post",data:r}),editModel:r=>e.request({url:"/service/report/model",method:"put",data:r}),deleteModel:r=>e.request({url:"/service/report/model",method:"delete",data:r}),getPropertiesByModelIdList:r=>e.request({url:"/service/report/model/properties",method:"get",data:r}),addProperty:r=>e.request({url:"/service/report/model/property",method:"post",data:r}),editProperty:r=>e.request({url:"/service/report/model/property",method:"put",data:r}),deleteProperty:r=>e.request({url:"/service/report/model/properties",method:"delete",data:r}),addPropertyFromCSV:r=>e.request({url:"/service/report/model/properties/xlsx",method:"post",data:r}),exportProperty:r=>e.request({url:"/service/report/model/properties",method:"get",data:r}),reportNodes:r=>e.request({url:"/service/report/device/cmd/report",method:"post",data:r})};export{r as T};
