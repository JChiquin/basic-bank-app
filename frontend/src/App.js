import './styles/App.css';
import Register from './components/Register';
import Logout from './components/Logout';
import Login from './components/Login';
import Appbar from './components/Appbar';
import OrderList from './components/MovementList';
import { Routes, Route } from 'react-router-dom';


function App() {
  return (
    <div className="App">
      <Appbar />
        <Routes>
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route path="/logout" element={<Logout />} />
          <Route path="/movements/:userID" element={
              <OrderList />
          } />
          <Route path="*" element={<Login/>} />
        </Routes>
    </div>
  );
}

export default App;
