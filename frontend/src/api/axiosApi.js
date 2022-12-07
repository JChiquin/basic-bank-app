import axios from 'axios'
import { getJWT } from "../utils/localStorage"

export const API_URL_BACKEND = process.env.REACT_APP_API_URL_BACKEND;
const AXIOS_TIMEOUT_MS = process.env.REACT_APP_AXIOS_TIMEOUT_MS || 10000;

const defaultHeaders = {
    Accept: "application/json",
    "Content-Type": "application/json",
};

export const apiHttp = async (method, endpoint, data = null, params = null, options = {}) => {

    // header options 
    options.headers = {
        ...defaultHeaders,
        ...options.headers
    }

    let jwt = getJWT()
    if (jwt) {
        options.headers["Authorization"] = `Bearer ${jwt}`
    }

    let serviceResponse = {}
    const url = `${API_URL_BACKEND}${endpoint}`

    const servicePromise = axios({
        method: method.toLowerCase(),
        url,
        params,
        data,
        timeout: AXIOS_TIMEOUT_MS,
        ...options
    });

    console.log(`${method.toUpperCase()} ${url}`);

    try {
        const materializedPromise = await servicePromise;
        console.log("promise", materializedPromise)
        serviceResponse = materializedPromise.data;
        serviceResponse.headers = materializedPromise.headers;
    } catch (error) {
        if (error.response) {
            console.log("apiHttp -> error.response", error.response)
            serviceResponse = error.response.data;
        } else {
            console.log("apiHttp -> error", error)
            serviceResponse = error
        }
    }

    return serviceResponse;
};
