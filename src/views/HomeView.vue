<template>
  <div class="home">
    <div class="box" ref="box">
      <div class="left" ref="menu">
        <div class="left-content" v-loading="leftLoading">
          <div class="fileCtl">
            <el-button type="success" icon="el-icon-link" size="small" @click="OpenUrl" style="width: 100%">添加连接</el-button>
          </div>
          <el-divider content-position="left"><span style="font-size: 13px; color: #2b4b6b;">Service -> Methods</span></el-divider>
          <div style="overflow: auto; height: calc(100vh - 90px);">
          <el-tree
            :data="data"
            :props="defaultProps"
            @node-contextmenu="onContextShow"
            @node-click="handleNodeClick">
            <span slot-scope="{ node, data }" style="user-select: none;">
              <el-image :src="data.icon" style="width:25px;margin-right: 5px;"></el-image>
              <span style="font-size: 14px;color: #2b4b6b; position: absolute; line-height: 26px;">{{ node.label}}</span>
            </span>
          </el-tree>
          </div>  
        </div>
      </div>
      <div class="resize" title="收缩侧边栏" ref="menuResize">
        ┊
      </div>
      <div class="right" ref="opera">
        <el-tabs v-model="activeName" type="card" class="cards">
          <el-tab-pane :label="activeTopic" name="first">
                <span slot="label"><i class="el-icon-monitor"></i> {{activeTopic}}</span>
                <div class="request">
                    <span class="tit">Response</span>
                    <div class="divider"></div>
                   <b-code-editor v-model="requestJson" :indent-unit="4" height="100%" style="text-align: left" v-loading="requestLoading"/>
                </div>
                <el-tooltip content="发送请求" placement="top" effect="light">
                <el-button class="SendBtn" icon="el-icon-caret-right" circle type="primary" style="z-index: 5" @click="request"></el-button>
                </el-tooltip>
                <div class="response">
                  <span class="tit">Response</span>
                  <div class="divider"></div> 
                  <b-code-editor v-model="responseJson" :indent-unit="4" height="100%" style="text-align: left" v-loading="responseLoading"/>
                </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div> 
    <Contextmenu ref="contextmenu">
      <el-row><el-button size="small" style="border: 0;width: 100%;" @click="refreshService">刷新</el-button></el-row>
      <el-row><el-button size="small" style="border: 0;width: 100%;" @click="copyName">复制</el-button></el-row>
      <el-row><el-button size="small" style="border: 0;width: 100%;" @click="deleteTreeNode">删除</el-button></el-row>
    </Contextmenu>
  </div>
  
</template>

<script>
import Contextmenu from "vue-context-menu";

