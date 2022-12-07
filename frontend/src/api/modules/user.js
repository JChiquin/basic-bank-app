import { apiHttp } from "../axiosApi"

export const loginAPI = (loginValues) => apiHttp("POST", `/v1/public/client/user/login`, loginValues)

export const whoAmIAPI = (userId) => apiHttp("GET", `/v1/client/user/whoami/${userId}`)

export const registerAPI = (registerValues) => apiHttp("POST", `/v1/public/client/user/register`, registerValues)