import React from 'react';
import { Spin, Alert } from 'antd';
import 'antd/dist/reset.css';
import PingTable from './components/PingTable';
import { usePings } from './hooks/usePings';

const App: React.FC = () => {
  const { pings, loading, error } = usePings(10000);

  return (
    <div style={{ padding: '20px' }}>
      <h1>Статус контейнеров</h1>
      {error && <Alert message={error} type="error" showIcon style={{ marginBottom: '20px' }} />}
      {loading ? <Spin tip="Загрузка..." /> : <PingTable pings={pings} />}
    </div>
  );
};

export default App;
