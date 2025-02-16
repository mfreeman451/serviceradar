import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { AlertCircle } from 'lucide-react';

const SNMPDashboard = ({ nodeId, service }) => {
    const [selectedOid, setSelectedOid] = useState(null);
    const [timeRange, setTimeRange] = useState('1h');
    const [metrics, setMetrics] = useState([]);

    useEffect(() => {
        const fetchMetrics = async () => {
            if (!selectedOid) return;

            try {
                const response = await fetch(`/api/nodes/${nodeId}/metrics/snmp?oid=${selectedOid}&range=${timeRange}`);
                if (!response.ok) throw new Error('Failed to fetch metrics');
                const data = await response.json();
                setMetrics(data);
            } catch (error) {
                console.error('Error fetching SNMP metrics:', error);
            }
        };

        fetchMetrics();
        const interval = setInterval(fetchMetrics, 10000);
        return () => clearInterval(interval);
    }, [nodeId, selectedOid, timeRange]);

    const formatMetricValue = (value, type) => {
        switch (type) {
            case 'numeric':
                return typeof value === 'number' ? value.toLocaleString() : value;
            case 'boolean':
                return value ? 'Yes' : 'No';
            case 'counter':
                return typeof value === 'number' ? value.toLocaleString() : value;
            case 'gauge':
                return typeof value === 'number' ? value.toFixed(2) : value;
            case 'string':
                return String(value);
            default:
                return value;
        }
    };

    return (
        <div className="space-y-6">
            {/* Time range selector */}
            <div className="flex gap-2">
                {['1h', '6h', '24h', '7d'].map((range) => (
                    <button
                        key={range}
                        onClick={() => setTimeRange(range)}
                        className={`px-3 py-1 rounded transition-colors ${
                            timeRange === range
                                ? 'bg-blue-500 text-white'
                                : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200'
                        }`}
                    >
                        {range}
                    </button>
                ))}
            </div>

            {/* Metrics chart */}
            {selectedOid && metrics.length > 0 && (
                <div className="h-64">
                    <ResponsiveContainer width="100%" height="100%">
                        <LineChart data={metrics}>
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis
                                dataKey="timestamp"
                                type="number"
                                domain={['auto', 'auto']}
                                tickFormatter={(ts) => new Date(ts).toLocaleString()}
                            />
                            <YAxis
                                domain={['auto', 'auto']}
                                hide={!['numeric', 'counter', 'gauge'].includes(metrics[0]?.value_type)}
                            />
                            <Tooltip
                                labelFormatter={(ts) => new Date(ts).toLocaleString()}
                                formatter={(value, name, props) => [
                                    formatMetricValue(value, props.payload.value_type),
                                    selectedOid
                                ]}
                            />
                            <Legend />
                            {['numeric', 'counter', 'gauge'].includes(metrics[0]?.value_type) && (
                                <Line
                                    type="monotone"
                                    dataKey="value"
                                    stroke="#8884d8"
                                    dot={false}
                                    name={selectedOid}
                                />
                            )}
                        </LineChart>
                    </ResponsiveContainer>
                </div>
            )}

            {/* OID list */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {service.details && service.details.oid_status &&
                    Object.entries(service.details.oid_status).map(([oidName, status]) => (
                        <div
                            key={oidName}
                            className={`bg-white dark:bg-gray-800 rounded-lg shadow p-4 transition-colors 
                                cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700
                                ${selectedOid === oidName ? 'ring-2 ring-blue-500' : ''}`}
                            onClick={() => setSelectedOid(oidName)}
                        >
                            <div className="flex justify-between items-start">
                                <div>
                                    <h4 className="font-medium text-gray-800 dark:text-gray-200">{oidName}</h4>
                                    <div className="text-2xl font-bold mt-1 text-gray-900 dark:text-gray-100">
                                        {formatMetricValue(status.last_value, status.type)}
                                    </div>
                                </div>

                                {status.error_count > 0 && (
                                    <AlertCircle className="text-red-500 dark:text-red-400" size={20} />
                                )}
                            </div>
                        </div>
                    ))
                }
            </div>
        </div>
    );
};

export default SNMPDashboard;