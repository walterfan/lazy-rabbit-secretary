#!/usr/bin/env node

/**
 * Test script for loading secrets from the API
 * Usage: node test-secrets-api.js [username] [password]
 */

const https = require('https');

// Disable SSL certificate verification for self-signed certificates
process.env["NODE_TLS_REJECT_UNAUTHORIZED"] = 0;

// Configuration
const config = {
    useProxy: true, // Set to false to connect directly to backend at port 9090
    proxyUrl: 'https://localhost:5173',
    backendUrl: 'https://localhost:9090'
};

const baseUrl = config.useProxy ? config.proxyUrl : config.backendUrl;

// Get credentials from command line or use defaults
const username = process.argv[2] || 'admin';
const password = process.argv[3] || '';

if (!password) {
    console.error('Usage: node test-secrets-api.js [username] [password]');
    console.error('Example: node test-secrets-api.js admin mypassword');
    process.exit(1);
}

// Helper function to make HTTPS requests
function makeRequest(options, postData = null) {
    return new Promise((resolve, reject) => {
        const req = https.request(options, (res) => {
            let data = '';
            
            res.on('data', (chunk) => {
                data += chunk;
            });
            
            res.on('end', () => {
                try {
                    const jsonData = JSON.parse(data);
                    resolve({ statusCode: res.statusCode, data: jsonData });
                } catch (e) {
                    resolve({ statusCode: res.statusCode, data: data });
                }
            });
        });
        
        req.on('error', (error) => {
            reject(error);
        });
        
        if (postData) {
            req.write(postData);
        }
        
        req.end();
    });
}

// Sign in to get access token
async function signIn() {
    console.log(`\n1. Signing in as ${username}...`);
    
    const postData = JSON.stringify({ username, password });
    
    const options = {
        hostname: baseUrl.replace('https://', '').split(':')[0],
        port: baseUrl.includes(':') ? baseUrl.split(':')[2] : 443,
        path: '/api/v1/auth/login',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': postData.length
        }
    };
    
    try {
        const response = await makeRequest(options, postData);
        
        if (response.statusCode === 200 && response.data.access_token) {
            console.log('✓ Sign in successful!');
            console.log(`  Access Token: ${response.data.access_token.substring(0, 20)}...`);
            return response.data.access_token;
        } else {
            console.error('✗ Sign in failed:', response.data);
            return null;
        }
    } catch (error) {
        console.error('✗ Network error during sign in:', error.message);
        return null;
    }
}

// Load secrets from API
async function loadSecrets(accessToken) {
    console.log('\n2. Loading secrets from /api/v1/secrets...');
    
    const options = {
        hostname: baseUrl.replace('https://', '').split(':')[0],
        port: baseUrl.includes(':') ? baseUrl.split(':')[2] : 443,
        path: '/api/v1/secrets',
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json'
        }
    };
    
    try {
        const response = await makeRequest(options);
        
        if (response.statusCode === 200) {
            console.log('✓ Secrets loaded successfully!');
            
            const data = response.data;
            console.log(`\nTotal secrets: ${data.total || 0}`);
            
            if (data.items && data.items.length > 0) {
                console.log('\nSecrets List:');
                console.log('=' .repeat(80));
                
                data.items.forEach((secret, index) => {
                    console.log(`\nSecret ${index + 1}:`);
                    console.log(`  ID:          ${secret.id}`);
                    console.log(`  Name:        ${secret.name}`);
                    console.log(`  Group:       ${secret.group}`);
                    console.log(`  Path:        ${secret.path}`);
                    console.log(`  Description: ${secret.desc || 'N/A'}`);
                    console.log(`  Encryption:  ${secret.cipher_alg}`);
                    console.log(`  KEK Version: ${secret.kek_version}`);
                    console.log(`  Created By:  ${secret.created_by}`);
                    console.log(`  Updated:     ${new Date(secret.updated_at).toLocaleString()}`);
                });
                
                console.log('\n' + '=' .repeat(80));
            } else {
                console.log('\nNo secrets found in the system.');
            }
            
            return data;
        } else {
            console.error('✗ Failed to load secrets:', response.data);
            return null;
        }
    } catch (error) {
        console.error('✗ Network error while loading secrets:', error.message);
        return null;
    }
}

// Main execution
async function main() {
    console.log('Secrets API Test');
    console.log('================');
    console.log(`Using ${config.useProxy ? 'Vite proxy' : 'direct backend'} at ${baseUrl}`);
    
    // Step 1: Sign in
    const accessToken = await signIn();
    if (!accessToken) {
        console.error('\nFailed to authenticate. Please check your credentials.');
        process.exit(1);
    }
    
    // Step 2: Load secrets
    const secrets = await loadSecrets(accessToken);
    if (!secrets) {
        console.error('\nFailed to load secrets.');
        process.exit(1);
    }
    
    console.log('\n✓ Test completed successfully!');
}

// Run the test
main().catch(error => {
    console.error('Unexpected error:', error);
    process.exit(1);
});
