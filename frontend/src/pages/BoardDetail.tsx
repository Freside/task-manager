import { useEffect, useState } from "react";
import { observer } from "mobx-react-lite";
import { boardStore } from "../stores/BoardStore";
import { useParams } from "react-router-dom";
import { Card, List, Input, Button, message, Empty } from "antd";

const BoardDetail = observer(() => {
    const { boardId } = useParams<{ boardId: string }>();
    const [newColumnTitle, setNewColumnTitle] = useState("");
    const [newCardTitles, setNewCardTitles] = useState<Record<number, string>>({});

    useEffect(() => {
        if (boardId) {
            boardStore
                .fetchBoards()
                .then(async () => {
                    await boardStore.fetchColumns(Number(boardId));

                    // Безопасная проверка — если колонок нет, не вызываем fetchCards
                    if (boardStore.columns?.length) {
                        await Promise.all(
                            boardStore.columns.map(col =>
                                boardStore
                                    .fetchCards(col.id)
                                    .catch(err => message.error(err.message))
                            )
                        );
                    }
                })
                .catch(err => message.error(err.message));
        }
    }, [boardId]);

    const boardTitle =
        boardStore.boards.length > 0
            ? boardStore.boards.find(b => b.id === Number(boardId))?.title || "Доска"
            : "Загрузка...";

    const handleAddColumn = async () => {
        if (!newColumnTitle.trim()) return;
        try {
            await boardStore.createColumn(newColumnTitle.trim());
            message.success("Колонка создана!");
            setNewColumnTitle("");

            // После создания — обновляем список колонок
            await boardStore.fetchColumns(Number(boardId));
        } catch (err: any) {
            message.error(err.message);
        }
    };

    const handleAddCard = async (columnId: number) => {
        const title = newCardTitles[columnId]?.trim();
        if (!title) return;
        try {
            await boardStore.createCard(title, columnId);
            message.success("Карточка создана!");
            setNewCardTitles(prev => ({ ...prev, [columnId]: "" }));

            // После добавления карточки — обновляем карточки этой колонки
            await boardStore.fetchCards(columnId);
        } catch (err: any) {
            message.error(err.message);
        }
    };

    return (
        <div className="p-8 bg-gray-100 min-h-screen">
            <h1 className="text-2xl font-bold mb-4">{boardTitle}</h1>

            <div className="mb-4 flex gap-2">
                <Input
                    placeholder="Название новой колонки"
                    value={newColumnTitle}
                    onChange={e => setNewColumnTitle(e.target.value)}
                />
                <Button type="primary" onClick={handleAddColumn}>
                    Добавить колонку
                </Button>
            </div>

            <div className="flex gap-4 overflow-x-auto">
                {!boardStore.columns || boardStore.columns.length === 0 ? (
                    <Empty
                        description="На этой доске пока нет колонок"
                        image={Empty.PRESENTED_IMAGE_SIMPLE}
                    />
                ) : (
                    boardStore.columns.map(column => (
                        <Card
                            key={column.id}
                            title={column.title}
                            className="w-64 flex-shrink-0"
                        >
                            <List
                                dataSource={boardStore.cards[column.id] || []}
                                renderItem={card => (
                                    <List.Item key={card.id}>{card.title}</List.Item>
                                )}
                            />

                            <div className="mt-2 flex gap-1">
                                <Input
                                    placeholder="Название карточки"
                                    value={newCardTitles[column.id] || ""}
                                    onChange={e =>
                                        setNewCardTitles(prev => ({
                                            ...prev,
                                            [column.id]: e.target.value,
                                        }))
                                    }
                                />
                                <Button
                                    type="primary"
                                    onClick={() => handleAddCard(column.id)}
                                >
                                    Добавить
                                </Button>
                            </div>
                        </Card>
                    ))
                )}
            </div>
        </div>
    );
});

export default BoardDetail;