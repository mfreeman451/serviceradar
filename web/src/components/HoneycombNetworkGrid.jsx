/*
 * Copyright 2025 Carver Automation Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useState, useRef, useEffect } from 'react';
import { HexGrid, Layout, Hexagon } from 'react-hexgrid';
import { Search } from 'lucide-react';

// Super simple implementation
const HoneycombNetworkGrid = ({ networks = [], searchTerm = '', onSearchChange, statusData = {} }) => {
    const [hoveredNetwork, setHoveredNetwork] = useState(null);
    const containerRef = useRef(null);
    const [containerSize, setContainerSize] = useState({ width: 800, height: 400 });

    // Track container size
    useEffect(() => {
        if (!containerRef.current) return;

        const updateSize = () => {
            if (containerRef.current) {
                const width = containerRef.current.offsetWidth;
                setContainerSize({
                    width,
                    height: Math.min(500, width * 0.5)
                });
            }
        };

        updateSize();
        window.addEventListener('resize', updateSize);
        return () => window.removeEventListener('resize', updateSize);
    }, []);

    // Filter networks
    const filteredNetworks = networks.filter(network =>
        !searchTerm || network.toLowerCase().includes(searchTerm.toLowerCase())
    );

    // Get status colors
    const getStatusColor = (status) => {
        switch(status) {
            case 'online': return '#10b981'; // Green
            case 'offline': return '#ef4444'; // Red
            case 'warning': return '#f59e0b'; // Amber
            default: return '#6b7280'; // Gray
        }
    };

    // Extremely simple hex grid - just place them in a spiral-like pattern
    const createHexGrid = () => {
        if (filteredNetworks.length === 0) return [];

        // Calculate center of the container
        const centerX = containerSize.width / 2;
        const centerY = containerSize.height / 2;

        // Hexagon size - smaller for more networks
        const size = Math.min(30, Math.max(20, 50 - filteredNetworks.length / 5));

        // Spiral layout variables
        let angle = 0;
        const step = Math.min(2, Math.max(0.5, filteredNetworks.length / 20));
        const hexSize = size * 1.8;

        return filteredNetworks.map((network, i) => {
            // Create a spiral pattern
            const radius = Math.sqrt(i) * hexSize / 1.5;
            const x = centerX + radius * Math.cos(angle);
            const y = centerY + radius * Math.sin(angle);

            // Status color
            const status = statusData[network]?.status || 'unknown';
            const color = getStatusColor(status);

            // Extract last octet or CIDR
            const displayText = network.split('/')[0].split('.')[3] || '';

            // Increment angle for spiral
            angle += step;

            return {
                network,
                status,
                x,
                y,
                color,
                displayText,
                size
            };
        });
    };

    // Create hexagons
    const hexGrid = createHexGrid();

    return (
        <div className="w-full" ref={containerRef}>
            {/* Search bar */}
            <div className="mb-4 relative">
                <input
                    type="text"
                    placeholder="Search networks..."
                    className="w-full p-2 pl-8 border rounded text-gray-700 dark:text-gray-200 dark:bg-gray-700 border-gray-300 dark:border-gray-600"
                    value={searchTerm}
                    onChange={(e) => onSearchChange(e.target.value)}
                />
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-400" />
            </div>

            {filteredNetworks.length === 0 ? (
                <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                    No networks found matching "{searchTerm}"
                </div>
            ) : (
                <div className="relative">
                    {/* Tooltip for hovered network */}
                    {hoveredNetwork && (
                        <div className="absolute top-0 right-0 bg-white dark:bg-gray-800 p-3 rounded-lg shadow-md z-10 border dark:border-gray-700">
                            <h4 className="font-medium text-gray-800 dark:text-gray-200">{hoveredNetwork}</h4>
                            <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                                Status: <span className={`font-medium ${
                                statusData[hoveredNetwork]?.status === 'online' ? 'text-green-600 dark:text-green-400' :
                                    statusData[hoveredNetwork]?.status === 'offline' ? 'text-red-600 dark:text-red-400' :
                                        statusData[hoveredNetwork]?.status === 'warning' ? 'text-yellow-600 dark:text-yellow-400' :
                                            'text-gray-600 dark:text-gray-400'
                            }`}>
                  {statusData[hoveredNetwork]?.status || 'Unknown'}
                </span>
                            </div>
                            {statusData[hoveredNetwork]?.responseTime && (
                                <div className="text-sm text-gray-600 dark:text-gray-400">
                                    Response time: {statusData[hoveredNetwork]?.responseTime}
                                </div>
                            )}
                        </div>
                    )}

                    {/* Absolute simplest implementation - just use HexGrid as a layout */}
                    <div style={{ height: `${containerSize.height}px`, width: '100%', position: 'relative' }}>
                        <HexGrid width={containerSize.width} height={containerSize.height}>
                            <Layout size={{ x: 3, y: 3 }} flat={false} origin={{ x: 0, y: 0 }}>
                                {filteredNetworks.map((network, i) => {
                                    const status = statusData[network]?.status || 'unknown';
                                    // Calculate spiral position
                                    const angle = 0.5 * i;
                                    const q = Math.round(Math.cos(angle) * Math.sqrt(i) * 0.6);
                                    const r = Math.round(Math.sin(angle) * Math.sqrt(i) * 0.6);
                                    const s = -q - r;

                                    return (
                                        <Hexagon
                                            key={i}
                                            q={q}
                                            r={r}
                                            s={s}
                                            fill={getStatusColor(status)}
                                            onMouseEnter={() => setHoveredNetwork(network)}
                                            onMouseLeave={() => setHoveredNetwork(null)}
                                            className="hex-cell"
                                            style={{
                                                filter: "drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1))",
                                                transition: "all 0.2s ease",
                                                stroke: "#4b5563",
                                                strokeWidth: 0.2
                                            }}
                                        >
                                            <text
                                                x="0"
                                                y="0.25"
                                                style={{
                                                    fill: 'white',
                                                    fontSize: '0.4rem',
                                                    textAnchor: 'middle',
                                                    fontWeight: 'bold',
                                                    pointerEvents: 'none'
                                                }}
                                            >
                                                {network.split('/')[0].split('.')[3] || ''}
                                            </text>
                                        </Hexagon>
                                    );
                                })}
                            </Layout>
                        </HexGrid>
                    </div>
                </div>
            )}

            {/* Legend */}
            <div className="flex flex-wrap gap-4 mt-4 justify-center">
                <div className="flex items-center">
                    <div className="w-4 h-4 rounded-full bg-green-500 mr-2"></div>
                    <span className="text-sm text-gray-700 dark:text-gray-300">Online</span>
                </div>
                <div className="flex items-center">
                    <div className="w-4 h-4 rounded-full bg-red-500 mr-2"></div>
                    <span className="text-sm text-gray-700 dark:text-gray-300">Offline</span>
                </div>
                <div className="flex items-center">
                    <div className="w-4 h-4 rounded-full bg-yellow-500 mr-2"></div>
                    <span className="text-sm text-gray-700 dark:text-gray-300">Warning</span>
                </div>
                <div className="flex items-center">
                    <div className="w-4 h-4 rounded-full bg-gray-500 mr-2"></div>
                    <span className="text-sm text-gray-700 dark:text-gray-300">Unknown</span>
                </div>
            </div>
        </div>
    );
};

export default HoneycombNetworkGrid;