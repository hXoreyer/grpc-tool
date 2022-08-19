import axios from 'axios'

/*
let g_port

await (window as any).api.invoke('send-port').then(res=>{
    g_port = res
})

*/

const service = axios.create({
    baseURL: 'http://localhost:10580',
    timeout: 60000
})


service.interceptors.request.use(
    (config) => {
        
        return config
    },
    (error) => {
        
        return Promise.reject(error)
    }
)

service.interceptors.response.use(
	function (response) {

		return response
	},
	function (error) {

		return Promise.reject(error)
	}
)

export default service
