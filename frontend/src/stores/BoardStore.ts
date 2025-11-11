import { makeAutoObservable } from "mobx";

class BoardStore {
    boards: any[] = [];
    columns: any[] = [];
    cards: Record<number, any[]> = {}; // ключ = columnId
    currentBoardId: number | null = null;

    constructor() {
        makeAutoObservable(this);
    }

    // Получение досок
    async fetchBoards() {
        const res = await fetch("/api/boards", { credentials: "include" });
        if (!res.ok) throw new Error("Ошибка при загрузке досок");
        this.boards = await res.json();
    }

    // Создание доски
    async createBoard(title: string) {
        const res = await fetch("/api/boards", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title }),
            credentials: "include",
        });
        if (!res.ok) throw new Error("Не удалось создать доску");
        const board = await res.json();
        this.boards.push(board);
    }

    // Получение колонок по доске
    async fetchColumns(boardId: number) {
        const res = await fetch(`/api/boards/${boardId}/columns`, { credentials: "include" });
        if (!res.ok) throw new Error("Ошибка при загрузке колонок");
        this.columns = await res.json();
        this.currentBoardId = boardId;
    }

    // Создание колонки
    async createColumn(title: string) {
        if (!this.currentBoardId) throw new Error("Не выбрана доска");
        const res = await fetch("/api/columns", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, board_id: this.currentBoardId }),
            credentials: "include",
        });
        if (!res.ok) throw new Error("Не удалось создать колонку");
        const column = await res.json();
        this.columns.push(column);
    }

    // Получение карточек по колонке
    async fetchCards(columnId: number) {
        const res = await fetch(`/api/columns/${columnId}/cards`, { credentials: "include" });
        if (!res.ok) throw new Error("Ошибка при загрузке карточек");
        this.cards[columnId] = await res.json() || [];
    }

    // Создание карточки
    async createCard(title: string, columnId: number) {
        const res = await fetch("/api/cards", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, column_id: columnId }),
            credentials: "include",
        });
        if (!res.ok) throw new Error("Не удалось создать карточку");
        const card = await res.json();
        if (!this.cards[columnId]) this.cards[columnId] = [];
        this.cards[columnId].push(card);
    }
}

export const boardStore = new BoardStore();
