import React, { useState, useEffect } from 'react';
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
} from 'recharts';

const NodeTimeline = ({ nodeId }) => {
    const [availabilityData, setAvailabilityData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(`/api/nodes/${nodeId}/history`);
                if (!response.ok) throw new Error('Failed to fetch node history');
                const data = await response.json();

                // Transform the history data for the chart
                const timelineData = data.map((point) => ({
                    timestamp: new Date(point.timestamp).getTime(),
                    status: point.is_healthy ? 1 : 0,
                    tooltipTime: new Date(point.timestamp).toLocaleString(),
                }));

                setAvailabilityData(timelineData);
                setLoading(false);
            } catch (err) {
                console.error('Error fetching history:', err);
                setError(err.message);
                setLoading(false);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 10000);
        return () => clearInterval(interval);
    }, [nodeId]);

    const CustomTooltip = ({ active, payload }) => {
        if (!active || !payload || !payload.length) return null;

        const data = payload[0].payload;
        return (
            <div className="bg-white dark:bg-gray-700 p-4 rounded shadow-lg border dark:border-gray-600 dark:text-gray-100">
                <p className="text-sm font-semibold">{data.tooltipTime}</p>
                <p className="text-sm">
                    Status: {data.status === 1 ? 'Online' : 'Offline'}
                </p>
            </div>
        );
    };

    if (loading && availabilityData.length === 0) {
        return <div className="text-center p-4">Loading timeline...</div>;
    }

    if (error && availabilityData.length === 0) {
        return <div className="text-red-500 text-center p-4">{error}</div>;
    }

    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4 transition-colors">
            <h3 className="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-100">
                Node Availability Timeline
            </h3>
            <div className="h-48">
                <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={availabilityData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis
                            dataKey="timestamp"
                            type="number"
                            domain={['auto', 'auto']}
                            tickFormatter={(ts) => new Date(ts).toLocaleTimeString()}
                        />
                        <YAxis
                            domain={[0, 1]}
                            ticks={[0, 1]}
                            tickFormatter={(value) => (value === 1 ? 'Online' : 'Offline')}
                        />
                        <Tooltip content={<CustomTooltip />} />
                        <Line
                            type="stepAfter"
                            dataKey="status"
                            stroke="#8884d8"
                            dot={false}
                            isAnimationActive={false}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
};

export default NodeTimeline;
