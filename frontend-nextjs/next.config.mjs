/** @type {import('next').NextConfig} */
const nextConfig = {
    experimental: {
        after: true
    },
    images: {
        remotePatterns: [
            {
                hostname: 'python_app'
            },
            {
                hostname: 'nginx'
            },
            {
                hostname: 'host.docker.internal'
            },
        ]
    }
};

export default nextConfig;
