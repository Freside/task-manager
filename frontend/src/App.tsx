import { BrowserRouter, Routes, Route } from "react-router-dom";
import Register from "./pages/Register";
import Login from "./pages/Login";
import Boards from "./pages/Boards";
import BoardDetail from "./pages/BoardDetail.tsx";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/register" element={<Register />} />
                <Route path="/login" element={<Login />} />
                <Route path="/boards" element={<Boards />} />
                <Route path="/boards/:boardId" element={<BoardDetail />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;
