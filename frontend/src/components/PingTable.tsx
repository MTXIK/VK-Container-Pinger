import React from 'react';
import { Table } from 'antd';
import { Ping } from '../hooks/usePings';

interface PingTableProps {
  pings: Ping[];
}

const PingTable: React.FC<PingTableProps> = ({ pings }) => {
  const columns = [
    {
      title: 'IP адрес',
      dataIndex: 'ip_address',
      key: 'ip_address',
    },
    {
      title: 'Имя контейнера',
      dataIndex: 'container_name',
      key: 'container_name',
    },
    {
      title: 'Время пинга (мс)',
      dataIndex: 'ping_time',
      key: 'ping_time',
      render: (value: number) => value.toFixed(2),
    },
    {
      title: 'Дата последнего пинга',
      dataIndex: 'last_success',
      key: 'last_success',
      render: (text: string) => new Date(text).toLocaleString(),
    },
  ];

  return <Table rowKey="id" dataSource={pings} columns={columns} pagination={{ pageSize: 10 }} />;
};

export default PingTable;
