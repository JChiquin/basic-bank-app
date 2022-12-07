import * as React from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import Menu from '@mui/material/Menu';
import MenuIcon from '@mui/icons-material/Menu';
import Container from '@mui/material/Container';
import Button from '@mui/material/Button';
import MenuItem from '@mui/material/MenuItem';
import logo from '../static/images/logo.png'
import { Link } from "react-router-dom";
import {
    selectIsLogged
} from '../redux/user/userSlice';
import { useSelector } from 'react-redux';

const guestPages = [
    {
        name: 'Registrarse',
        link: '/register'
    },
    {
        name: 'Ingresar',
        link: '/login'
    }
];
const userPages = [
    {
        name: 'Movimientos',
        link: '/movements'
    },
    {
        name: 'Cerrar sesión',
        link: '/logout'
    },
];

const ResponsiveAppBar = () => {
    const [anchorElNav, setAnchorElNav] = React.useState(null);
    const [pages, setPages] = React.useState([]);
    const isLogged = useSelector(selectIsLogged);

    React.useEffect(() => {
        if (isLogged) {
            setPages(userPages)
        } else {
            setPages(guestPages)
        }
    }, [isLogged])

    const handleOpenNavMenu = (event) => {
        setAnchorElNav(event.currentTarget);
    };

    const handleCloseNavMenu = () => {
        setAnchorElNav(null);
    };


    return (
        <AppBar elevation={1} position="static" sx={{ background: "#eef7f1", py: 1 }}>
            <Container maxWidth="xl">
                <Toolbar disableGutters variant="regular">
                    <img src={logo} alt="bank" width="90" />
                    <Box sx={{ flexGrow: 1 }} />
                    <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
                        {pages.map((page) => (
                            <Button
                                to={page.link}
                                component={Link}
                                key={page.name}
                                sx={{ my: 2, color: 'black', display: 'block', textTransform: 'none' }}
                            >
                                {page.name}
                            </Button>
                        ))}
                    </Box>
                    <Box sx={{ display: { xs: 'flex', md: 'none' } }}>
                        <IconButton
                            size="large"
                            aria-label="account of current user"
                            aria-controls="menu-appbar"
                            aria-haspopup="true"
                            onClick={handleOpenNavMenu}

                        >
                            <MenuIcon />
                        </IconButton>
                        <Menu
                            id="menu-appbar"
                            anchorEl={anchorElNav}
                            anchorOrigin={{
                                vertical: 'bottom',
                                horizontal: 'left',
                            }}
                            keepMounted
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'left',
                            }}
                            open={Boolean(anchorElNav)}
                            onClose={handleCloseNavMenu}
                            sx={{
                                display: { xs: 'block', md: 'none' },
                            }}
                        >
                            {pages.map((page) => (
                                <Link to={page.link} key={page.name}>
                                    <MenuItem onClick={handleCloseNavMenu}>
                                        <Typography textAlign="center">{page.name}</Typography>
                                    </MenuItem>
                                </Link>
                            ))}
                        </Menu>
                    </Box>
                </Toolbar>
            </Container>
        </AppBar>
    );
};
export default ResponsiveAppBar;
