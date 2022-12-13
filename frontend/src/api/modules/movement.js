import { apiHttp } from "../axiosApi"

export const getMovementsAPI = (pagination, userID) => apiHttp("GET", `/v1/client/movement/${userID}`, null, pagination)
  