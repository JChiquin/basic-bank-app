import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import LoadingButton from '@mui/lab/LoadingButton';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import {
    login,
    selectIsLogged,
    selectUserErrorMessage,
    selectUserLoading,
    selectUserLogged,
} from '../redux/user/userSlice';
import { useDispatch, useSelector } from 'react-redux';
import { hasFieldsErrors, isObjNotEmpty } from '../utils/formValidation'
import { useSnackbar } from 'notistack';

const theme = createTheme();

const DEFAULT_FORM = {
    email: "",
    password: ""
}
const FORM_VALIDATORS = {
    email: ["required", "email"],
    password: ["required", { maxLength: 16 }]
}

export default function SignIn() {
    const { enqueueSnackbar } = useSnackbar();
    const [form, setForm] = React.useState(DEFAULT_FORM)
    const [formErrors, setFormErrors] = React.useState(DEFAULT_FORM)
    const isLogged = useSelector(selectIsLogged);
    const userLogger = useSelector(selectUserLogged);
    const userErrorMessage = useSelector(selectUserErrorMessage);
    const userUserLoading = useSelector(selectUserLoading);
    const dispatch = useDispatch()
    const navigate = useNavigate()

    const handleChangeForm = (event) => {
        setForm(Object.assign({}, form, { [event.target.name]: event.target.value }))
        setFormErrors(Object.assign({}, formErrors, { [event.target.name]: null }))
    }

    const handleSubmit = () => {
        const errors = hasFieldsErrors(form, FORM_VALIDATORS)
        if (isObjNotEmpty(errors)) {
            setFormErrors(errors)
            return
        }
        dispatch(login(form))
    };
    
    const handleFormSubmit = (e) => {
        e.preventDefault()
        handleSubmit()
    };

    React.useEffect(() => {
        if (isLogged) {
            navigate(`/movements/${userLogger.id}`)
        }
    }, [isLogged])
    
    React.useEffect(() => {
        if (userErrorMessage) {
            enqueueSnackbar(userErrorMessage, { variant: "error" })
        }
    }, [userErrorMessage])

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
                        Ingreso
                    </Typography>
                    <Box component="form" noValidate onSubmit={handleFormSubmit}>
                        <Box sx={{ mt: 1 }}>
                            <TextField
                                margin="normal"
                                required
                                value={form.email}
                                onChange={handleChangeForm}
                                fullWidth
                                id="email"
                                label="Correo"
                                name="email"
                                autoComplete="email"
                                autoFocus
                                error={!!formErrors.email}
                                helperText={formErrors.email}
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                value={form.password}
                                onChange={handleChangeForm}
                                name="password"
                                label="Contraseña"
                                type="password"
                                id="password"
                                autoComplete="current-password"
                                error={!!formErrors.password}
                                helperText={formErrors.password}
                            />
                            <LoadingButton
                                type="submit"
                                fullWidth
                                variant="contained"
                                loading={userUserLoading}
                                sx={{ mt: 3, mb: 2 }}
                            >
                                Iniciar sesión
                            </LoadingButton>
                            <Grid container justifyContent="flex-end">
                                <Grid item>
                                    <Link component={RouterLink} to="/register" variant="body2">
                                        ¿No estás registrado? Registrate
                                    </Link>
                                </Grid>
                            </Grid>
                        </Box>
                    </Box>
                </Box>
            </Container>
        </ThemeProvider>
    );
}
