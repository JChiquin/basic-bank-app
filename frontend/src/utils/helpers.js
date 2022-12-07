export function ccyFormat(num) {
    return `${num.toFixed(2)}`;
}

export const getAPIError = (errors = []) => errors[0]?.error