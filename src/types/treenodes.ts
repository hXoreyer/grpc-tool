import {TreeOption} from 'naive-ui'
import { reactive } from 'vue'

class TreeNode {
    data: TreeOption[] = reactive([])

    PushOption = (val:TreeOption)=>{
        this.data.push(val)
    }

    PushOptionToChild = (parent:TreeOption,val:TreeOption) => {
        for(let index = 0; index < this.data.length; index++) {
            if(this.data[index].key === parent.key) {
                this.data[index].children.push(val)
            }
        }
    }

    DeleteOption = (val:TreeOption)=> {
        for(let index = 0; index < this.data.length; index++) {
            if(this.data[index].key === val.key) {
                this.data.splice(index,1)
            }
        }
    }

    DeleteOptionByIndex = (index:number) => {
        this.data.splice(index,1)    
    }

    SetOptions = (options:TreeOption[])=> {
        this.data = []
        this.data = options
    }

    Clear = () => {
        this.data.splice(0,this.data.length)
    }
}

export default TreeNode