import { makeAutoObservable } from "mobx";

interface AuthData {
    username: string;
    expiry: number;
}

class AuthStore {
    user: string | null = null;

    constructor() {
        makeAutoObservable(this);
        this.loadFromStorage();
    }

    // Загружаем статус авторизации из localStorage
    private loadFromStorage() {
        const data = localStorage.getItem("auth");
        if (!data) return;

        const parsed: AuthData = JSON.parse(data);
        if (Date.now() < parsed.expiry) {
            this.user = parsed.username;
        } else {
            localStorage.removeItem("auth");
            this.user = null;
        }
    }

    private saveToStorage(username: string) {
        const expiry = Date.now() + 2 * 60 * 60 * 1000; // 2 часа
        const data: AuthData = { username, expiry };
        localStorage.setItem("auth", JSON.stringify(data));
    }

    // Логин через бэкенд, cookie ставится автоматически
    async login(username: string, password: string) {
        const res = await fetch("/api/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password }),
            credentials: "include", // <- отправка cookie
        });

        if (!res.ok) throw new Error("Неверный логин или пароль");

        this.user = username;
        this.saveToStorage(username);
    }

    // Регистрация
    async register(username: string, password: string) {
        const res = await fetch("/api/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password }),
            credentials: "include",
        });

        if (!res.ok) throw new Error("Не удалось зарегистрироваться");

        this.user = username;
        this.saveToStorage(username);
    }

    // Логаут без запроса к бэку
    logout() {
        this.user = null;
        localStorage.removeItem("auth");
    }

    get isAuthenticated() {
        return !!this.user;
    }
}

export const authStore = new AuthStore();