export default {
  name: 'HomeView',
  components:{
    Contextmenu,
  },
  data() {
    return {
      data: [],
        defaultProps: {
          children: 'children',
          label: 'label'
        },
      visible: false,
      top: 0,
      left: 0,
      leftLoading: false,
      activeName: 'first',
      activeTopic:'undefined',
      requestJson: '{}',
      responseJson: '{}',
      requestLoading: false,
      responseLoading: false,
      activeMethod: {},
      rightclickVal: {}
    }
  },
  methods:{
    OpenUrl(){
      this.$prompt('请输入grpc地址', '添加', {
        confirmButtonText: '添加',
        cancelButtonText: '取消',
      }).then(({value}) => {
        console.log(typeof value)
        if(typeof value == 'string' && value.length > 0) {
          this.AddFile(value)
          return
        }
        this.$message.info("未输入地址")
      })
    },
    onContextShow(e,data) {
      this.rightclickVal = data
      this.$refs.contextmenu.open();
    },
    handleNodeClick(data) {
      this.activeTopic = data.label
      console.log(data)
      this.activeMethod.service = data.father
      this.activeMethod.method = data.label
      this.activeMethod.url = data.url
      if(data.father != undefined) {
        this.getMethodParams(data.father, data.label, data.url)
      }
    },
    dropSize() {
      let that = this
      let min = 250
      this.$refs.menuResize.onmousedown = function () {
        document.onmousemove = function (e) {
          let clientX = e.clientX;
          if(clientX >= document.body.clientWidth * 0.8){
            clientX = document.body.clientWidth * 0.8
          }
          if(clientX <= min){
            clientX = min;
          }
          that.$refs.menu.style.width = clientX +"px";
          that.$refs.opera.style.width = document.body.clientWidth - clientX + "px"
        }
        document.onmouseup = function () {
        
          document.onmousemove = null;
          document.onmouseup = null;
          that.releaseCapture && that.releaseCapture()
        }
      }
    },
    async AddFile(url){
      let that = this
      let ax = this.$axios.create({
        baseUrl: 'http://127.0.0.1:10580',
      })
      ax.interceptors.request.use(function (config) {
        that.leftLoading = true

        return config
      })
      ax.interceptors.response.use(function (config) {
        that.leftLoading = false

        return config
      },function(error){
        that.leftLoading = false
        return Promise.reject(error)
      })
      try {
        const {
          data: res
        } = await ax.post('LinkMethods', {url: url})
        let insert = {
          label: '',
          icon: '',
          key: '',
          children: []
        }
        insert.label = res.services.name
        insert.icon = res.services.icon
        insert.key = res.services.key
        for(let i = 0; i < res.services.methods.length; i++){
          let ins = {
            label: res.services.methods[i].name,
            icon: res.services.methods[i].icon,
            inputType: res.services.methods[i].inputType,
            outputType: res.services.methods[i].outputType,
            father: res.services.methods[i].father,
            url: res.services.methods[i].url
          }
          insert.children.push(ins)
        }
        this.data = []
        this.data.push(insert)
        this.setFile()
      }catch(error) {
        console.log(error)
        this.$message.error(error.response.data)
      }
    },
    async getMethodParams(serviceName, methodName,url) {
      let that = this
      let ax = this.$axios.create({
        baseUrl: 'http://127.0.0.1:10580',
      })
      ax.interceptors.request.use(function (config) {
        that.requestLoading = true

        return config
      })
      ax.interceptors.response.use(function (config) {
        that.requestLoading = false

        return config
      },function(error){
        that.requestLoading = false
        return Promise.reject(error)
      })

      try {
        const {
          data: res
        } = await ax.post('MethodParam', {url: url, service: serviceName, method: methodName})

        this.requestJson = JSON.stringify(res.methods,null,'\t')
      }catch(error) {
        this.requestJson = JSON.stringify("Method no longer exists",null,'\t')
      }
    },
    async Query(serviceName, methodName,url,data) {
      let that = this
      let ax = this.$axios.create({
        baseUrl: 'http://127.0.0.1:10580',
      })
      ax.interceptors.request.use(function (config) {
        that.responseLoading = true

        return config
      })
      ax.interceptors.response.use(function (config) {
        that.responseLoading = false

        return config
      },function(error){
        that.responseLoading = false
        return Promise.reject(error)
      })

      try {
        const {
          data: res
        } = await ax.post('Call', {url: url, service: serviceName, method: methodName ,data: data})

        this.responseJson = JSON.stringify(res, null,'\t')
      }catch(error) {
        console.log(error)
         this.responseJson = JSON.stringify(error.response.data, null,'\t')
      }
    },
    request(){
      let js = JSON.parse(this.requestJson)
      let jsf =JSON.stringify(js)
      console.log(jsf)
      this.Query(this.activeMethod.service,this.activeMethod.method,this.activeMethod.url,jsf)
    },
    async setFile(){
      const {
        data: res
      } = await this.$axios.post('set',{data: this.data})
    },
    async getFile(){
      try {
        const {
        data: res
        } = await this.$axios.post('get')
        console.log(res)
        this.data = res.data
      }catch(error) {
        this.data = []
      }
    },
    deleteTreeNode() {
      if(this.rightclickVal.key != undefined){
        for(let index in this.data) {
          if(this.data[index] == this.rightclickVal) {
            this.data.splice(index,1)
            this.setFile()
          }
        }
      }else if(this.rightclickVal.father != undefined) {
        for(let index in this.data) {
          if(this.data[index].label == this.rightclickVal.father) {
            for(let id in this.data[index].children) {
              if(this.data[index].children[id] == this.rightclickVal) {
                this.data[index].children.splice(id,1)
                this.setFile()
              }
            }
          }
        }
      }
    },
    copyName() {
      var input = document.createElement("input")
      input.value = this.rightclickVal.label
      document.body.appendChild(input)
      input.select()
      document.execCommand("Copy")
      document.body.removeChild(input)
      this.$message.success("复制成功")
    },
    async refreshService(){
      if(this.rightclickVal.key != undefined) {
        let that = this
        let ax = this.$axios.create({
          baseUrl: 'http://127.0.0.1:10580',
        })
        ax.interceptors.request.use(function (config) {
          that.leftLoading = true

          return config
        })
        ax.interceptors.response.use(function (config) {
          that.leftLoading = false

          return config
        },function(error){
          that.leftLoading = false
          that.$message.error(error.response.data)
          return Promise.reject(error)
        })
        let url = this.rightclickVal.key.split("::")
        try {
          const {
            data: res
          } = await ax.post('LinkMethods', {url: url[0]})
          let change = {
            label: '',
            icon: '',
            key: '',
            children: []
          }
          change.label = res.services.name
          change.icon = res.services.icon
          change.key = res.services.key
          for(let i = 0; i < res.services.methods.length; i++){
            let ins = {
              label: res.services.methods[i].name,
              icon: res.services.methods[i].icon,
              inputType: res.services.methods[i].inputType,
              outputType: res.services.methods[i].outputType,
              father: res.services.methods[i].father,
              url: res.services.methods[i].url
            }
            change.children.push(ins)
          }

          console.log("before",this.data)
          for(let index in this.data) {
            if(this.data[index] == this.rightclickVal) {
              this.data.splice(index,1,change)
            }
          }
          this.setFile()
        }catch(error){
          this.$message.error(error)
        }
      }
    }
  },
  mounted(){
    this.dropSize();
    this.getFile()
  }
}
</script>

