import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Navbar from './components/navbar/Navbar';
import Flows from './components/flow/Flows';
import WorkflowBuilder from './components/flow/flow-builder/FlowBuilder';

function App() {
  return (
    <BrowserRouter>
        <Navbar />
        <Routes>
            <Route path="/flow" element={<Flows />} />
            <Route path="/flow/new" element={<WorkflowBuilder />} />
            <Route path="/flow/:id" element={<WorkflowBuilder />} />
        </Routes>
    </BrowserRouter>
  );
}

export default App;
