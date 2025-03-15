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

// src/lib/networkUtils.js

/**
 * Check if an IP address is within a CIDR subnet
 * @param {string} ip - IP address to check
 * @param {string} cidr - CIDR notation subnet (e.g. "192.168.1.0/24")
 * @returns {boolean} - Whether the IP is in the subnet
 */
export function isIpInCidr(ip, cidr) {
    try {
        // Extract network and mask from CIDR notation
        const [subnetIp, maskBits] = cidr.split('/');
        const mask = parseInt(maskBits, 10);

        // Validate inputs
        if (mask < 0 || mask > 32) return false;

        // Convert IP addresses to integer representations
        const ipInt = ipToInt(ip);
        const subnetIpInt = ipToInt(subnetIp);

        // Calculate subnet mask
        const subnetMask = (0xffffffff << (32 - mask)) >>> 0;

        // Compare network portions
        const ipNetwork = ipInt & subnetMask;
        const subnetNetwork = subnetIpInt & subnetMask;

        return ipNetwork === subnetNetwork;
    } catch (e) {
        console.error('Error in isIpInCidr:', e);
        return false;
    }
}

/**
 * Convert an IP address to its integer representation
 * @param {string} ip - IP address (e.g. "192.168.1.1")
 * @returns {number} - Integer representation of the IP
 */
export function ipToInt(ip) {
    try {
        const octets = ip.split('.');
        if (octets.length !== 4) throw new Error('Invalid IP address format');

        return ((parseInt(octets[0], 10) << 24) |
            (parseInt(octets[1], 10) << 16) |
            (parseInt(octets[2], 10) << 8) |
            parseInt(octets[3], 10)) >>> 0; // >>> 0 ensures unsigned 32-bit integer
    } catch (e) {
        console.error('Error in ipToInt:', e);
        return 0;
    }
}

/**
 * Convert an integer to IP address
 * @param {number} int - Integer representation of IP
 * @returns {string} - IP address in string format
 */
export function intToIp(int) {
    try {
        return [
            (int >>> 24) & 255,
            (int >>> 16) & 255,
            (int >>> 8) & 255,
            int & 255
        ].join('.');
    } catch (e) {
        console.error('Error in intToIp:', e);
        return '0.0.0.0';
    }
}

/**
 * Group hosts by network
 * @param {Array} hosts - Array of host objects with IP addresses
 * @param {Array} networks - Array of network CIDR strings
 * @returns {Object} - Object with networks as keys and arrays of hosts as values
 */
export function groupHostsByNetwork(hosts, networks) {
    if (!hosts || !networks) return {};

    const result = {};

    // Initialize empty arrays for each network
    networks.forEach(network => {
        result[network] = [];
    });

    // Assign hosts to networks
    hosts.forEach(host => {
        if (!host.host) return;

        // Find which network this host belongs to
        for (const network of networks) {
            if (isIpInCidr(host.host, network)) {
                result[network].push(host);
                break;
            }
        }
    });

    return result;
}

/**
 * Simple version of network consolidation that focuses on clarity for UI
 * @param {Array} ips - Array of IP addresses as strings
 * @returns {Array} - Array of network objects with CIDR notation
 */
export function consolidateNetworks(ips) {
    if (!ips || !ips.length) return [];

    // Group IPs by first 3 octets to find /24 networks
    const networkGroups = {};
    ips.forEach(ip => {
        if (!ip) return;

        // Get first three octets as a potential /24 network
        const parts = ip.split('.');
        if (parts.length !== 4) return;

        const networkPrefix = `${parts[0]}.${parts[1]}.${parts[2]}`;
        if (!networkGroups[networkPrefix]) {
            networkGroups[networkPrefix] = [];
        }
        networkGroups[networkPrefix].push(ip);
    });

    // Convert to network objects with more data for UI display
    const networks = Object.entries(networkGroups).map(([prefix, addresses]) => ({
        network: `${prefix}.0/24`,
        hosts: addresses.length,
        coverage: addresses.length / 256, // Coverage relative to a /24 network
        addresses
    }));

    // Sort networks by size (number of hosts) descending
    networks.sort((a, b) => b.hosts - a.hosts);

    return networks;
}

/**
 * Converts list of IPs to optimal CIDR blocks
 * Simplified for better visualization
 * @param {Array} ips - List of IP addresses
 * @returns {Array} - List of optimal CIDR blocks
 */
export function ipsToCIDRBlocks(ips) {
    if (!ips || !ips.length) return [];

    // Create basic network groups for better visualization
    const networks = [];

    // First try to group by class C networks (/24)
    const classC = {};

    ips.forEach(ip => {
        const parts = ip.split('.');
        if (parts.length !== 4) return;

        const c = `${parts[0]}.${parts[1]}.${parts[2]}`;
        if (!classC[c]) classC[c] = [];
        classC[c].push(ip);
    });

    // Create networks from class C groups
    Object.entries(classC).forEach(([prefix, ips]) => {
        networks.push({
            network: `${prefix}.0/24`,
            hosts: ips.length,
            ips
        });
    });

    // If we have too many small networks, try to consolidate to class B (/16)
    if (networks.length > 10) {
        const classB = {};

        ips.forEach(ip => {
            const parts = ip.split('.');
            if (parts.length !== 4) return;

            const b = `${parts[0]}.${parts[1]}`;
            if (!classB[b]) classB[b] = [];
            classB[b].push(ip);
        });

        // If a class B has enough IPs, use that instead of its component class Cs
        const consolidatedNetworks = [];

        Object.entries(classB).forEach(([prefix, ips]) => {
            // If this classB has more than 5 IPs, add it
            if (ips.length > 5) {
                consolidatedNetworks.push({
                    network: `${prefix}.0.0/16`,
                    hosts: ips.length,
                    ips
                });

                // Remove the component class C networks
                const prefixMatch = new RegExp(`^${prefix}\\.`);
                for (let i = networks.length - 1; i >= 0; i--) {
                    if (prefixMatch.test(networks[i].network)) {
                        networks.splice(i, 1);
                    }
                }
            }
        });

        // Add the consolidated classB networks
        networks.push(...consolidatedNetworks);
    }

    // Sort networks by size for better display
    networks.sort((a, b) => b.hosts - a.hosts);

    return networks;
}