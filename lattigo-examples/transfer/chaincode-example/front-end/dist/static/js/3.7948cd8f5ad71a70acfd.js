webpackJsonp([3],{Lsv8:function(t,e,a){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var n=a("mvHQ"),r=a.n(n),o={name:"Bank",data:function(){return{form:{}}},mounted:function(){this.init()},methods:{init:function(){null!==localStorage.getItem("bank")&&(this.form=JSON.parse(localStorage.getItem("bank")))},storage:function(){localStorage.setItem("bank",r()(this.form))},createBank:function(){var t=this;this.$axios.post("/v1/bank",r()(this.form)).then(function(e){0===e.data.validCode?t.$message.success("创建成功"):t.$message.error(e.data.msg)}).catch(function(t){console.log(t)})},getSecrete:function(){var t=document.createElement("iframe");t.src="/v1/bank/"+this.form.bankId+"/key",t.style.display="none",document.body.appendChild(t)}}},l={render:function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticStyle:{padding:"5% 20%",width:"50%","text-align":"center"},attrs:{id:"bank"}},[a("el-form",{staticStyle:{padding:"4% 20% 5% 15%","background-color":"#fff","border-radius":"8px"},attrs:{model:t.form,"label-width":"100px"}},[a("div",{staticStyle:{"font-size":"17px",margin:"7% 0px"}},[a("b",[t._v("创建银行")])]),t._v(" "),a("el-form-item",{attrs:{label:"银行名称："}},[a("el-input",{attrs:{placeholder:"请输入银行名称"},model:{value:t.form.bankId,callback:function(e){t.$set(t.form,"bankId",e)},expression:"form.bankId"}})],1),t._v(" "),a("el-row",[a("el-col",{staticStyle:{"text-align":"right"},attrs:{span:8}},[a("el-button",{attrs:{type:"primary"},on:{click:t.createBank}},[t._v("创建银行")])],1),t._v(" "),a("el-col",{staticStyle:{"text-align":"right"},attrs:{span:8}},[a("el-button",{attrs:{type:"success"},on:{click:t.storage}},[t._v("保存记录")])],1),t._v(" "),a("el-col",{staticStyle:{"text-align":"right"},attrs:{span:8}},[a("el-button",{attrs:{type:"primary"},on:{click:t.getSecrete}},[t._v("下载密钥")])],1)],1)],1)],1)},staticRenderFns:[]},s=a("VU/8")(o,l,!1,null,null,null);e.default=s.exports}});
//# sourceMappingURL=3.7948cd8f5ad71a70acfd.js.map