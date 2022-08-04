<script lang="ts" setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { CaretForward } from '@vicons/ionicons5'
import Editor from './Editor.vue'
import useVStore from '../api/useVStore';
import { GetMethodParam, Query } from '../types/request'
import { jsonType } from '../store/index'
import { useNotification } from 'naive-ui';

const notification = useNotification()
const store = useVStore()
let pvalue = ref<string>('1')
const panels = reactive<string[]>([])
let loading = ref<boolean>(false)
const closable = computed(() => {
    return panels.length > 1
})

//localhost:8090
const handleAdd = (name: string, b: boolean) => {
    if (panels.length === 1 && panels[0] === 'undefined') {
        store.commit('setUndefined')
        panels.splice(0, 1)
    }
    panels.push(name)
    pvalue.value = name
    if (b) store.commit('addTabVal', name)

    let me = store.getters.getMenuVal
    for (let i = 0; i < me.length; i++) {
        for (let j = 0; j < me[i].childs.length; j++) {
            let keyname = me[i].key + ":" + me[i].childs[j].name
            if (name === keyname) {
                GetParam(me[i].childs[j].father, me[i].childs[j].name, me[i].childs[j].url, keyname)
                return
            }
        }
    }
}


watch(() => store.getters.getNewTab, (newVal: string, oldVal) => {
    let tabs = store.getters.getTabVal
    for (let i = 0; i < tabs.length; i++) {
        if (tabs[i] === newVal) {
            pvalue.value = newVal
            return
        }
    }
    handleAdd(newVal, true)
}, { immediate: true })

const handleClose = (name: string) => {
    const nameIndex = panels.findIndex((panelName) => panelName === name)
    if (!~nameIndex) return
    panels.splice(nameIndex, 1)
    store.commit('deleteTabVal',nameIndex)
    if (name === pvalue.value) {
        pvalue.value = panels.at(-1)
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
        store.commit('addJsonVal', newVal)
    } catch (error) {
        notification['error']({
            content: '错误',
            meta: '获取参数失败',
            duration: 2500,
            keepAliveOnHover: true
        })
    }
}

//localhost:8090
const changeVal = (name: string) => {
    const index = panels.findIndex((panelName) => panelName === name)
    pvalue.value = panels[index]
}

onMounted(() => {
    let tabs = store.getters.getTabVal
    if (panels.length === 0 && tabs.length > 0) {
        for (let i = 0; i < tabs.length; i++) {
            panels.push(tabs)
        }
        pvalue.value = tabs[0]
    }
})

const query = async (serviceName, methodName, url, data, key) => {
    loading.value = true
    try {
        const {
            data: res
        } = await Query(serviceName, methodName, url, data)
        let st = store.getters.getReturnVal
        if (st.length === 0) {
            let newVal: jsonType = {
                name: key,
                data: res
            }
            st.push(newVal)
        }
        for (let i = 0; i < st.length; i++) {
            if (key === st[i].name) {
                st[i].data = res
                loading.value = false
                return
            }
        }
        let newVal: jsonType = {
            name: key,
            data: res
        }
        st.push(newVal)
    } catch (error) {
        notification['error']({
            content: '错误',
            meta: '请求失败',
            duration: 2500,
            keepAliveOnHover: true
        })
        loading.value = false
    }
}

//localhost:8090
const RequestForm = () => {
    let st = store.getters.getJsonVal
    let m = store.getters.getMenuVal
    for (let i = 0; i < st.length; i++) {
        if (st[i].name === pvalue.value) {
            for (let j = 0; j < m.length; j++) {
                for (let n = 0; n < m[j].childs.length; n++) {
                    let keys = m[j].key + ":" + m[j].childs[n].name
                    if (keys === pvalue.value) {
                        let dat = JSON.stringify(st[i].data)
                        query(m[j].childs[n].father, m[j].childs[n].name, m[j].childs[n].url, dat, pvalue.value)
                    }
                }
            }
        }
    }
}

</script>
<template>
    <n-tabs v-model:value="pvalue" type="card" :closable="closable" @close="handleClose"
        @add="handleAdd" @update:value="changeVal">
        <n-tab-pane v-for="panel in panels" :key="panel" :name="panel" display-directive="show">
            <div class="content">
                <div class="left">
                    <div class="qrstyle">request</div>
                    <Editor :keyname="panel" type="request"></Editor>
                </div>
                <n-popover trigger="hover">
                    <template #trigger>
                        <n-button circle secondary type="error" class="staticBtn" size="large" @click="RequestForm">
                            <template #icon>
                                <n-icon>
                                    <CaretForward />
                                </n-icon>
                            </template>
                        </n-button>
                    </template>
                    <span>发送请求</span>
                </n-popover>
                <div class="right">
                    <div class="qrstyle">response</div>
                    <n-spin :show="loading">
                        <Editor :keyname="panel" type="response"></Editor>
                    </n-spin>
                </div>
            </div>
        </n-tab-pane>
    </n-tabs>
</template>
<style scoped>
.content {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    position: relative;
}

.left {
    width: 49%;
    height: calc(100vh - 98px - 2rem);
    border: 0;
    display: flex;
    flex-direction: column;
    text-align: center;
    user-select: none;
    color: #2f4f6f;
    font-size: 1rem;
}

.right {
    max-width: 49%;
    width: 49%;
    height: calc(100vh - 98px - 2rem);
    display: flex;
    flex-direction: column;
    text-align: center;
    user-select: none;
    color: #2f4f6f;
    font-size: 1rem;
}

.staticBtn {
    position: absolute;
    left: calc(50% - 20px);
    top: 50%;
    z-index: 50;
}

.qrstyle {
    border: 1px solid rgb(231, 222, 222);
}
</style>