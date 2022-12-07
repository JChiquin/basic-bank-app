import * as React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TablePagination from '@mui/material/TablePagination';
import TableRow from '@mui/material/TableRow';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import Card from '@mui/material/Card';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import { ccyFormat, getAPIError } from '../utils/helpers'
import { getMovementsAPI } from '../api/modules/movement'
import { useSnackbar } from 'notistack';
import { red, green } from '@mui/material/colors';

import {
    selectUserLogged
} from '../redux/user/userSlice';
import { useSelector } from 'react-redux';

const headCells = [
    {
        id: 'id',
        numeric: false,
        disablePadding: false,
        label: 'Nro de Movimiento',
    },
    {
        id: 'created_at',
        numeric: false,
        disablePadding: false,
        label: 'Fecha',
    },
    {
        id: 'amount',
        numeric: true,
        disablePadding: false,
        label: 'Monto',
    },
    {
        id: 'balance',
        numeric: true,
        disablePadding: false,
        label: 'Saldo',
    },
];

function EnhancedTableHead(props) {
    return (
        <TableHead>
            <TableRow>
                {headCells.map((headCell) => (
                    <TableCell
                        key={headCell.id}
                        align={headCell.numeric ? 'right' : 'left'}
                        padding={headCell.disablePadding ? 'none' : 'normal'}

                    >
                        {headCell.label}
                    </TableCell>
                ))}
            </TableRow>
        </TableHead>
    );
}


const formatDate = (date) => {
    return new Date(date).toLocaleString("es")
}

const printAmount = (amount, multiplier) => {
    const color = multiplier == -1 ? red[500] : green[500]
    const arrow = multiplier == -1 ? "↓"  : "↑"
    return <Typography sx={{ color }}>
        $ {ccyFormat(amount)} {arrow}
    </Typography>
}

export default function MovementList() {
    const { enqueueSnackbar } = useSnackbar();
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(10);
    const [totalCount, setTotalCount] = React.useState(0);

    const userLogged = useSelector(selectUserLogged);

    const [movements, setMovements] = React.useState([]);

    const getMovements = async () => {
        let apiParams = {
            page_size: rowsPerPage,
            page: page + 1
        }
        const response = await getMovementsAPI(apiParams)
        if (response?.errors?.length) {
            return enqueueSnackbar(getAPIError(response.errors), { variant: "error" })
        }
        setMovements(response.data)
        setTotalCount(Number(response.headers["x-pagination-total-count"]))
    }

    React.useEffect(() => {
        getMovements()
    }, [page, rowsPerPage])

    const handleChangePage = (event, newPage) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };

    return (
        <Container component="main" maxWidth="lg" sx={{ mt: 4 }}>
            <Typography variant="h5">
                Hola, {userLogged.first_name} {userLogged.last_name}. Tu saldo es {ccyFormat(movements[0]?.balance || 0)}
            </Typography>
            <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
                <Paper sx={{ width: '100%', mb: 2 }}>
                    <TableContainer>
                        <Table
                            sx={{ minWidth: 750 }}
                            aria-labelledby="tableTitle"
                            size={'medium'}
                        >
                            <EnhancedTableHead
                                rowCount={totalCount}
                            />
                            <TableBody>
                                {movements.map((row, index) => {
                                    const labelId = `enhanced-table-checkbox-${index}`;
                                    return (
                                        <TableRow
                                            hover
                                            role="checkbox"
                                            tabIndex={-1}
                                            key={row.id}
                                        >
                                            <TableCell
                                                component="th"
                                                id={labelId}
                                                scope="row"
                                            >
                                                {row.id}
                                            </TableCell>
                                            <TableCell align="left">{formatDate(row.created_at)}</TableCell>
                                            <TableCell align="right">{printAmount(row.amount, row.multiplier)}</TableCell>
                                            <TableCell align="right">{ccyFormat(row.balance)}</TableCell>
                                        </TableRow>
                                    );
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>
                    <TablePagination
                        rowsPerPageOptions={[5, 10, 25]}
                        component="div"
                        count={totalCount}
                        rowsPerPage={rowsPerPage}
                        page={page}
                        onPageChange={handleChangePage}
                        onRowsPerPageChange={handleChangeRowsPerPage}
                    />
                </Paper>
            </Box>
            <Box sx={{ display: { xs: 'block', md: 'none' } }}>
                {movements.map((row, i) => (
                    <Card key={i} sx={{ p: 1, textAlign: 'left', mb: 1 }}>
                        <Grid container spacing={2}>
                            <Grid item xs={8}>
                                <Typography>
                                    <strong>Nro de Movimiento: </strong>{row.id}
                                </Typography>
                                <Typography>
                                    <strong>Fecha: </strong> {formatDate(row.created_at)}
                                </Typography>
                                <Typography>
                                    <strong>Monto: </strong> {printAmount(row.amount, row.multiplier)}
                                </Typography>
                                <Typography>
                                    <strong>Balance: </strong> {ccyFormat(row.balance)}
                                </Typography>
                            </Grid>
                        </Grid>
                    </Card>
                ))}
            </Box>
        </Container>
    );
}
