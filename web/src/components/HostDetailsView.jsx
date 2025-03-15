import { useState, useMemo, useEffect } from 'react';
import { ChevronDown, ChevronUp } from 'lucide-react';
import '../styles/NetworkSweepView.css';

// HostDetailsView subcomponent defined within NetworkSweepView.jsx
const HostDetailsView = ({ host }) => {
    const [expanded, setExpanded] = useState(false);

    const formatResponseTime = (ns) => {
        if (!ns || ns === 0) return 'N/A';
        const ms = ns / 1000000; // Convert nanoseconds to milliseconds
        return `${ms.toFixed(2)}ms`;
    };

    const toggleExpanded = () => {
        setExpanded(!expanded);
    };

    return (
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow transition-colors">
            <div className="flex justify-between items-center">
                <h4 className="text-base sm:text-lg font-semibold text-gray-800 dark:text-gray-200">
                    {host.host}
                </h4>
                <div className="flex items-center gap-2">
                    <span
                        className={`px-2 py-1 text-xs sm:text-sm rounded transition-colors ${
                            host.available
                                ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100'
                                : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-100'
                        }`}
                    >
                        {host.available ? 'Online' : 'Offline'}
                    </span>
                    <button
                        onClick={toggleExpanded}
                        className="p-1 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 sm:hidden"
                    >
                        {expanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
                    </button>
                </div>
            </div>

            <div className={`mt-3 ${expanded ? 'block' : 'hidden sm:block'}`}>
                {host.icmp_status && (
                    <div className="mt-3 bg-gray-50 dark:bg-gray-700 p-3 rounded transition-colors">
                        <h5 className="font-medium mb-2 text-gray-800 dark:text-gray-200">
                            ICMP Status
                        </h5>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-sm">
                            <div>
                                <span className="text-gray-600 dark:text-gray-400">
                                    Response Time:
                                </span>
                                <span className="ml-2 font-medium text-gray-800 dark:text-gray-200">
                                    {formatResponseTime(host.icmp_status.round_trip)}
                                </span>
                            </div>
                            <div>
                                <span className="text-gray-600 dark:text-gray-400">
                                    Packet Loss:
                                </span>
                                <span className="ml-2 font-medium text-gray-800 dark:text-gray-200">
                                    {host.icmp_status.packet_loss}%
                                </span>
                            </div>
                        </div>
                    </div>
                )}

                {host.port_results?.length > 0 && (
                    <div className="mt-4">
                        <h5 className="font-medium text-gray-800 dark:text-gray-200">
                            Open Ports
                        </h5>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 mt-2">
                            {host.port_results
                                .filter((port) => port.available)
                                .map((port) => (
                                    <div
                                        key={port.port}
                                        className="text-sm bg-gray-50 dark:bg-gray-700 p-2 rounded transition-colors"
                                    >
                                        <span className="font-medium text-gray-800 dark:text-gray-200">
                                            Port {port.port}
                                        </span>
                                        {port.service && (
                                            <span className="text-gray-600 dark:text-gray-400 ml-1">
                                                ({port.service})
                                            </span>
                                        )}
                                    </div>
                                ))}
                        </div>
                    </div>
                )}

                <div className="mt-4 text-xs text-gray-500 dark:text-gray-400">
                    <div>First seen: {new Date(host.first_seen).toLocaleString()}</div>
                    <div>Last seen: {new Date(host.last_seen).toLocaleString()}</div>
                </div>
            </div>
        </div>
    );
};

export default HostDetailsView;