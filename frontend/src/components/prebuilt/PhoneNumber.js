import { forwardRef } from 'react'
import TextField from '@mui/material/TextField'


const phoneInput = (props, ref) => {
    return (

        <TextField
            {...props}
            inputRef={ref}
            fullWidth
            size='small'
            label='Número de teléfono'
            variant="standard"
            name='phone'
        />
    )
}
export default forwardRef(phoneInput)