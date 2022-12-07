export function getJWT() {
    return localStorage.getItem("bank_jwt")
}

export function setJWT(value) {
    return localStorage.setItem("bank_jwt", value)
}

export function removeJWT() {
    return localStorage.removeItem("bank_jwt")
}