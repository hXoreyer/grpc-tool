import { createStore } from 'vuex'

export interface jsonType {
    name: string
    data: object
}

export interface menuChildType {
    father: string
    inputType: string
    outputType: string
    name: string
    url: string
}

export interface menuType {
    key: string
    name: string
    childs: menuChildType[]
}

export interface jsonArrType {
    arr: jsonType[]
}

export interface menuArrType {
    arr: menuType[]
}

export interface storeType {
    jsonVals: jsonArrType
    returnVals: jsonArrType
    menuVals: menuArrType
    tabVals: string[],
    newTab: string
}


const store = createStore<storeType>({
    state() {
        return{
            jsonVals: {
                arr: []
            },
            returnVals: {
                arr: []
            },
            menuVals: {
                arr:[]
            },
            tabVals: [],
            newTab: 'undefined'
        }
    },
    mutations: {
        addJsonVal(state,val) {
            state.jsonVals.arr.push(val)
        },
        setJsonVal(state,val){
            for(let i = 0; i < state.jsonVals.arr.length; i++) {
                if(state.jsonVals.arr[i].name === val.name) {
                    state.jsonVals.arr[i].data = val.data
                }
            }
        },
        addreturnVal(state,val) {
            state.jsonVals.arr.push(val)
        },
        setreturnVal(state,val){
            for(let i = 0; i < state.returnVals.arr.length; i++) {
                if(state.returnVals.arr[i].name === val.name) {
                    state.returnVals.arr[i].data = val.data
                }
            }
        },
        addMenuVal(state,val) {
            state.menuVals.arr.push(val)
        },
        deleteMenuVal(state,val) {
            state.menuVals.arr.splice(val,1)
        },
        deleteMenuChildVal(state,val){
            state.menuVals.arr[val.index].childs.splice(val.ci,1)
        },
        deleteMenuByKey(state,key) {
            for(let i = 0; i < state.menuVals.arr.length; i++) {
                if(state.menuVals.arr[i].key === key) {
                    state.menuVals.arr.splice(i,1)
                }
            }
        },
        addTabVal(state,val) {
            state.tabVals.push(val)
        },
        deleteTabVal(state,val){
            state.tabVals.splice(val,1)
        },
        setUndefined(state) {
            state.tabVals.splice(0,1)
        },
        setNewTab(state,val){
            state.newTab = val
        },
        setMenuVal(state,val) {
            for(let i = 0; i < state.menuVals.arr.length; i++) {
                if(state.menuVals.arr[i].key === val.key) {
                    state.menuVals.arr[i].childs = val.childs
                }
            }
        }
    },
    getters: {
        getJsonVal(state) {
            return state.jsonVals.arr
        },
        getReturnVal(state) {
            return state.returnVals.arr
        },
        getMenuVal(state) {
            return state.menuVals.arr
        },
        getTabVal(state){
            return state.tabVals
        },
        getNewTab(state) {
            return state.newTab
        }
    }
})


export default store 
