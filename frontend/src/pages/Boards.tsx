import { useEffect, useState } from "react";
import { observer } from "mobx-react-lite";
import { boardStore } from "../stores/BoardStore";
import { authStore } from "../stores/AuthStore";
import {Button, Input, Card, List, message} from "antd";
import { useNavigate } from "react-router-dom";

const Boards = observer(() => {
    const [newBoardTitle, setNewBoardTitle] = useState("");
    const navigate = useNavigate();

    useEffect(() => {
        if (!authStore.isAuthenticated) {
            navigate("/login");
            return;
        }

        boardStore.fetchBoards().catch((err) => message.error(err.message));
    }, []);

    const handleCreateBoard = async () => {
        if (!newBoardTitle.trim()) return;

        try {
            await boardStore.createBoard(newBoardTitle.trim());
            message.success("Доска создана!");
            setNewBoardTitle("");
        } catch (err: any) {
            message.error(err.message);
        }
    };

    const handleOpenBoard = (boardId: number) => {
        navigate(`/boards/${boardId}`); // позже сделаем BoardDetail
    };

    return (
        <div className="p-8 bg-gray-100 min-h-screen">
            <div className="flex justify-between mb-6">
                <h1 className="text-2xl font-bold">Мои доски</h1>
                <Button
                    danger
                    onClick={() => {
                        authStore.logout();       // сброс авторизации
                        navigate("/login");       // редирект на страницу логина
                    }}
                >
                    Выйти
                </Button>
            </div>

            <div className="mb-6 flex gap-2">
                <Input
                    placeholder="Название новой доски"
                    value={newBoardTitle}
                    onChange={(e) => setNewBoardTitle(e.target.value)}
                />
                <Button type="primary" onClick={handleCreateBoard}>
                    Создать
                </Button>
            </div>

            <List
                grid={{ gutter: 16, column: 3 }}
                dataSource={boardStore.boards || []} // <- пустой массив по дефолту
                renderItem={(board) => board ? (
                    <List.Item key={board.id}>
                        <Card hoverable onClick={() => handleOpenBoard(board.id)} title={board.title} />
                    </List.Item>
                ) : null}
            />

        </div>
    );
});

export default Boards;
