import { useState } from "react";
import { observer } from "mobx-react-lite";
import { authStore } from "../stores/AuthStore";
import { Button, Input, Form, message, Card } from "antd";
import { useNavigate } from "react-router-dom";

const Login = observer(() => {
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const onFinish = async (values: { username: string; password: string }) => {
        setLoading(true);
        try {
            await authStore.login(values.username, values.password);
            message.success("Вход выполнен!");
            navigate("/boards"); // переход на страницу досок
        } catch (err: any) {
            message.error(err.message || "Ошибка входа");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex items-center justify-center h-screen bg-gray-100">
            <Card title="Вход" className="w-96">
                <Form layout="vertical" onFinish={onFinish}>
                    <Form.Item
                        label="Имя пользователя"
                        name="username"
                        rules={[{ required: true, message: "Введите имя пользователя" }]}
                    >
                        <Input placeholder="Username" />
                    </Form.Item>

                    <Form.Item
                        label="Пароль"
                        name="password"
                        rules={[{ required: true, message: "Введите пароль" }]}
                    >
                        <Input.Password placeholder="Password" />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" block loading={loading}>
                            Войти
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
});

export default Login;
