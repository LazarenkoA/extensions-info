import React from "react";
import DatabaseList from '../components/DatabaseItemComponent';
import AddDatabaseBtn from '../components/DatabaseFormModal';
import MainContent from  '../components/MainContent';
import { DatabaseProvider } from "../components/ContexDatabaset";
import {WebSocketProvider} from '../services/contexWebSocket';
import './Home.css';
import './Icon.css'

const Home = () => {
  return (
      <WebSocketProvider>
        <DatabaseProvider>
          <div className="app">
            <div className="main-layout">
              <aside className="sidebar">
                <div className="sidebar-header">
                  <h2 className="sidebar-title">Базы данных</h2>
                  <AddDatabaseBtn/>
                </div>
                <div className="sidebar-content">
                  <div className="database-list">
                    <DatabaseList />
                  </div>
                </div>
              </aside>
              <main className="main-content">
                <MainContent />
              </main>
            </div>
          </div>
        </DatabaseProvider>
      </WebSocketProvider>
  );
};

export default Home;