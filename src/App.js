import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Navbar from './components/navbar/Navbar';
import Flows from './components/flow/Flows';
import WorkflowBuilder from './components/flow/flow-builder/FlowBuilder';
import { PATHS } from './constants/constants';
import ProtectedRoutes from './components/auth/ProtectedRoutes';
import Login from './components/auth/Login';

function App() {
  return (
    <BrowserRouter>
        <Navbar />
        <Routes>
            <Route path="/login" element={<Login />} />
            <Route element={<ProtectedRoutes />}>
                <Route path={PATHS.FLOWS} element={<Flows />} />
                <Route path={PATHS.ADD_FLOW} element={<WorkflowBuilder />} />
                <Route path={PATHS.OPEN_FLOW} element={<WorkflowBuilder />} />
            </Route>
            
        </Routes>
    </BrowserRouter>
  );
}

export default App;
