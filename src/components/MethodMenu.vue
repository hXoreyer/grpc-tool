<template>
  <n-input-group>
    <n-input v-model:value="pattern" placeholder="搜索" size="small" style="margin-bottom: 1rem" />
    <!---
    <n-popover trigger="hover">
      <template #trigger>
        <n-button size="small" type="info">
          <template #icon>
            <n-icon>
              <Refresh />
            </n-icon>
          </template>
        </n-button>
      </template>
      <span>刷新</span>
    </n-popover>
    -->
  </n-input-group>
  <n-tree :show-irrelevant-nodes="true" :pattern="pattern" :data="data.data" block-line v-if="bShow"
    :node-props="ClickTree" />
  <n-empty description="无数据" v-if="!bShow">
    <template #icon>
      <n-icon>
        <ReceiptOutline />
      </n-icon>
    </template>
  </n-empty>
  <n-dropdown trigger="manual" placement="bottom-start" :show="showDropdown" :options="(options as any)" :x="x" :y="y"
    @select="HandleSelect" @clickoutside="handleClickoutside" />
</template>

<script setup lang="ts">
import { h, ref, reactive, watch, onMounted } from 'vue'
import { DropdownOption, NButton, NIcon, TreeOption, useNotification } from 'naive-ui'
import { InformationCircleOutline, ReceiptOutline, Infinite, Refresh } from '@vicons/ionicons5'
import TreeNode from '../types/treenodes'
import useVStore from '../api/useVStore'
import { GetMethodParam, SetFile } from '../types/request'
import { jsonType } from '../store'

const notification = useNotification()
const store = useVStore()
const data: TreeNode = reactive(new TreeNode)
let bShow = ref(false)
let pattern = ref('')
let showDropdown = ref(false)
let x = ref(0)
let y = ref(0)
let options = ref<DropdownOption[]>([])
let clickOption = reactive<TreeOption>({})
options.value = [
  { key: 'delete', label: '删除' },
  { key: 'update', label: '刷新' },
  { key: 'copy', label: '复制' }
]

const handleClickoutside = () => {
  showDropdown.value = false
}

const addMenu = (val) => {
  let node: TreeOption = reactive({
    label: val.name,
    key: val.key,
    prefix: () => h(NIcon, { size: '1.2rem', color: '#cf2f74' }, () => h(InformationCircleOutline)),
    children: []
  })


  for (let i = 0; i < val.childs.length; i++) {
    let key = val.key + ":" + val.childs[i].name
    let child = {
      label: val.childs[i].name,
      key: key,
      prefix: () => h(NIcon, { size: '1.2rem', color: '#b2de27' }, () => h(Infinite)),
    }
    node.children.push(child)
  }
  data.PushOption(node)
}

const SetMenu = (vals) => {
  for (let i = 0; i < vals.length; i++) {
    addMenu(vals[i])
  }
}

watch(store.getters.getMenuVal, (newVal, oldVal) => {
  data.Clear()
  SetMenu(newVal)
}, { immediate: true, deep: true })


watch(data.data, (newVal, oldVal) => {
  if (newVal.length > 0) {
    bShow.value = true
  } else {
    bShow.value = false
  }
})

onMounted(() => {
  if (store.getters.getMenuVal.length > 0 && data.data.length === 0) {
    data.Clear()
    SetMenu(store.getters.getMenuVal)
    bShow.value = true
  } else if (store.getters.getMenuVal.length > 0 && data.data.length > 0) {
    bShow.value = true
  }
})

const setFile = async () => {
  try {
    const {
      data: res
    } = await SetFile({ data: store.getters.getMenuVal })
    console.log(res)
  } catch (error) {
    console.log(error)
  }

}

const GetParam = async (serviceName, methodName, url, key) => {
  try {
    const {
      data: res
    } = await GetMethodParam(serviceName, methodName, url)

    let newVal: jsonType = {
      name: key,
      data: res.methods
    }
    store.commit('setJsonVal', newVal)
  } catch (error) {
    notification['error']({
      content: '错误',
      meta: '刷新错误',
      duration: 2500,
      keepAliveOnHover: true
    })

  }
}

const HandleSelect = (key: string | number, option: DropdownOption) => {
  showDropdown.value = false
  if (clickOption.children) return
  let st = store.getters.getMenuVal
  if (key === 'delete') {
    for (let i = 0; i < st.length; i++) {
      for (let j = 0; j < st[i].childs.length; j++) {
        let name = st[i].key + ":" + st[i].childs[j].name
        if (name === clickOption.key) {
          store.commit('deleteMenuChildVal', { index: i, ci: j })
          if (store.getters.getMenuVal[i].childs.length === 0) {
            store.commit('deleteMenuVal', i)
          }
          setFile()
          return
        }
      }
    }
  } else if (key === 'copy') {
    var input = document.createElement("input")
    for (let i = 0; i < st.length; i++) {
      for (let j = 0; j < st[i].childs.length; j++) {
        let name = st[i].key + ":" + st[i].childs[j].name
        if (name === clickOption.key) {
          input.value = st[i].childs[j].name
          break
        }
      }
    }
    document.body.appendChild(input)
    input.select()
    document.execCommand("Copy")
    document.body.removeChild(input)
    notification['success']({
      content: '复制成功',
      meta: "已成功将方法名称置入剪切板",
      duration: 2500,
      keepAliveOnHover: true
    })
  } else if (key === 'update') {
    let me = store.getters.getMenuVal
    for (let i = 0; i < me.length; i++) {
      for (let j = 0; j < me[i].childs.length; j++) {
        let keyname = me[i].key + ":" + me[i].childs[j].name
        if (clickOption.key === keyname) {
          GetParam(me[i].childs[j].father, me[i].childs[j].name, me[i].childs[j].url, keyname)
          return
        }
      }
    }
  }
}


const ClickTree = ({ option }: { option: TreeOption }) => {
  return {
    onClick() {
      if (option.children) return
      store.commit('setNewTab', option.key)
    },
    onContextmenu(e: MouseEvent): void {
      clickOption = option
      console.log(option)
      showDropdown.value = true
      x.value = e.clientX
      y.value = e.clientY
      e.preventDefault()
    }

  }

}

</script>