import React, { useState, useEffect, useRef, useMemo } from 'react';
import { ChevronLeft, ChevronRight, Search, X } from 'lucide-react';
import { isIpInCidr } from '../lib/networkUtils';

const HoneycombNetworkGrid = ({ networks = [], onClose, sweepDetails }) => {
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(0);
    const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
    const [hoveredNetwork, setHoveredNetwork] = useState(null);
    const [consolidatedNetworks, setConsolidatedNetworks] = useState([]);
    const containerRef = useRef(null);
    const canvasRef = useRef(null);
    const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
    const [hexagons, setHexagons] = useState([]);

    // Set initial networks
    useEffect(() => {
        setConsolidatedNetworks(networks);
    }, [networks]);

    // Get network status from sweep data
    const getNetworkStatus = (network) => {
        if (!sweepDetails?.hosts || !network) {
            return { responds: false, pingTime: null };
        }

        const hostsInNetwork = sweepDetails.hosts.filter(host =>
            host.host && isIpInCidr(host.host, network)
        );

        const respondingHosts = hostsInNetwork.filter(host =>
            host.icmp_status && host.icmp_status.available
        );

        if (respondingHosts.length === 0) {
            return {
                responds: false,
                pingTime: null,
                respondingCount: 0,
                totalCount: hostsInNetwork.length
            };
        }

        const totalPingTime = respondingHosts.reduce((sum, host) => {
            return sum + (host.icmp_status.round_trip || 0);
        }, 0);

        const avgPingTime = totalPingTime / respondingHosts.length / 1000000;

        return {
            responds: true,
            pingTime: avgPingTime,
            respondingCount: respondingHosts.length,
            totalCount: hostsInNetwork.length
        };
    };

    // Get color based on network type
    const getNetworkColor = (network) => {
        if (!network) return '#4f46e5';
        return '#3b82f6'; // Consistent blue for all networks
    };

    // Adjust color brightness
    const adjustColor = (color, amount) => {
        const clamp = (val) => Math.min(255, Math.max(0, val));

        let r = parseInt(color.substring(1, 3), 16);
        let g = parseInt(color.substring(3, 5), 16);
        let b = parseInt(color.substring(5, 7), 16);

        r = clamp(r + amount);
        g = clamp(g + amount);
        b = clamp(b + amount);

        return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
    };

    // Calculate grid dimensions
    const calculateGridDimensions = useMemo(() => {
        if (!containerRef.current) return { hexSize: 50 };

        const width = containerRef.current.clientWidth;
        const height = containerRef.current.clientHeight;

        const hexSize = Math.min(width, height) / 5; // Keep the large size
        const hexWidth = hexSize * Math.sqrt(3);
        const hexHeight = hexSize * 2;
        const horizontalSpacing = hexWidth;
        const verticalSpacing = hexHeight * 0.75;

        return {
            hexSize,
            horizontalSpacing,
            verticalSpacing,
        };
    }, [dimensions]);

    // Filter and sort networks based on search and response status
    const filteredNetworks = useMemo(() => {
        let networks = [...consolidatedNetworks];

        if (searchTerm) {
            networks = networks.filter(network =>
                network.toLowerCase().includes(searchTerm.toLowerCase())
            );
        }

        return networks.sort((a, b) => {
            const statusA = getNetworkStatus(a);
            const statusB = getNetworkStatus(b);

            if (statusA.responds && !statusB.responds) return -1;
            if (!statusA.responds && statusB.responds) return 1;

            if (statusA.responds && statusB.responds) {
                return statusA.pingTime - statusB.pingTime;
            }

            return statusB.totalCount - statusA.totalCount;
        });
    }, [consolidatedNetworks, searchTerm]);

    // Pagination
    const hexesPerPage = 5; // Cluster of 5 hexagons per page
    const totalPages = Math.ceil(filteredNetworks.length / hexesPerPage);
    const currentNetworks = useMemo(() => {
        const startIndex = currentPage * hexesPerPage;
        return filteredNetworks.slice(startIndex, startIndex + hexesPerPage);
    }, [currentPage, filteredNetworks, hexesPerPage]);

    // Navigation
    const nextPage = () => {
        if (currentPage < totalPages - 1) {
            setCurrentPage(prev => prev + 1);
        }
    };

    const prevPage = () => {
        if (currentPage > 0) {
            setCurrentPage(prev => prev - 1);
        }
    };

    // Handle resize
    useEffect(() => {
        const handleResize = () => {
            if (containerRef.current) {
                setDimensions({
                    width: containerRef.current.clientWidth,
                    height: containerRef.current.clientHeight
                });
            }
        };

        handleResize();
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    // Track mouse movement
    const handleMouseMove = (e) => {
        if (!canvasRef.current) return;

        const rect = canvasRef.current.getBoundingClientRect();
        setMousePos({
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        });

        const hovered = hexagons.find(hex =>
            isPointInHexagon(mousePos.x, mousePos.y, hex.x, hex.y, hex.size)
        );

        if (hovered?.network !== hoveredNetwork) {
            setHoveredNetwork(hovered?.network || null);
        }
    };

    // Check if point is inside a hexagon
    const isPointInHexagon = (x, y, hexX, hexY, hexSize) => {
        const dx = Math.abs(x - hexX);
        const dy = Math.abs(y - hexY);

        const rx = hexSize * Math.sqrt(3) / 2;
        const ry = hexSize;

        if (dx > rx || dy > ry) return false;

        return 3 * hexSize * hexSize / 4 >= (dx * dx + dy * dy);
    };

    // Draw the canvas with network hexagons
    useEffect(() => {
        if (!canvasRef.current || !containerRef.current) return;

        const canvas = canvasRef.current;
        const ctx = canvas.getContext('2d');
        if (!ctx) return;

        canvas.width = containerRef.current.clientWidth;
        canvas.height = containerRef.current.clientHeight;

        const width = canvas.width;
        const height = canvas.height;

        ctx.clearRect(0, 0, width, height);

        // Removed drawHexagonalBackground to eliminate the background grid

        const { hexSize, horizontalSpacing, verticalSpacing } = calculateGridDimensions;

        // Define positions for a cluster of up to 5 hexagons (like Synadia)
        const positions = [
            { x: 0, y: 0 }, // Center hexagon
            { x: horizontalSpacing, y: 0 }, // Right
            { x: -horizontalSpacing, y: 0 }, // Left
            { x: horizontalSpacing / 2, y: verticalSpacing }, // Bottom-right
            { x: -horizontalSpacing / 2, y: verticalSpacing }, // Bottom-left
        ];

        // Calculate the bounding box of the cluster
        const clusterWidth = 2 * horizontalSpacing;
        const clusterHeight = verticalSpacing + hexSize;

        // Center the cluster in the canvas
        const offsetX = (width - clusterWidth) / 2;
        const offsetY = (height - clusterHeight) / 2;

        const hexagonData = [];

        let index = 0;
        for (const pos of positions) {
            if (index >= currentNetworks.length) break;

            const network = currentNetworks[index++];
            if (!network) continue;

            const x = offsetX + pos.x + horizontalSpacing; // Adjust for centering
            const y = offsetY + pos.y + hexSize; // Adjust for centering

            const isHovered = hoveredNetwork === network;
            const networkStatus = getNetworkStatus(network);

            drawNetworkHexagon(
                ctx,
                x,
                y,
                hexSize * (isHovered ? 1.1 : 1),
                network,
                networkStatus,
                isHovered
            );

            hexagonData.push({
                x,
                y,
                network,
                size: hexSize,
                status: networkStatus,
            });
        }

        setHexagons(hexagonData);
    }, [currentNetworks, dimensions, calculateGridDimensions, hoveredNetwork, mousePos]);

    // Draw an individual network hexagon
    function drawNetworkHexagon(ctx, x, y, size, network, networkStatus, isHovered) {
        let baseColor = getNetworkColor(network);

        if (!networkStatus.responds) {
            baseColor = adjustColor(baseColor, -50);
        }

        ctx.beginPath();
        for (let i = 0; i < 6; i++) {
            const angle = 2 * Math.PI / 6 * i;
            const px = x + size * Math.cos(angle);
            const py = y + size * Math.sin(angle);

            if (i === 0) ctx.moveTo(px, py);
            else ctx.lineTo(px, py);
        }
        ctx.closePath();

        // Enhanced gradient for more depth
        const gradient = ctx.createRadialGradient(x, y, 0, x, y, size);
        gradient.addColorStop(0, adjustColor(baseColor, 20));
        gradient.addColorStop(1, adjustColor(baseColor, -40));
        ctx.fillStyle = gradient;
        ctx.fill();

        ctx.save();
        ctx.clip();
        ctx.shadowBlur = 15;
        ctx.shadowColor = 'rgba(0, 0, 0, 0.5)';
        ctx.shadowOffsetX = 0;
        ctx.shadowOffsetY = 0;
        ctx.fill();
        ctx.restore();

        if (isHovered) {
            ctx.strokeStyle = 'rgba(255, 255, 255, 1)';
            ctx.lineWidth = 4; // Thicker border on hover
            ctx.shadowColor = 'rgba(255, 255, 255, 0.8)';
            ctx.shadowBlur = 12;
            ctx.stroke();
            ctx.shadowBlur = 0;

            ctx.beginPath();
            for (let i = 0; i < 6; i++) {
                const angle = 2 * Math.PI / 6 * i;
                const px = x + (size - 4) * Math.cos(angle);
                const py = y + (size - 4) * Math.sin(angle);

                if (i === 0) ctx.moveTo(px, py);
                else ctx.lineTo(px, py);
            }
            ctx.closePath();
            ctx.strokeStyle = 'rgba(255, 255, 255, 0.6)';
            ctx.lineWidth = 2;
            ctx.stroke();
        } else {
            ctx.strokeStyle = 'rgba(255, 255, 255, 0.3)';
            ctx.lineWidth = 2;
            ctx.stroke();
        }

        const textX = x;
        const textY = y;

        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillStyle = 'white';

        if (isHovered) {
            ctx.font = 'bold 18px sans-serif';
            ctx.fillText(network, textX, textY - 25);

            if (networkStatus.responds) {
                ctx.font = 'bold 16px sans-serif';
                ctx.fillText(`${networkStatus.pingTime.toFixed(1)}ms`, textX, textY + 5);

                ctx.font = '16px sans-serif';
                ctx.fillText(
                    `${networkStatus.respondingCount}/${networkStatus.totalCount}`,
                    textX,
                    textY + 30
                );
            } else {
                ctx.font = 'bold 16px sans-serif';
                ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
                ctx.fillText('No Response', textX, textY + 5);

                if (networkStatus.totalCount > 0) {
                    ctx.font = '16px sans-serif';
                    ctx.fillText(`${networkStatus.totalCount} hosts`, textX, textY + 30);
                }
            }
        } else {
            ctx.font = 'bold 16px sans-serif';
            ctx.fillText(network, textX, textY - 20);

            if (networkStatus.responds) {
                ctx.font = '16px sans-serif';
                ctx.fillText(`${networkStatus.pingTime.toFixed(1)}ms`, textX, textY + 5);
            } else {
                ctx.font = '16px sans-serif';
                ctx.fillStyle = 'rgba(255, 255, 255, 0.7)';
                ctx.fillText('No Response', textX, textY + 5);
            }
        }
    }

    return (
        <div className="bg-gray-800 rounded-lg text-white p-4 flex flex-col h-full" style={{ minHeight: '400px' }}>
            {/* Header */}
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-medium">Networks ({consolidatedNetworks.length})</h3>
                <button
                    onClick={onClose}
                    className="text-gray-400 hover:text-white transition-colors"
                    aria-label="Close network view"
                >
                    <X size={18} />
                </button>
            </div>

            {/* Search */}
            <div className="relative mb-4">
                <input
                    type="text"
                    placeholder="Search networks..."
                    className="w-full bg-gray-700 text-white border border-gray-600 rounded px-8 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={searchTerm}
                    onChange={(e) => {
                        setSearchTerm(e.target.value);
                        setCurrentPage(0);
                    }}
                />
                <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                {searchTerm && (
                    <button
                        onClick={() => {
                            setSearchTerm('');
                            setCurrentPage(0);
                        }}
                        className="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-white"
                    >
                        <X size={16} />
                    </button>
                )}
            </div>

            {/* Canvas container */}
            <div
                ref={containerRef}
                className="flex-grow relative"
                style={{ height: '300px' }}
                onMouseMove={handleMouseMove}
                onMouseLeave={() => setHoveredNetwork(null)}
            >
                <canvas
                    ref={canvasRef}
                    className="w-full h-full"
                />
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex justify-between items-center mt-4">
                    <div className="text-sm text-gray-400">
                        Showing {currentPage * hexesPerPage + 1}-
                        {Math.min((currentPage + 1) * hexesPerPage, filteredNetworks.length)} of {filteredNetworks.length}
                    </div>
                    <div className="flex items-center gap-2">
                        <button
                            onClick={prevPage}
                            disabled={currentPage === 0}
                            className={`p-2 rounded-full ${
                                currentPage === 0
                                    ? 'text-gray-600 cursor-not-allowed'
                                    : 'text-gray-300 hover:bg-gray-700'
                            }`}
                        >
                            <ChevronLeft size={20} />
                        </button>
                        <span className="text-sm text-gray-300">
                            {currentPage + 1} / {totalPages}
                        </span>
                        <button
                            onClick={nextPage}
                            disabled={currentPage >= totalPages - 1}
                            className={`p-2 rounded-full ${
                                currentPage >= totalPages - 1
                                    ? 'text-gray-600 cursor-not-allowed'
                                    : 'text-gray-300 hover:bg-gray-700'
                            }`}
                        >
                            <ChevronRight size={20} />
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default HoneycombNetworkGrid;