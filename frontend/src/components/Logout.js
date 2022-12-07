import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import {
    logout,
    selectIsLogged
} from '../redux/user/userSlice';
import { useDispatch, useSelector } from 'react-redux';


export default function SignOut() {
    const isLogged = useSelector(selectIsLogged);
    const dispatch = useDispatch()
    const navigate = useNavigate()

    React.useEffect(() => {
        dispatch(logout())
      }, [])

    React.useEffect(() => {
        if (!isLogged) {
            navigate("/login")
        }
      }, [isLogged])

      return null
}
