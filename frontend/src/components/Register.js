import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { DEFAULT_TYPES } from '../constant'
import PhoneInput from 'react-phone-number-input'
import 'react-phone-number-input/style.css'
import CustomPhoneNumber from './prebuilt/PhoneNumber'
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';
import FormLabel from '@mui/material/FormLabel';
import LoadingButton from '@mui/lab/LoadingButton';
import {
    login,
    selectIsLogged
} from '../redux/user/userSlice';
import { useDispatch, useSelector } from 'react-redux';
import { registerAPI } from '../api/modules/user';
import { useSnackbar } from 'notistack';

import { hasFieldsErrors, isObjNotEmpty } from '../utils/formValidation'
import { getAPIError } from '../utils/helpers'


const theme = createTheme();

const DEFAULT_FORM = {
    document_type: DEFAULT_TYPES[0],
    document_number: '',
    first_name: '',
    last_name: '',
    birth_date: '',
    phone_number: '',
    email: '',
    password: '',
    repeatPassword: ''
}


export default function SignUp() {


    const { enqueueSnackbar } = useSnackbar();
    const isLogged = useSelector(selectIsLogged);
    const navigate = useNavigate()
    const dispatch = useDispatch()

    const [form, setForm] = React.useState(DEFAULT_FORM)
    const [loading, setLoading] = React.useState(false)
    const [formErrors, setFormErrors] = React.useState(DEFAULT_FORM)

    const FORM_VALIDATORS = {
        document_number: ["required", { maxLength: 20 }],
        first_name: ["required", { maxLength: 40 }],
        last_name: ["required", { maxLength: 40 }],
        birth_date: ["required"],
        phone_number: ["required", { maxLength: 20 }],
        email: ["required", "email"],
        password: ["required", { maxLength: 16 }, { minLength: 8 }],
        repeatPassword: ["required", { sameAs: form.password }],
    }

    const handleChangeForm = (event) => {
        setForm(Object.assign({}, form, { [event.target.name]: event.target.value }))
        setFormErrors(Object.assign({}, formErrors, { [event.target.name]: null }))
    }

    const handleChangePhoneNumber = (value) => {
        setForm(Object.assign({}, form, { phone_number: value }))
        setFormErrors(Object.assign({}, formErrors, { phone_number: null }))
    }

    const handleSubmit = async () => {
        const errors = hasFieldsErrors(form, FORM_VALIDATORS)
        if (isObjNotEmpty(errors)) {
            setFormErrors(errors)
            return
        }
        let apiValues = {
            ...form,
            birth_date: new Date(form.birth_date)
        }
        setLoading(true)
        const response = await registerAPI(apiValues)
        setLoading(false)
        if (response.errors.length) {
            return enqueueSnackbar(getAPIError(response.errors), { variant: "error" })
        }
        const loginForm = {
            email: form.email,
            password: form.password
        }
        dispatch(login(loginForm))
    };

    React.useEffect(() => {
        if (isLogged) {
            navigate("/checkout/1")
        }
    }, [isLogged])


    return (
        <ThemeProvider theme={theme}>
            <Container component="main" maxWidth="xs">
                <CssBaseline />
                <Box
                    sx={{
                        marginTop: 8,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                    }}
                >
                    <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
                        <LockOutlinedIcon />
                    </Avatar>
                    <Typography component="h1" variant="h5">
                        Registro
                    </Typography>
                    <Box sx={{ mt: 3 }}>
                        <Grid container spacing={2}>
                            <Grid item xs={2}>
                                <TextField
                                    name="document_type"
                                    value={form.document_type}
                                    onChange={handleChangeForm}
                                    select
                                    fullWidth
                                    variant="standard"
                                    margin="normal"
                                >
                                    {DEFAULT_TYPES.map((option) => (
                                        <MenuItem key={option} value={option}>
                                            {option}
                                        </MenuItem>
                                    ))}
                                </TextField>
                            </Grid>
                            <Grid item xs={10}>
                                <TextField
                                    name="document_number"
                                    value={form.document_number}
                                    onChange={handleChangeForm}
                                    variant="standard"
                                    required
                                    fullWidth
                                    label="Número de documento"
                                    error={!!formErrors.document_number}
                                    helperText={formErrors.document_number}
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="first_name"
                                    value={form.first_name}
                                    onChange={handleChangeForm}
                                    variant="standard"
                                    autoComplete="given-name"
                                    required
                                    fullWidth
                                    label="Nombres"
                                    autoFocus
                                    error={!!formErrors.first_name}
                                    helperText={formErrors.first_name}
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="last_name"
                                    value={form.last_name}
                                    onChange={handleChangeForm}
                                    variant="standard"
                                    required
                                    fullWidth
                                    label="Apellidos"
                                    autoComplete="family-name"
                                    error={!!formErrors.last_name}
                                    helperText={formErrors.last_name}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    name="birth_date"
                                    value={form.birth_date}
                                    onChange={handleChangeForm}
                                    label="Fecha de nacimiento"
                                    type="date"
                                    fullWidth
                                    variant="standard"
                                    placeholder="dd-mm-yyyy"
                                    inputProps={{
                                        "placeholder": "DD MMMM YYYY"
                                    }}
                                    InputLabelProps={{
                                        shrink: true,
                                    }}
                                    error={!!formErrors.birth_date}
                                    helperText={formErrors.birth_date}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <PhoneInput
                                    international
                                    defaultCountry='VE'
                                    name="phone_number"
                                    value={form.phone_number}
                                    onChange={handleChangePhoneNumber}
                                    placeholder='Enter phone number'
                                    inputComponent={CustomPhoneNumber}
                                    error={!!formErrors.phone_number}
                                    helperText={formErrors.phone_number}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    name="email"
                                    value={form.email}
                                    onChange={handleChangeForm}
                                    variant="standard"
                                    required
                                    fullWidth
                                    id="email"
                                    label="Correo"
                                    autoComplete="email"
                                    error={!!formErrors.email}
                                    helperText={formErrors.email}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    name="password"
                                    value={form.password}
                                    onChange={handleChangeForm}
                                    variant="standard"
                                    required
                                    fullWidth
                                    label="Contraseña"
                                    type="password"
                                    id="password"
                                    autoComplete="new-password"
                                    error={!!formErrors.password}
                                    helperText={formErrors.password}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    name="repeatPassword"
                                    value={form.repeatPassword}
                                    onChange={handleChangeForm}
                                    variant="standard"
                                    required
                                    fullWidth
                                    label="Repetir contraseña"
                                    type="password"
                                    id="repeatPassword"
                                    error={!!formErrors.repeatPassword}
                                    helperText={formErrors.repeatPassword}
                                />
                            </Grid>
                        </Grid>
                        <LoadingButton
                            type="submit"
                            fullWidth
                            loading={loading}
                            variant="contained"
                            sx={{ mt: 3, mb: 2 }}
                            onClick={handleSubmit}
                        >
                            Registrarse
                        </LoadingButton>
                        <Grid container justifyContent="flex-end">
                            <Grid item>
                                <Link component={RouterLink} to="/login" variant="body2">
                                    ¿Ya tienes una cuenta? Ingresa
                                </Link>
                            </Grid>
                        </Grid>
                    </Box>
                </Box>
            </Container>
        </ThemeProvider>
    );
}
