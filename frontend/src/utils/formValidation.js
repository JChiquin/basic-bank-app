export function hasEmptyField(form, fieldsExcluded = []) {
    let binaryAnd = 1;
    for (let key in form) {
        if (fieldsExcluded.includes(key))
            continue;
        binaryAnd &= isDirtyField(form[key])
        if (!binaryAnd)
            break;
    }
    return !binaryAnd
}

export function hasArrayEmptyField(array) {
    let binaryAnd = 1;
    for (let item of array) {
        binaryAnd &= !hasEmptyField(item)
        if (!binaryAnd)
            break;
    }
    return !binaryAnd
}

export function isValidEmail(email) {
    var re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

    return re.test(email);
}

//isNotEmptyField
export function isDirtyField(field) {
    return !!String(field).length && field != null && field != undefined
}


const messagesDefault = {
    required: "Campo requerido",
    email: "Debe ser un correo válido",
    minLength: (value) => `Mínimo ${value} caracteres`,
    maxLength: (value) => `Máximo ${value} caracteres`,
    sameAs: "Campo no coincide",
    minValue: (value) => `Valor mínimo ${value}`,
    maxValue: (value) => `Valor máximo ${value}`
}
export function hasFieldError(field, validators = []) {
    let isValid = true
    for (let validator of validators) {
        let message = ""
        let value = null
        let key = validator
        if (typeof validator === "object") {
            key = Object.keys(validator)[0]
            value = validator[key]
        }
        message = messagesDefault[key]
        switch (key) {
            case "required":
                isValid = isDirtyField(field)
                break;
            case "email":
                isValid = isValidEmail(field)
                break;
            case "minLength":
                isValid = field.length >= value
                message = messagesDefault[key](value)
                break;
            case "maxLength":
                isValid = field.length <= value
                message = messagesDefault[key](value)
                break;
            case "sameAs":
                isValid = field == value
                break;
            case "minValue":
                isValid = Number(field) >= value
                message = messagesDefault[key](value)
                break;
            case "maxValue":
                isValid = Number(field) <= value
                message = messagesDefault[key](value)
                break;
        }
        if (!isValid)
            return message
    }
    return null
}

export function hasFieldsErrors(fields, validators) {
    let result = {}
    for (const fieldName in fields) {
        result[fieldName] = hasFieldError(fields[fieldName], validators[fieldName])
    }
    return result
}

export function isObjNotEmpty(obj = {}) {
    return Object.values(obj).some(isDirtyField)
}

export function hasArrayOfFieldsErrors(arrayfields, validators) {
    let result = []
    for (const fields of arrayfields) {
        result.push(hasFieldsErrors(fields, validators))
    }
    return result
}

export function isArrayOfObjNotEmpty(arrayObj) {
    return arrayObj.some(isObjNotEmpty)
}