<style lang="scss">
.home {
  padding: 0;
  margin: 0;
  height: 100vh;
  
  .box {
    display: flex;
    flex-direction:row;

    .left {
      width: 250px;
      border-right: 1px solid rgb(231, 222, 222);
      height: 100vh;
      overflow: hidden;

      .left-content {
        margin: 10px;

        .fileCtl {
          position: static;
          z-index: 5;
          margin-bottom: 10px;
          padding: 10px;
        }
      }
    }

    .resize {
        cursor: col-resize;
        float: left;
        position: relative;
        line-height: 100vh;
        font-size: 15px;
        margin-left: 5px;
        margin-right: 5px;
        user-select: none;
        background: white;
    }
    .right {
      height: calc(100vh - 1px);
      width: calc(100% - 250px);
      margin-right: 10px;
      .cards {
        height: calc(100vh - 11px);
        padding-top: 10px;
        user-select: none;
        .SendBtn {
          position: absolute;
          left: calc(50% - 20px);
          top: 50%;
        }
        .tit {
          user-select: none;
          margin-top: 10px;
          height: 34px;
          line-height: 34px;
          font-weight: bold;
          color: #2b4b6b;
        }
        .divider {
          width: 100%;
          height: 1px;
          border-top: 1px solid #e0e0e0;
        }
        .request {
          width: 49%;
          height: calc(100vh - 80px);
          border-right: 1px solid #e0e0e0;
          border-top: 1px solid #e0e0e0;
          float: left;
          text-align: center;
          .CodeMirror{
            height: 100%;
            font-size: 14px;
          }
          .jsonCoder {
            text-align: left;
          }
        }

        .response {
          width: 49%;
          height: calc(100vh - 80px);
          border-top: 1px solid #e0e0e0;
          border-left: 1px solid #e0e0e0;
          float: right;
          .CodeMirror{
            height: 100%;
            font-size: 14px;
          }
          .jsonCoder {
            text-align: left;
          }
        }
      }
    }
  }
  .ctx-menu {
    min-width: 100px;

    span {
      font-size: 14px;
    }
  }
}
</style>