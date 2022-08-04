<script lang="ts" setup>
import { LinkOutline } from '@vicons/ionicons5/'
import MethodMenu from './MethodMenu.vue';
import { ref } from 'vue'
import { AddFile, SetFile } from '../types/request'
import { useNotification } from 'naive-ui';
import useVStore from '../api/useVStore';
import { menuType, menuChildType } from '../store/index'

const store = useVStore()
const notification = useNotification()
let showModal = ref<boolean>(false)
let address = ref<string>('')
let loading = ref<boolean>(false)
const noSideSpace = (value: string) => !value.startsWith(' ') && !value.endsWith(' ')

const AddLink = async () => {
    showModal.value = false
    loading.value = true
    if (address.value.length < 2) return
    try {
        const {
            data: res
        } = await AddFile(address.value)
        let vals: menuType[] = store.getters.getMenuVal
        let da = res.services
        for (let i = 0; i < vals.length; i++) {
            if (vals[i].key === da.key) {
                notification['warning']({
                    content: '服务重复',
                    meta: "此连接的服务已存在",
                    duration: 2500,
                    keepAliveOnHover: true
                })
                loading.value = false
                return
            }
        }
        let newVal: menuType = {
            key: da.key,
            name: da.name,
            childs: []
        }
        for (let i = 0; i < da.methods.length; i++) {
            let method: menuChildType = {
                father: '',
                inputType: '',
                name: '',
                outputType: '',
                url: ''
            }
            method.father = da.methods[i].father
            method.inputType = da.methods[i].inputType
            method.name = da.methods[i].name
            method.outputType = da.methods[i].outputType
            method.url = da.methods[i].url
            newVal.childs.push(method)
        }
        store.commit('addMenuVal', newVal)
        setFile()
        loading.value = false

    } catch (error) {
        notification['error']({
            content: '获取错误',
            meta: '获取失败',
            duration: 2500,
            keepAliveOnHover: true
        })
        loading.value = false
    }

}


const setFile = async () => {
    try {
        const {
            data: res
        } = await SetFile({data: store.getters.getMenuVal})
        console.log(res)
    }catch(error) {
        console.log(error)
    }

}


</script>
<template>
    <div class="left-content">
        <div class="fileCtl">
            <n-button style="width:100%" color="#8a2be2" size="medium" @click="showModal = true">
                <template #icon>
                    <n-icon>
                        <LinkOutline />
                    </n-icon>
                </template>
                添加连接
            </n-button>
        </div>
        <n-divider>
            <span style="font-size: 13px; color: #2b4b6b;user-select: none;">Service -> Methods</span>
        </n-divider>
        <n-spin :show="loading">
            <MethodMenu></MethodMenu>
        </n-spin>
    </div>
    <n-modal v-model:show="showModal" preset="dialog" title="Dialog">
        <template #header>
            <div>添加连接</div>
        </template>
        <n-input style="margin-top: 1rem" placeholder="http[s]://地址" :allow-input="noSideSpace" v-model:value="address"
            clearable></n-input>
        <template #action>
            <n-button type="primary" @click="AddLink">添加</n-button>
        </template>
    </n-modal>
</template>
<style scoped>
.left-content {
    margin: 10px;
    height: 100%;
}

.n-divider--title-position-center {
    margin: 1rem 0;
}
</style>