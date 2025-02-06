import React from 'react';
import { Table, Tag } from 'antd';
import { Ping } from '../hooks/usePings';

interface PingTableProps {
  pings: Ping[];
}

const PingTable: React.FC<PingTableProps> = ({ pings }) => {
  const pollInterval = 10000; // 10 секунд
  const threshold = 15000;     // 15 секунд

  const columns = [
    {
      title: 'IP адрес',
      dataIndex: 'ip_address',
      key: 'ip_address',
      sorter: (a: Ping, b: Ping) => a.ip_address.localeCompare(b.ip_address),
    },
    {
      title: 'Имя контейнера',
      dataIndex: 'container_name',
      key: 'container_name',
      sorter: (a: Ping, b: Ping) => a.container_name.localeCompare(b.container_name),
    },
    {
      title: 'Время пинга (мс)',
      dataIndex: 'ping_time',
      key: 'ping_time',
      sorter: (a: Ping, b: Ping) => a.ping_time - b.ping_time,
      render: (value: number) => value.toFixed(2),
    },
    {
      title: 'Дата последнего пинга',
      dataIndex: 'last_success',
      key: 'last_success',
      sorter: (a: Ping, b: Ping) =>
        new Date(a.last_success).getTime() - new Date(b.last_success).getTime(),
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: 'Статус',
      key: 'status',
      render: (_: any, record: Ping) => {
        const lastPing = new Date(record.last_success).getTime();
        const diff = Date.now() - lastPing;
        const isDown = diff > pollInterval + threshold;
        return <Tag color={isDown ? 'red' : 'green'}>{isDown ? 'Probably Down' : 'Probably Up'}</Tag>;
      },
    },
  ];

  return <Table rowKey="id" dataSource={pings} columns={columns} pagination={{ pageSize: 10 }} />;
};

export default PingTable;