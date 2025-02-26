import React from 'react';

// Helper functions for formatting
const formatResponseTime = (time) => {
    if (!time && time !== 0) return 'N/A';
    return `${(time / 1000000).toFixed(2)}ms`;
};

const formatPacketLoss = (loss) => {
    if (typeof loss !== 'number') return 'N/A';
    return `${loss.toFixed(1)}%`;
};

// Individual ping status component
const PingStatus = ({ details }) => {
    const getPingDetails = () => {
        try {
            return typeof details === 'string' ? JSON.parse(details) : details;
        } catch (e) {
            console.error('Error parsing ping details:', e);
            return null;
        }
    };

    const pingData = getPingDetails();

    if (!pingData) {
        return (
            <div className="text-gray-500 dark:text-gray-400 transition-colors">
                No ping data available
            </div>
        );
    }

    return (
        <div className="grid grid-cols-2 gap-2 text-sm transition-colors">
            <div className="font-medium text-gray-600 dark:text-gray-400">Response Time:</div>
            <div className="text-gray-800 dark:text-gray-100">
                {formatResponseTime(pingData.response_time)}
            </div>

            <div className="font-medium text-gray-600 dark:text-gray-400">Packet Loss:</div>
            <div className="text-gray-800 dark:text-gray-100">
                {formatPacketLoss(pingData.packet_loss)}
            </div>

            <div className="font-medium text-gray-600 dark:text-gray-400">Status:</div>
            <div
                className={`font-medium ${
                    pingData.available
                        ? 'text-green-600 dark:text-green-400'
                        : 'text-red-600 dark:text-red-400'
                }`}
            >
                {pingData.available ? 'Available' : 'Unavailable'}
            </div>
        </div>
    );
};

export { PingStatus };
