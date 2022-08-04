import { useStore } from "vuex";


export default function useVStore() {
    const store = useStore()
    return store
}