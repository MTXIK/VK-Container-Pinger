import React, { useState } from 'react';
import { ConfigProvider, Switch, Spin, Alert, theme } from 'antd';
import PingTable from './components/PingTable';
import { usePings } from './hooks/usePings';
import 'antd/dist/reset.css';

const App: React.FC = () => {
  const { pings, loading, error } = usePings(10000);
  const [isDark, setIsDark] = useState<boolean>(true);

  return (
    <ConfigProvider theme={{ algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm }}>
      <div style={{ padding: '20px', backgroundColor: isDark ? '#141414' : '#fff', minHeight: '100vh' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h1 style={{ color: isDark ? '#fff' : '#000' }}>Статус контейнеров</h1>
          <Switch
            checkedChildren="Dark"
            unCheckedChildren="Light"
            checked={isDark}
            onChange={setIsDark}
          />
        </div>
        {error && <Alert message={error} type="error" showIcon style={{ marginBottom: '20px' }} />}
        {loading ? <Spin tip="Загрузка..." /> : <PingTable pings={pings} />}
      </div>
    </ConfigProvider>
  );
};

export default App;