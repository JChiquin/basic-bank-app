import './styles/App.css';
import Register from './components/Register';
import Logout from './components/Logout';
import Login from './components/Login';
import Appbar from './components/Appbar';
import MovementList from './components/MovementList';
import { Routes, Route, Navigate, useLocation, useNavigate } from 'react-router-dom';
import {
  selectIsLogged,
  whoAmI
} from './redux/user/userSlice';
import { useSelector, useDispatch } from 'react-redux';
import { getJWT } from "./utils/localStorage"
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';

function RequireAuth({ children }) {
  const isLogged = useSelector(selectIsLogged);
  let location = useLocation();

  if (!isLogged) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return children;
}

function KeepLogged({ children }) {
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const isLogged = useSelector(selectIsLogged);
  const jwt = getJWT()

  if (jwt && !isLogged) { //There is a JWT but user is not in redux
    dispatch(whoAmI())
    .unwrap()
    .catch(() => {
      navigate("/login")
    })


    return <Box sx={{ display: 'flex', justifyContent: "center", mt: 5 }}>
      <CircularProgress size={100} />
    </Box>
  }
  return children
}

function NotFound() {
  const isLogged = useSelector(selectIsLogged);

  if (isLogged) {
    return <Navigate replace to="/movements" />
  }
  return <Navigate replace to="/login" />
}

function App() {
  return (
    <div className="App">
      <Appbar />
      <KeepLogged>
        <Routes>
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route path="/logout" element={<Logout />} />
          <Route path="/movements" element={
            <RequireAuth>
              <MovementList />
            </RequireAuth>
          } />
          <Route path="*" element={<NotFound/>} />
        </Routes>
      </KeepLogged>
    </div>
  );
}

export default App;
