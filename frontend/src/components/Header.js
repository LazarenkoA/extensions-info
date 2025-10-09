import React, { useState, useEffect } from 'react';
import SettingsModal from './SettingsModal'

// Theme Toggle Component
const ThemeToggle = ({ theme, onToggle }) => {
    const icon = theme === 'light' ? '🌙' : '☀️';

    return (
        <button
            className="theme-toggle"
            onClick={onToggle}
            title={theme === 'light' ? 'Переключить на тёмную тему' : 'Переключить на светлую тему'}
        >
            <span className="theme-icon">{icon}</span>
        </button>
    );
};


const Header = () => {
    const [theme, setTheme] = useState('light');
    const [showSettingsModal, setShowSettingsModal] = useState(false);

    // Initialize theme
    useEffect(() => {
        // Simulate theme persistence (would use localStorage in real app)
        const savedTheme = 'light'; // localStorage.getItem('theme') || 'light';
        setTheme(savedTheme);
        document.documentElement.setAttribute('data-color-scheme', savedTheme);
    }, []);

    const toggleTheme = () => {
        const newTheme = theme === 'light' ? 'dark' : 'light';
        setTheme(newTheme);
        document.documentElement.setAttribute('data-color-scheme', newTheme);
    };


    return (
        <header className="header">
            <div className="container">
                <div className="header-content">
                    <h1 className="logo">Анализатор расширений 1С</h1>
                    <div className="header-actions">
                        <button className="icon-button" title="Глобальные настройки"
                                onClick={() => setShowSettingsModal(true)}
                        >⚙️</button>
                        <div className="header-right">
                            <ThemeToggle theme={theme} onToggle={toggleTheme}/>
                        </div>
                    </div>

                    {showSettingsModal && (
                        <SettingsModal onClose={() =>  setShowSettingsModal(false)} />
                    )}
                </div>
            </div>
        </header>
    )
}

export default Header;