<script lang="ts" setup>
import JsonEditorVue from 'json-editor-vue3'
import { reactive, toRefs, watch } from 'vue';
import useVStore from '../api/useVStore.js';
import { jsonType } from '../store/index.js';

const store = useVStore()
const props = defineProps({
    keyname: String,
    type: String
})
let {keyname, type} = toRefs(props)
let jsonVal = reactive<jsonType>({
    name: keyname.value,
    data: {}
})
const changeJson = (data:any) => {
    if(type.value === "request") {
        store.commit('setJsonVal',jsonVal)
    }
}

let listen = () => {
    return type.value === "request" ? store.getters.getJsonVal : store.getters.getReturnVal
}

watch(listen() , (newVal,oldVal)=> {
    for(let i = 0; i < newVal.length; i++) {
        if(newVal[i].name === keyname.value) {
            jsonVal.data = newVal[i].data
            return
        }
    }
})


</script>
<template>
    <JsonEditorVue v-model="jsonVal.data" @update:modelValue="changeJson" currentMode="code"/>
</template>
<style>
.jsoneditor-menu, .jsoneditor-poweredBy, .jsoneditor-statusbar, .full-screen.show.right {
    display: none;
}

.jsoneditor {
    border: 0;
    height: calc(100vh - 123px - 2rem);
}


div.jsoneditor-outer.has-main-menu-bar {
    margin-top: 0;
    padding-top: 0;
    border: 1px solid rgb(231, 222, 222);
    border-top: 0;
}

.ace_gutter-layer, .ace_scroller {
    margin-top: .5rem;
}

</style>