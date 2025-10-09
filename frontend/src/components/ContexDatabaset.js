import React, { createContext, useState, useContext } from "react";

// создаём контекст
const ContexDatabaset = createContext(null);

// провайдер, который оборачивает все дочерние компоненты
export const DatabaseProvider = ({ children }) => {
    const [selectedDb, setSelectedDb] = useState(null);

    return (
        <ContexDatabaset.Provider value={{ selectedDb, setSelectedDb }}>
            {children}
        </ContexDatabaset.Provider>
    );
};

// хук для удобного использования
export const useDatabase = () => useContext(ContexDatabaset);