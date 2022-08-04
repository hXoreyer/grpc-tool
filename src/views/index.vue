<template>
    <div class="container">
        <div class="title">
            <div class="title-left" ref="title">
                <div class="title-icon">
                    <img src="../assets/icon.ico" style="width: 1.3rem">
                </div>
                <div class="title-content">
                    <span style="font-weight:bold; font-size: 1.3rem;color: #2f4f6f;user-select: none;">grpc-tool</span>
                </div>
            </div>
            <div class="title-right">
                <div class="right-icon-min right-icon" @click="minApp">
                    <n-icon size="1.3rem">
                        <RemoveOutline />
                    </n-icon>

                </div>

                <div class="right-icon-close right-icon" @click="quitApp">
                    <n-icon size="1.3rem">
                        <CloseOutline />
                    </n-icon>
                </div>
            </div>
        </div>
        <n-notification-provider>
            <div class="box">
                <div class="left" ref="menu">
                    <LeftMenu></LeftMenu>
                </div>
                <div class="resize" title="收缩侧边栏" ref="menuResize">
                    ┊
                </div>
                <div class="right" ref="opera">
                    <Tabs></Tabs>
                </div>
            </div>
        </n-notification-provider>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import useDrag from '../api/useDrag';
import LeftMenu from '../components/LeftMenu.vue'
import Tabs from '../components/Tabs.vue'
import { CloseOutline, RemoveOutline } from '@vicons/ionicons5'
import { GetFile } from '../types/request'
import store from '../store';



const menuResize = ref<HTMLDivElement | null>(null)
let menu = ref<HTMLDivElement | null>(null)
let opera = ref<HTMLDivElement | null>(null)
let title = ref<HTMLDivElement | null>(null)


const dropSize = () => {
    let min = 300
    if (menuResize.value) {
        menuResize.value.onmousedown = function () {
            document.onmousemove = function (e) {
                let clientX = e.clientX;
                if (clientX >= document.body.clientWidth * 0.6) {
                    clientX = document.body.clientWidth * 0.6
                }
                if (clientX <= min) {
                    clientX = min;
                }
                if (menu.value && opera.value) {
                    menu.value.style.width = clientX + "px";
                    opera.value.style.width = document.body.clientWidth - clientX + "px"
                }
            }
            document.onmouseup = function () {
                document.onmousemove = null;
                document.onmouseup = null;
            }
        }
    }
}

const quitApp = () => {
    (window as any).api.send('quit-app')
}

const minApp = () => {
    (window as any).api.send('min-app')
}

const getFile = async () => {
    try {
        const {
            data: res
        } = await GetFile()
        for(let i = 0; i < res.data.data.length; i++) {
            store.commit('addMenuVal', res.data.data[i])
        }
    }catch(error) {
        console.log(error)
    }
}

onMounted(() => {
    dropSize()
    if (title.value) {
        let fn = useDrag()
        fn(title.value)
    }
    getFile()
})

</script>

<style scoped>
.container {
    box-sizing: border-box;
    height: calc(100vh - 6px);
}

.box {
    display: flex;
    align-items: center;
}

.title {
    height: 50px;
    border-bottom: 1px solid rgb(231, 222, 222);
    box-shadow: 0 1px 1px rgba(10, 10, 10, 0.1);
    display: flex;
    flex-direction: row;
}

.title-left {
    display: flex;
    flex-direction: row;
    margin-left: 1rem;
    align-items: center;
    width: 100%;
}

.title-icon {
    margin-right: .35rem;
    display: flex;
    align-items: center;
    padding-top: .2rem;
    margin-right: .5rem;
}


.title-right {
    display: flex;
    align-items: center;
    margin-right: 1rem;
}

.right-icon {
    height: 1.3rem;
    padding: .3rem;
    border-radius: 5px;
    cursor: pointer;
    transition-duration: 0.5s;
    margin-right: 0.3rem;
}

.right-icon:hover {
    background: rgba(35, 35, 35, 0.3);
    color: white;
    transition-duration: 0.5s;
}

.right-icon-close:hover {
    background: #ff4c30;
}

.left {
    min-width: 300px;
    width: 300px;
    border-right: 1px solid rgb(231, 222, 222);
    overflow: hidden;
    height: calc(100vh - 57px - 1rem);
    padding-top: 1rem;
}

/*移动按钮样式*/
.resize {
    cursor: col-resize;
    float: left;
    position: relative;
    line-height: 100%;
    font-size: 15px;
    margin-left: 5px;
    margin-right: 5px;
    user-select: none;
    background: white;
}

.right {
    width: calc(100% - 330px);
    margin-right: 10px;
    height: calc(100vh - 57px - 1rem);
    padding-top: 1rem;
}
</style>