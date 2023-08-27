import { Navigate, Outlet } from "react-router-dom";
import { checkLoginState } from "./auth-utils";
const ProtectedRoutes = () => {
    const loggedIn = checkLoginState();
    return (loggedIn ? <Outlet /> : <Navigate to="/login" replace/>)
}
export default ProtectedRoutes;