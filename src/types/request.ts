import HttpRequest from '../request/index'

export const AddFile = (url: string)=> {
    return HttpRequest.post('LinkMethods', {url: url})
}

export const Query = (serviceName: string, methodName: string, url: string, data: string)=> {
    return HttpRequest.post('Call', {url: url, service: serviceName, method: methodName ,data: data})
}

export const GetMethodParam = (serviceName, methodName,url)=> {
    return HttpRequest.post('MethodParam', {url: url, service: serviceName, method: methodName})
}

export const SetFile = (data)=>{
    return HttpRequest.post('set',{data: data})
}

export const GetFile = ()=>{
    return HttpRequest.post('get')
}