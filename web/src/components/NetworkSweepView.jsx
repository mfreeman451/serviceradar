import { useState, useMemo, useEffect } from 'react';
import ExportButton from './ExportButton';
import { Filter, Search, ChevronDown, ChevronUp, Info, X, ChevronLeft, ChevronRight } from 'lucide-react';
import { HexGrid, Layout, Hexagon, Text, GridGenerator } from 'react-hexgrid';
import { consolidateNetworks, isIpInCidr } from '../lib/networkUtils';
import HostDetailsView from './HostDetailsView';
import './hexgrid.css';

const compareIPAddresses = (ip1, ip2) => {
    const ip1Parts = ip1.split('.').map(Number);
    const ip2Parts = ip2.split('.').map(Number);
    for (let i = 0; i < 4; i++) {
        if (ip1Parts[i] !== ip2Parts[i]) {
            return ip1Parts[i] - ip2Parts[i];
        }
    }
    return 0;
};


const NetworkSweepView = ({ nodeId, service, standalone = false }) => {
    const [viewMode, setViewMode] = useState('summary');
    const [searchTerm, setSearchTerm] = useState('');
    const [showFilters, setShowFilters] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [networkSearchTerm, setNetworkSearchTerm] = useState('');
    const [showNetworkInfo, setShowNetworkInfo] = useState(false);
    const [hexGridPage, setHexGridPage] = useState(1);
    const hostsPerPage = 10;
    const hexesPerPage = 5;

    const sweepDetails = typeof service.details === 'string'
        ? JSON.parse(service.details)
        : service.details;

    const networks = useMemo(() => {
        if (!sweepDetails?.network) return [];
        return sweepDetails.network.split(',');
    }, [sweepDetails]);

    const networkCount = networks.length;

    const filteredNetworks = useMemo(() => {
        if (networkSearchTerm === '') return networks.slice(0, 5);
        return networks.filter(network =>
            network.toLowerCase().includes(networkSearchTerm.toLowerCase())
        ).slice(0, 10);
    }, [networks, networkSearchTerm]);

    const sortAndFilterHosts = (hosts) => {
        if (!hosts) return [];
        return [...hosts]
            .filter((host) =>
                host.host.toLowerCase().includes(searchTerm.toLowerCase())
            )
            .sort((a, b) => compareIPAddresses(a.host, b.host));
    };

    const getRespondingHosts = (hosts) => {
        if (!hosts) return [];
        return hosts.filter((host) => {
            if (host.available) return true;
            const hasOpenPorts = host.port_results?.some((port) => port.available);
            const hasICMPResponse = host.icmp_status?.available === true;
            return hasOpenPorts || hasICMPResponse;
        });
    };

    const respondingHosts = getRespondingHosts(sweepDetails.hosts);

    const filteredHosts = sweepDetails.hosts
        ? sortAndFilterHosts(respondingHosts).filter((host) =>
            host.host.toLowerCase().includes(searchTerm.toLowerCase())
        )
        : [];

    const paginatedHosts = useMemo(() => {
        const startIndex = (currentPage - 1) * hostsPerPage;
        return filteredHosts.slice(startIndex, startIndex + hostsPerPage);
    }, [filteredHosts, currentPage]);

    const totalPages = Math.ceil(filteredHosts.length / hostsPerPage);

    const consolidatedNetworkData = useMemo(() => {
        const hostIPs = sweepDetails.hosts?.map(host => host.host) || [];
        return consolidateNetworks(hostIPs);
    }, [sweepDetails.hosts]);

    const getNetworkResponseTime = (network) => {
        const hostsInNetwork = sweepDetails.hosts?.filter(host =>
            host.host && isIpInCidr(host.host, network.network)
        ) || [];
        const respondingHosts = hostsInNetwork.filter(host =>
            host.icmp_status?.available
        );
        if (respondingHosts.length === 0) return 'N/A';
        const totalPingTime = respondingHosts.reduce((sum, host) => sum + (host.icmp_status?.round_trip || 0), 0);
        return `${(totalPingTime / respondingHosts.length / 1000000).toFixed(1)}ms`;
    };

    const paginatedNetworkData = useMemo(() => {
        const startIndex = (hexGridPage - 1) * hexesPerPage;
        return consolidatedNetworkData.slice(startIndex, startIndex + hexesPerPage);
    }, [consolidatedNetworkData, hexGridPage]);

    const totalHexPages = Math.ceil(consolidatedNetworkData.length / hexesPerPage);

    const [hexGridSize, setHexGridSize] = useState({ width: 600, height: 400 });


    useEffect(() => {
        const handleResize = () => {
            const width = window.innerWidth;
            if (width < 640) {
                setHexGridSize({ width: 350, height: 250 });
            } else if (width < 1024) {
                setHexGridSize({ width: 500, height: 350 });
            } else {
                setHexGridSize({ width: 700, height: 450 });
            }
        };

        handleResize();
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    const toggleFilters = () => setShowFilters(!showFilters);

    return (
        <div
            className={`space-y-4 ${
                !standalone &&
                'bg-white dark:bg-gray-800 rounded-lg shadow p-4 transition-colors'
            }`}
        >
            <div
                className={`
                    ${standalone
                    ? 'bg-white dark:bg-gray-800 rounded-lg shadow p-4 transition-colors'
                    : ''}
                `}
            >
                <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mb-4">
                    <div>
                        <h3 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                            Network Sweep
                        </h3>
                        <div className="flex items-center gap-2">
                            <span className="text-sm text-gray-600 dark:text-gray-400">
                                {networkCount} {networkCount === 1 ? 'network' : 'networks'} scanned
                            </span>
                            <button
                                onClick={() => setShowNetworkInfo(!showNetworkInfo)}
                                className="text-blue-500 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
                            >
                                <Info size={16} />
                            </button>
                        </div>
                    </div>

                    <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2">
                        <div className="flex flex-wrap gap-2">
                            <button
                                onClick={() => setViewMode('summary')}
                                className={`px-3 py-1 rounded transition-colors ${
                                    viewMode === 'summary'
                                        ? 'bg-blue-500 text-white'
                                        : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200'
                                }`}
                            >
                                Summary
                            </button>
                            <button
                                onClick={() => setViewMode('hosts')}
                                className={`px-3 py-1 rounded transition-colors ${
                                    viewMode === 'hosts'
                                        ? 'bg-blue-500 text-white'
                                        : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200'
                                }`}
                            >
                                Host Details
                            </button>
                            <button
                                onClick={() => setViewMode('hexgrid')}
                                className={`px-3 py-1 rounded transition-colors ${
                                    viewMode === 'hexgrid'
                                        ? 'bg-blue-500 text-white'
                                        : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200'
                                }`}
                            >
                                Hex Grid
                            </button>
                        </div>
                        <div className="hidden sm:block">
                            <ExportButton sweepDetails={sweepDetails} />
                        </div>
                    </div>
                </div>

                {showNetworkInfo && (
                    <div className="bg-blue-50 dark:bg-blue-900/30 p-4 rounded-lg mb-4 relative">
                        <button
                            onClick={() => setShowNetworkInfo(false)}
                            className="absolute top-2 right-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                        >
                            <X size={16} />
                        </button>
                        <h4 className="text-md font-medium text-gray-800 dark:text-gray-200 mb-2">
                            Networks Being Monitored
                        </h4>
                        <div className="relative mb-3">
                            <input
                                type="text"
                                placeholder="Search networks..."
                                className="w-full px-3 py-2 border rounded text-gray-700 dark:text-gray-200 dark:bg-gray-800 border-gray-300 dark:border-gray-600 pl-8"
                                value={networkSearchTerm}
                                onChange={(e) => setNetworkSearchTerm(e.target.value)}
                            />
                            <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-400" />
                        </div>
                        <div className="max-h-40 overflow-y-auto bg-white dark:bg-gray-800 rounded p-2 mb-2">
                            {filteredNetworks.length > 0 ? (
                                <ul className="list-disc list-inside">
                                    {filteredNetworks.map((network, index) => (
                                        <li key={index} className="text-sm text-gray-700 dark:text-gray-300 py-1">
                                            {network}
                                        </li>
                                    ))}
                                </ul>
                            ) : (
                                <p className="text-sm text-gray-500 dark:text-gray-400 italic">
                                    No matching networks found
                                </p>
                            )}
                        </div>
                        {networkSearchTerm === '' && networkCount > 5 && (
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                                Showing 5 of {networkCount} networks. Use the search to find specific networks.
                            </p>
                        )}
                        {networkSearchTerm !== '' && filteredNetworks.length < networkCount && (
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                                Showing {filteredNetworks.length} of {networkCount} networks matching "{networkSearchTerm}".
                            </p>
                        )}
                    </div>
                )}

                {viewMode === 'hosts' && (
                    <div className="flex items-center gap-2 mt-4">
                        <div className="relative flex-1">
                            <input
                                type="text"
                                placeholder="Search hosts..."
                                className="w-full p-2 border rounded text-gray-700 dark:text-gray-200 dark:bg-gray-800 border-gray-300 dark:border-gray-600 pl-8"
                                value={searchTerm}
                                onChange={(e) => {
                                    setSearchTerm(e.target.value);
                                    setCurrentPage(1);
                                }}
                            />
                            <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-400" />
                        </div>
                        <button
                            onClick={toggleFilters}
                            className="sm:hidden p-2 border rounded text-gray-700 dark:text-gray-200 dark:bg-gray-800 border-gray-300 dark:border-gray-600"
                        >
                            <Filter className="h-4 w-4" />
                        </button>
                    </div>
                )}

                <div className="sm:hidden mt-4">
                    <ExportButton sweepDetails={sweepDetails} />
                </div>
                <div className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                    Last sweep: {new Date(sweepDetails.last_sweep * 1000).toLocaleString()}
                </div>
            </div>

            {respondingHosts.length > 0 && viewMode !== 'hexgrid' && (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4 transition-colors">
                    <h4 className="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-200">
                        ICMP Status Summary
                    </h4>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                        <div className="bg-gray-50 dark:bg-gray-700 p-4 rounded transition-colors">
                            <div className="text-sm text-gray-600 dark:text-gray-400">
                                ICMP Responding
                            </div>
                            <div className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100">
                                {
                                    respondingHosts.filter(
                                        (h) =>
                                            h.icmp_status?.available && h.icmp_status?.packet_loss === 0
                                    ).length
                                }
                                <span className="text-sm text-gray-500 dark:text-gray-300 ml-2">
                                    hosts
                                </span>
                            </div>
                        </div>
                        <div className="bg-gray-50 dark:bg-gray-700 p-4 rounded transition-colors">
                            <div className="text-sm text-gray-600 dark:text-gray-400">
                                Average Response Time
                            </div>
                            <div className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100">
                                {(() => {
                                    const respondingToICMP = respondingHosts.filter(
                                        (h) =>
                                            h.icmp_status?.available &&
                                            h.icmp_status?.packet_loss === 0 &&
                                            h.icmp_status?.round_trip > 0
                                    );
                                    if (respondingToICMP.length === 0) return 'N/A';
                                    const avg =
                                        respondingToICMP.reduce(
                                            (acc, h) => acc + h.icmp_status.round_trip,
                                            0
                                        ) /
                                        respondingToICMP.length /
                                        1000000;
                                    return `${avg.toFixed(2)}ms`;
                                })()}
                            </div>
                        </div>
                        <div className="bg-gray-50 dark:bg-gray-700 p-4 rounded transition-colors">
                            <div className="text-sm text-gray-600 dark:text-gray-400">
                                Open Services
                            </div>
                            <div className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100">
                                {respondingHosts.reduce(
                                    (acc, host) =>
                                        acc +
                                        (host.port_results?.filter((port) => port.available)
                                            ?.length || 0),
                                    0
                                )}
                                <span className="text-sm text-gray-500 dark:text-gray-300 ml-2">
                                    ports
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {viewMode === 'summary' ? (
                <div
                    className={`
                        ${standalone
                        ? 'bg-white dark:bg-gray-800 rounded-lg shadow p-4 transition-colors'
                        : ''}
                    `}
                >
                    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                        {sweepDetails.ports
                            ?.sort((a, b) => b.available - a.available)
                            .map((port) => (
                                <div
                                    key={port.port}
                                    className="bg-gray-50 dark:bg-gray-700 p-3 rounded transition-colors"
                                >
                                    <div className="font-medium text-gray-800 dark:text-gray-200">
                                        Port {port.port}
                                    </div>
                                    <div className="text-sm text-gray-600 dark:text-gray-400">
                                        {port.available} hosts responding
                                    </div>
                                    <div className="mt-1 bg-gray-200 dark:bg-gray-600 rounded-full h-2">
                                        <div
                                            className="bg-blue-500 rounded-full h-2"
                                            style={{
                                                width: `${
                                                    (port.available / sweepDetails.total_hosts) * 100
                                                }%`,
                                            }}
                                        />
                                    </div>
                                </div>
                            ))}
                    </div>
                </div>
            ) : viewMode === 'hosts' ? (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {paginatedHosts.map((host) => (
                            <HostDetailsView key={host.host} host={host} />
                        ))}
                    </div>
                    {totalPages > 1 && (
                        <div className="flex justify-center mt-4">
                            <div className="flex items-center space-x-2">
                                <button
                                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                                    disabled={currentPage === 1}
                                    className={`px-3 py-1 rounded ${
                                        currentPage === 1
                                            ? 'bg-gray-200 dark:bg-gray-700 text-gray-400 dark:text-gray-500 cursor-not-allowed'
                                            : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-600'
                                    }`}
                                >
                                    Previous
                                </button>
                                <div className="text-sm text-gray-600 dark:text-gray-400">
                                    Page {currentPage} of {totalPages}
                                </div>
                                <button
                                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                                    disabled={currentPage === totalPages}
                                    className={`px-3 py-1 rounded ${
                                        currentPage === totalPages
                                            ? 'bg-gray-200 dark:bg-gray-700 text-gray-400 dark:text-gray-500 cursor-not-allowed'
                                            : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-600'
                                    }`}
                                >
                                    Next
                                </button>
                            </div>
                        </div>
                    )}
                    <div className="text-center text-sm text-gray-500 dark:text-gray-400">
                        Showing {filteredHosts.length > 0 ? (currentPage - 1) * hostsPerPage + 1 : 0}-
                        {Math.min(currentPage * hostsPerPage, filteredHosts.length)} of {filteredHosts.length} hosts
                    </div>
                </>
            ) : viewMode === 'hexgrid' ? (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4 transition-colors">
                    <h4 className="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-200">
                        Network Hex Grid
                    </h4>
                    {paginatedNetworkData.length === 0 ? (
                        <div className="text-gray-500 dark:text-gray-400 text-center py-4">
                            No networks available to display
                        </div>
                    ) : (
                        <>
                            <HexGrid width={hexGridSize.width} height={hexGridSize.height}>
                                <Layout size={{ x: 15, y: 15 }} flat={false} origin={{ x: 0, y: 0 }}>
                                    {GridGenerator.hexagon(2).slice(0, paginatedNetworkData.length).map((hex, index) => {
                                        const networkData = paginatedNetworkData[index];
                                        if (!networkData) return null;
                                        const hasHosts = networkData.hosts > 0;
                                        console.log(`Network: ${networkData.network}, hasHosts: ${hasHosts}, fill: ${hasHosts ? '#3F51B5' : '#757575'}`);
                                        return (
                                            <Hexagon
                                                key={index}
                                                q={hex.q}
                                                r={hex.r}
                                                s={hex.s}
                                                fill={hasHosts ? '#3F51B5' : '#757575'}
                                                className={hasHosts ? 'cursor-pointer' : ''}
                                            >
                                                <Text className="font-bold">
                                                    {networkData.network.split('/')[0]}
                                                </Text>
                                                <Text y={1.5}>
                                                    {getNetworkResponseTime(networkData) || 'N/A'}
                                                </Text>
                                                <Text y={2.5} className="fill-opacity-70">
                                                    {`${networkData.hosts} hosts`}
                                                </Text>
                                            </Hexagon>
                                        );
                                    })}
                                </Layout>
                            </HexGrid>
                            {totalHexPages > 1 && (
                                <div className="flex justify-between items-center mt-4">
                                    <div className="text-sm text-gray-500 dark:text-gray-400">
                                        Showing {(hexGridPage - 1) * hexesPerPage + 1}-
                                        {Math.min(hexGridPage * hexesPerPage, consolidatedNetworkData.length)} of {consolidatedNetworkData.length} networks
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <button
                                            onClick={() => setHexGridPage(prev => Math.max(prev - 1, 1))}
                                            disabled={hexGridPage === 1}
                                            className={`p-2 rounded-full ${
                                                hexGridPage === 1
                                                    ? 'text-gray-400 dark:text-gray-500 cursor-not-allowed'
                                                    : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                                            }`}
                                        >
                                            <ChevronLeft size={20} />
                                        </button>
                                        <span className="text-sm text-gray-600 dark:text-gray-300">
                                            {hexGridPage} / {totalHexPages}
                                        </span>
                                        <button
                                            onClick={() => setHexGridPage(prev => Math.min(prev + 1, totalHexPages))}
                                            disabled={hexGridPage === totalHexPages}
                                            className={`p-2 rounded-full ${
                                                hexGridPage === totalHexPages
                                                    ? 'text-gray-400 dark:text-gray-500 cursor-not-allowed'
                                                    : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                                            }`}
                                        >
                                            <ChevronRight size={20} />
                                        </button>
                                    </div>
                                </div>
                            )}
                        </>
                    )}
                </div>
            ) : null}
        </div>
    );
};

export default NetworkSweepView;