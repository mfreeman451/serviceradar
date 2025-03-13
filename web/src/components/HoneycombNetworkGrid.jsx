import React, { useState, useEffect, useRef, useMemo } from 'react';
import { ChevronLeft, ChevronRight, Search, X } from 'lucide-react';

const HoneycombNetworkGrid = ({ networks = [], onClose }) => {
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(0);
    const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
    const [hoveredNetwork, setHoveredNetwork] = useState(null);
    const containerRef = useRef(null);
    const canvasRef = useRef(null);

    // Detect mouse position for hover effects
    const [mousePos, setMousePos] = useState({ x: 0, y: 0 });

    const calculateGridDimensions = useMemo(() => {
        if (!containerRef.current) return { cols: 8, rows: 4 };

        const containerWidth = containerRef.current.clientWidth;
        const containerHeight = containerRef.current.clientHeight;
        const hexSize = 50; // Base hex radius size

        const cols = Math.floor(containerWidth / (hexSize * 1.8));
        const rows = Math.floor(containerHeight / (hexSize * 1.5));

        return {
            cols: Math.max(Math.min(cols, 12), 4),
            rows: Math.max(Math.min(rows, 6), 3)
        };
    }, [dimensions]);

    const filteredNetworks = useMemo(() => {
        if (!searchTerm) return networks;
        return networks.filter(network =>
            network && typeof network === 'string' && network.toLowerCase().includes(searchTerm.toLowerCase())
        );
    }, [networks, searchTerm]);

    const hexesPerPage = calculateGridDimensions.cols * calculateGridDimensions.rows;
    const totalPages = Math.ceil(filteredNetworks.length / hexesPerPage);
    const currentNetworks = useMemo(() => {
        const startIndex = currentPage * hexesPerPage;
        return filteredNetworks.slice(startIndex, startIndex + hexesPerPage);
    }, [currentPage, filteredNetworks, hexesPerPage]);

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

    // Track mouse movement for hover effects
    const handleMouseMove = (e) => {
        if (!canvasRef.current) return;

        const rect = canvasRef.current.getBoundingClientRect();
        setMousePos({
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        });
    };

    // Check if a point is inside a hexagon
    const isPointInHexagon = (x, y, hexX, hexY, hexSize) => {
        const dx = Math.abs(x - hexX);
        const dy = Math.abs(y - hexY);

        // Quick check using bounding rectangle
        if (dx > hexSize * Math.sqrt(3) / 2 || dy > hexSize) return false;

        // More precise check using distance
        return hexSize * hexSize >= dx * dx + dy * dy * 3;
    };

    // Get ping data for a network
    const getNetworkStatus = (network) => {
        // Use a consistent hash from the network string for demo
        let hash = 0;
        for (let i = 0; i < network.length; i++) {
            hash = ((hash << 5) - hash) + network.charCodeAt(i);
            hash |= 0; // Convert to 32bit integer
        }

        // Some networks don't respond to ping (about 10% in this demo)
        const responds = (Math.abs(hash) % 10) !== 0;

        if (!responds) {
            return {
                responds: false,
                pingTime: null
            };
        }

        // Generate a ping time between 25-125ms based on hash
        return {
            responds: true,
            pingTime: Math.abs(hash % 100) + 25
        };
    };

    // Get IP-based background color
    const getNetworkColor = (network) => {
        if (network.startsWith('10.')) return '#7c3aed'; // Purple
        if (network.startsWith('172.')) return '#2563eb'; // Blue
        if (network.startsWith('192.168.')) return '#0891b2'; // Teal
        return '#4f46e5'; // Indigo
    };

    useEffect(() => {
        if (!canvasRef.current || !currentNetworks.length) return;

        const canvas = canvasRef.current;
        const ctx = canvas.getContext('2d');
        if (!ctx) {
            console.error('Failed to get 2D context');
            return;
        }

        // Set canvas size to match container dimensions
        canvas.width = dimensions.width || containerRef.current.clientWidth;
        canvas.height = dimensions.height || containerRef.current.clientHeight;

        ctx.clearRect(0, 0, canvas.width, canvas.height);

        // Background pattern
        drawHexagonalBackground(ctx, canvas.width, canvas.height);

        const hexSize = 45; // Hexagon radius
        const width = calculateGridDimensions.cols;
        const height = calculateGridDimensions.rows;

        // Hexagon geometry
        const hexWidth = hexSize * Math.sqrt(3);
        const hexHeight = hexSize * 2;
        const horizontalSpacing = hexWidth;
        const verticalSpacing = hexHeight * 0.75;

        // Offsets to center the grid
        const offsetX = (canvas.width - (width - 1) * horizontalSpacing) / 2;
        const offsetY = (canvas.height - (height - 1) * verticalSpacing) / 2;

        // Store hexagon data for hover detection
        const hexagons = [];

        // Draw hexagons
        let index = 0;
        for (let row = 0; row < height; row++) {
            for (let col = 0; col < width; col++) {
                if (index >= currentNetworks.length) continue;

                const network = currentNetworks[index++];
                if (!network || typeof network !== 'string') continue;

                // Calculate center position of the hexagon with staggered rows
                const x = offsetX + col * horizontalSpacing + (row % 2 === 1 ? horizontalSpacing / 2 : 0);
                const y = offsetY + row * verticalSpacing;

                // Store hexagon data for hover detection
                hexagons.push({
                    x, y, network, size: hexSize
                });

                // Get network status (ping data)
                const networkStatus = getNetworkStatus(network);

                // Check if this hexagon is being hovered
                const isHovered = hoveredNetwork === network;

                // Draw the hexagon
                drawHexagon(ctx, x, y, hexSize * (isHovered ? 1.1 : 1), network, networkStatus, isHovered);
            }
        }

        // Check if mouse is hovering over any hexagon
        const hoveredHex = hexagons.find(h =>
            isPointInHexagon(mousePos.x, mousePos.y, h.x, h.y, h.size)
        );

        if (hoveredHex?.network !== hoveredNetwork) {
            setHoveredNetwork(hoveredHex?.network || null);
        }

    }, [currentNetworks, dimensions, calculateGridDimensions, mousePos, hoveredNetwork]);

    // Draw background pattern of subtle hexagons
    const drawHexagonalBackground = (ctx, width, height) => {
        const hexSize = 30;
        const hexWidth = hexSize * Math.sqrt(3);
        const hexHeight = hexSize * 2;
        const horizontalSpacing = hexWidth;
        const verticalSpacing = hexHeight * 0.75;

        const cols = Math.ceil(width / horizontalSpacing) + 2;
        const rows = Math.ceil(height / verticalSpacing) + 2;

        ctx.strokeStyle = 'rgba(255, 255, 255, 0.05)';
        ctx.lineWidth = 1;

        for (let row = -1; row < rows; row++) {
            for (let col = -1; col < cols; col++) {
                const x = col * horizontalSpacing + (row % 2 === 1 ? horizontalSpacing / 2 : 0);
                const y = row * verticalSpacing;

                ctx.beginPath();
                for (let i = 0; i < 6; i++) {
                    const angle = 2 * Math.PI / 6 * i - Math.PI / 6;
                    const px = x + hexSize * Math.cos(angle);
                    const py = y + hexSize * Math.sin(angle);
                    if (i === 0) {
                        ctx.moveTo(px, py);
                    } else {
                        ctx.lineTo(px, py);
                    }
                }
                ctx.closePath();
                ctx.stroke();
            }
        }
    };

    // Draw a single hexagon with network info
    const drawHexagon = (ctx, x, y, size, network, networkStatus, isHovered) => {
        // Get color based on network type
        let fillColor = getNetworkColor(network);

        // If the network doesn't respond to ping, use a muted color
        if (!networkStatus.responds) {
            fillColor = adjustColor(fillColor, -40); // Make it darker/muted
        }

        // Draw hexagon path
        ctx.beginPath();
        for (let i = 0; i < 6; i++) {
            const angle = 2 * Math.PI / 6 * i - Math.PI / 6; // Start at the right corner
            const px = x + size * Math.cos(angle);
            const py = y + size * Math.sin(angle);
            if (i === 0) {
                ctx.moveTo(px, py);
            } else {
                ctx.lineTo(px, py);
            }
        }
        ctx.closePath();

        // Fill with gradient for more depth
        const gradient = ctx.createRadialGradient(x, y, 0, x, y, size);
        gradient.addColorStop(0, fillColor);
        gradient.addColorStop(1, adjustColor(fillColor, -20));
        ctx.fillStyle = gradient;
        ctx.fill();

        // Stroke with glow effect if hovered
        if (isHovered) {
            ctx.strokeStyle = 'white';
            ctx.lineWidth = 3;
            ctx.shadowColor = 'rgba(255, 255, 255, 0.8)';
            ctx.shadowBlur = 10;
        } else {
            ctx.strokeStyle = 'rgba(30, 41, 59, 0.8)';
            ctx.lineWidth = 1.5;
            ctx.shadowBlur = 0;
        }
        ctx.stroke();
        ctx.shadowBlur = 0;

        // Add text
        ctx.fillStyle = 'white';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';

        // Display shortened network ID
        const displayText = getDisplayText(network);

        // On hover, show the full network address
        if (isHovered) {
            ctx.font = 'bold 14px sans-serif';
            ctx.fillText(network, x, y - 8);

            if (networkStatus.responds) {
                ctx.font = 'bold 16px sans-serif';
                ctx.fillText(`${networkStatus.pingTime.toFixed(1)}ms`, x, y + 12);
            } else {
                ctx.font = 'bold 14px sans-serif';
                ctx.fillStyle = 'rgba(255, 255, 255, 0.7)';
                ctx.fillText('No Response', x, y + 12);
            }
        } else {
            // Regular view (not hovered)
            ctx.font = 'bold 12px sans-serif';
            ctx.fillText(displayText, x, y - 5);

            if (networkStatus.responds) {
                ctx.font = '12px sans-serif';
                ctx.fillText(`${networkStatus.pingTime.toFixed(1)}ms`, x, y + 12);
            } else {
                ctx.font = '11px sans-serif';
                ctx.fillStyle = 'rgba(255, 255, 255, 0.7)';
                ctx.fillText('No Response', x, y + 12);
            }
        }
    };

    // Helper to get display text for network
    const getDisplayText = (network) => {
        // For IP addresses with CIDR notation (e.g., 192.168.1.0/24)
        if (network.includes('/')) {
            const parts = network.split('/');
            const ipParts = parts[0].split('.');

            // For subnet, show first three octets + CIDR
            // Example: 192.168.1.0/24 → 192.168.1/24
            if (ipParts[3] === '0') {
                return `${ipParts[0]}.${ipParts[1]}.${ipParts[2]}/${parts[1]}`;
            }

            // For host addresses, show meaningful part
            // If it's a standard private network, show last two octets
            if (ipParts[0] === '10' ||
                (ipParts[0] === '172' && parseInt(ipParts[1]) >= 16 && parseInt(ipParts[1]) <= 31) ||
                (ipParts[0] === '192' && ipParts[1] === '168')) {
                return `${ipParts[2]}.${ipParts[3]}`;
            }

            // Otherwise show full IP without CIDR
            return parts[0];
        }

        // For IP addresses without CIDR
        const ipParts = network.split('.');
        if (ipParts.length === 4) {
            // If it's a standard private network, show last two octets
            if (ipParts[0] === '10' ||
                (ipParts[0] === '172' && parseInt(ipParts[1]) >= 16 && parseInt(ipParts[1]) <= 31) ||
                (ipParts[0] === '192' && ipParts[1] === '168')) {
                return `${ipParts[2]}.${ipParts[3]}`;
            }

            // Otherwise show full IP
            return network;
        }

        // For non-IP addresses, truncate if too long
        return network.length > 8 ? network.substring(0, 6) + '..' : network;
    };

    // Helper to adjust color brightness
    const adjustColor = (color, amount) => {
        const clamp = (val) => Math.min(255, Math.max(0, val));

        // Convert hex to RGB
        let r = parseInt(color.substring(1, 3), 16);
        let g = parseInt(color.substring(3, 5), 16);
        let b = parseInt(color.substring(5, 7), 16);

        // Adjust brightness
        r = clamp(r + amount);
        g = clamp(g + amount);
        b = clamp(b + amount);

        // Convert back to hex
        return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
    };

    return (
        <div className="bg-gray-800 rounded-lg text-white p-4 flex flex-col h-full" style={{ minHeight: '400px' }}>
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-medium">Networks ({filteredNetworks.length})</h3>
                <button
                    onClick={onClose}
                    className="text-gray-400 hover:text-white"
                    aria-label="Close network view"
                >
                    <X size={18} />
                </button>
            </div>

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

            <div
                ref={containerRef}
                className="flex-grow relative overflow-hidden"
                style={{ position: 'relative' }}
                onMouseMove={handleMouseMove}
                onMouseLeave={() => setHoveredNetwork(null)}
            >
                <canvas
                    ref={canvasRef}
                    className="w-full h-full"
                />
            </div>

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