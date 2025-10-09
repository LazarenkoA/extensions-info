
import React, { createContext, useContext, useEffect, useRef, useState } from 'react';

const WS_URL = process.env.REACT_WS_URL || 'http://localhost:8080/api/v1/ws';
const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {
    const wsRef = useRef(null);
    const heartbeatRef = useRef(null);
    const listeners  = useRef(new Set());

    const openWSConn = () => {
        const ws = new WebSocket(WS_URL); // твой эндпоинт
        wsRef.current = ws;

        ws.onopen = () => {
            console.log("✅ WS connected");
            // запускаем keepalive
            clearInterval(heartbeatRef.current); // очистим старый
            heartbeatRef.current = setInterval(() => {
                if (ws.readyState === WebSocket.OPEN) {
                    ws.send(JSON.stringify({ type: "ping" }));
                }
                if (ws.readyState === WebSocket.CLOSED) {
                    console.log("🔌 write to closed WS");
                    openWSConn();
                }
            }, 5_000); // каждые 5 сек
        };

        ws.onmessage = (event) => {
            listeners.current.forEach((h) => h(event.data)); // уведомляем подписчиков
        };

        ws.onerror = (err) => {
            console.error("❌ WS error", err);
        };
    }


    useEffect(() => {
       openWSConn()
        if(wsRef.current) {
            wsRef.current.onclose = () => {
                console.log("🔌 WS closed");
            };
        }

       return () => {
          wsRef.current?.close();
          clearInterval(heartbeatRef.current);
       };
    }, []);

    const sendMessage = (msg) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify(msg));
        }
    };

    const subscribe = (callback) => {
        listeners.current.add(callback);
        return () => listeners.current.delete(callback);
    };

    return (
        <WebSocketContext.Provider value={{ subscribe }}>
            {children}
        </WebSocketContext.Provider>
    );
};

export const useWebSocket = () => useContext(WebSocketContext);
