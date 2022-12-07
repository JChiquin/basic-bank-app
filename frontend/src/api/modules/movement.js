import { apiHttp } from "../axiosApi"

export const getMovementsAPI = (pagination) => apiHttp("GET", `/v1/client/movement`, null, pagination)
  