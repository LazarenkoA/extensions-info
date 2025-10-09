import React from "react";
import { useMutation, useQueryClient  } from '@tanstack/react-query';
import {addBaseSettings} from "../services/settingservice";

const DatabaseFormModal = () => {
    const queryClient = useQueryClient();
    const {useState, useEffect, useRef} = React;
    const [formData, setFormData] = useState({
        name: '',
        connectionString: '',
        username: '',
        password: ''
    });
    const [isOpen, setShowModal] = useState(false);
    const [error, setError] = useState('');

    const addDatabase = useMutation({
        mutationFn: addBaseSettings,
        onSuccess: () => {
            setFormData({ name: '', connectionString: '', username: '', password: '' });
            setError('');
            setShowModal(false);

            queryClient.invalidateQueries({ // для перересовки компонентов
                queryKey: ['databases']
            }).then(r => {});
        },
        onError: (err) => {
            console.log(err)

            if (err instanceof Error) {
                setError('Request error: '+ err.message);
            } else {
                setError('Не удалось добавить базу');
            }
        }
    });

    const handleSubmit = (e) => {
        e.preventDefault();
        if (!formData.name.trim() || !formData.connectionString.trim() ) {
            setError('Пожалуйста, заполните все обязательные поля');
            return;
        }

        addDatabase.mutate({
            name: formData.name,
            connectionString: formData.connectionString,
            username: formData.username,
            password: formData.password,
        });
    };
    const handleClose = () => {
        setFormData({ name: '', server: '', database: '', username: '', password: '' });
        setError('');
        setShowModal(false);
    };

    return (
        <div>
            <button
                className="add-database-btn"
                onClick={() => {setShowModal(true) }} >
                <span>+</span>
                Добавить базу
            </button>

            {isOpen && (
                <div className="modal">
                    <div className="modal-content">
                    <div className="modal-header">
                        <h2 className="modal-title">Добавить базу данных</h2>
                            <button className="close-btn" onClick={handleClose}>×</button>
                        </div>
                        <form onSubmit={handleSubmit}>
                            <div className="form-group">
                                <label className="form-label">Название базы *</label>
                                <input
                                    type="text"
                                    className="connection-input"
                                    placeholder=""
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                />
                            </div>

                            <div className="form-group">
                                <label className="form-label">Строка подключения *</label>
                                <input
                                    type="text"
                                    className="connection-input"
                                    placeholder=""
                                    value={formData.connectionString}
                                    onChange={(e) => setFormData({ ...formData, connectionString: e.target.value })}
                                />
                            </div>

                            <div className="form-group">
                                <label className="form-label">Пользователь</label>
                                <input
                                    type="text"
                                    className="connection-input"
                                    placeholder=""
                                    value={formData.Username}
                                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                                />
                            </div>

                            <div className="form-group">
                                <label className="form-label">Пароль</label>
                                <input
                                    type="password"
                                    className="connection-input"
                                    placeholder=""
                                    value={formData.Password}
                                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                />
                            </div>

                            {error && <div style={{ color: 'var(--color-error)', marginBottom: '16px', fontSize: '14px' }}>{error}</div>}

                            <div className="form-actions">
                                <button type="button" className="btn btn--secondary" onClick={handleClose}>
                                    Отмена
                                </button>
                                <button type="submit" className="btn btn--primary">
                                    Добавить
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default DatabaseFormModal;